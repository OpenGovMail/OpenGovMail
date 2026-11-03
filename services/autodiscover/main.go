package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	serverReadHeaderTimeout = 5 * time.Second
	serverReadTimeout       = 10 * time.Second
	serverWriteTimeout      = 15 * time.Second
	serverIdleTimeout       = 60 * time.Second
	serverMaxHeaderBytes    = 1 << 14  // 16 KB
	maxRequestBodyBytes     = 64 << 10 // 64 KiB
)

// Config mirrors the subset of conf/opengovmail.yaml this service needs.
type Config struct {
	Domains []DomainConfig `yaml:"domains"`
}

// DomainConfig is one mail domain from opengovmail.yaml.
type DomainConfig struct {
	Domain string `yaml:"domain"`
}

// service holds the validated domain settings shared by every handler.
type service struct {
	primary  string
	mailHost string
	domains  map[string]struct{}
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}
	if len(cfg.Domains) == 0 {
		return nil, fmt.Errorf("config %s has no domains", path)
	}

	for i := range cfg.Domains {
		cfg.Domains[i].Domain = normalizeHost(cfg.Domains[i].Domain)
		if cfg.Domains[i].Domain == "" {
			return nil, fmt.Errorf("config %s has an empty domain at index %d", path, i)
		}
	}
	return &cfg, nil
}

func newService(cfg *Config, mailHost string) (*service, error) {
	if cfg == nil || len(cfg.Domains) == 0 {
		return nil, errors.New("at least one domain is required")
	}

	s := &service{domains: make(map[string]struct{}, len(cfg.Domains))}
	for i, d := range cfg.Domains {
		domain := normalizeHost(d.Domain)
		if domain == "" {
			return nil, fmt.Errorf("domain at index %d is empty", i)
		}
		if s.primary == "" {
			s.primary = domain
		}
		s.domains[domain] = struct{}{}
	}

	s.mailHost = normalizeHost(mailHost)
	if s.mailHost == "" {
		s.mailHost = "mail." + s.primary
	}
	return s, nil
}

func newHandler(cfg *Config, mailHost string) (http.Handler, error) {
	s, err := newService(cfg, mailHost)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.healthz)
	mux.HandleFunc("GET /mail/config-v1.1.xml", s.thunderbird)
	mux.HandleFunc("GET /.well-known/autoconfig/mail/config-v1.1.xml", s.thunderbird)
	mux.HandleFunc("POST /autodiscover/autodiscover.xml", s.outlook)
	mux.HandleFunc("POST /Autodiscover/Autodiscover.xml", s.outlook)

	return loggingMiddleware(mux), nil
}

func (s *service) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *service) thunderbird(w http.ResponseWriter, r *http.Request) {
	address := strings.TrimSpace(r.URL.Query().Get("emailaddress"))
	domain, err := emailDomain(address)
	if err != nil {
		http.Error(w, "a valid emailaddress query parameter is required", http.StatusBadRequest)
		return
	}
	if !s.knownDomain(domain) {
		http.NotFound(w, r)
		return
	}

	cfg := thunderbirdClientConfig{
		Version: "1.1",
		Provider: thunderbirdProvider{
			ID:          domain,
			Domain:      domain,
			DisplayName: displayName(domain),
			IncomingServer: thunderbirdServer{
				Type:           "imap",
				Hostname:       s.mailHost,
				Port:           993,
				SocketType:     "SSL",
				Authentication: "password-cleartext",
				Username:       "%EMAILADDRESS%",
			},
			OutgoingServer: thunderbirdServer{
				Type:           "smtp",
				Hostname:       s.mailHost,
				Port:           587,
				SocketType:     "STARTTLS",
				Authentication: "password-cleartext",
				Username:       "%EMAILADDRESS%",
			},
		},
	}

	body, err := xml.MarshalIndent(cfg, "", "  ")
	if err != nil {
		http.Error(w, "failed to build configuration", http.StatusInternalServerError)
		return
	}
	writeXML(w, http.StatusOK, append([]byte(xml.Header), body...))
}

func (s *service) outlook(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			writeOutlookError(w, http.StatusRequestEntityTooLarge)
			return
		}
		writeOutlookError(w, http.StatusBadRequest)
		return
	}

	var req outlookPoxRequest
	if err := xml.Unmarshal(body, &req); err != nil {
		writeOutlookError(w, http.StatusBadRequest)
		return
	}

	address := strings.TrimSpace(req.Request.EMailAddress)
	domain, err := emailDomain(address)
	if err != nil {
		writeOutlookError(w, http.StatusBadRequest)
		return
	}
	if !s.knownDomain(domain) {
		writeOutlookError(w, http.StatusNotFound)
		return
	}

	response := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<Autodiscover xmlns="http://schemas.microsoft.com/exchange/autodiscover/responseschema/2006">
  <Response xmlns="http://schemas.microsoft.com/exchange/autodiscover/outlook/responseschema/2006a">
    <Account>
      <AccountType>email</AccountType>
      <Action>settings</Action>
      <Protocol>
        <Type>IMAP</Type>
        <Server>%s</Server>
        <Port>993</Port>
        <DomainRequired>off</DomainRequired>
        <LoginName>%s</LoginName>
        <SPA>off</SPA>
        <SSL>on</SSL>
        <Encryption>SSL</Encryption>
        <AuthRequired>on</AuthRequired>
      </Protocol>
      <Protocol>
        <Type>SMTP</Type>
        <Server>%s</Server>
        <Port>587</Port>
        <DomainRequired>off</DomainRequired>
        <LoginName>%s</LoginName>
        <SPA>off</SPA>
        <SSL>off</SSL>
        <Encryption>TLS</Encryption>
        <AuthRequired>on</AuthRequired>
      </Protocol>
    </Account>
  </Response>
</Autodiscover>
`, xmlEscape(s.mailHost), xmlEscape(address), xmlEscape(s.mailHost), xmlEscape(address))

	writeXML(w, http.StatusOK, []byte(response))
}

func (s *service) knownDomain(domain string) bool {
	_, ok := s.domains[domain]
	return ok
}

func writeXML(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	w.Write(body)
}

func writeOutlookError(w http.ResponseWriter, status int) {
	now := time.Now().UTC().Format(time.RFC3339)
	body := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<Autodiscover xmlns="http://schemas.microsoft.com/exchange/autodiscover/responseschema/2006">
  <Response>
    <Error Time="%s" Id="%d">
      <ErrorCode>600</ErrorCode>
      <Message>Invalid Request</Message>
      <DebugData />
    </Error>
  </Response>
</Autodiscover>
`, now, time.Now().UnixNano())
	writeXML(w, status, []byte(body))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		// Log the path only: the emailaddress query parameter is personal data.
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func emailDomain(address string) (string, error) {
	if address == "" || strings.TrimSpace(address) != address {
		return "", errors.New("empty or padded email address")
	}
	if strings.ContainsAny(address, " \t\r\n") {
		return "", errors.New("email address contains whitespace")
	}
	at := strings.LastIndex(address, "@")
	if at <= 0 || at == len(address)-1 {
		return "", errors.New("email address must contain a local part and a domain")
	}
	if strings.Contains(address[:at], "@") {
		return "", errors.New("email address contains more than one @")
	}
	domain := normalizeHost(address[at+1:])
	if domain == "" || strings.ContainsAny(domain, "@/") {
		return "", errors.New("invalid email domain")
	}
	return domain, nil
}

func normalizeHost(host string) string {
	return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(host)), ".")
}

func displayName(domain string) string {
	return domain
}

func xmlEscape(value string) string {
	var buf bytes.Buffer
	xml.EscapeText(&buf, []byte(value))
	return buf.String()
}

// outlookPoxRequest is the POX request Outlook sends to /autodiscover/autodiscover.xml.
type outlookPoxRequest struct {
	XMLName xml.Name `xml:"Autodiscover"`
	Request struct {
		EMailAddress             string `xml:"EMailAddress"`
		AcceptableResponseSchema string `xml:"AcceptableResponseSchema"`
	} `xml:"Request"`
}

type thunderbirdClientConfig struct {
	XMLName  xml.Name            `xml:"clientConfig"`
	Version  string              `xml:"version,attr"`
	Provider thunderbirdProvider `xml:"emailProvider"`
}

type thunderbirdProvider struct {
	ID             string            `xml:"id,attr"`
	Domain         string            `xml:"domain"`
	DisplayName    string            `xml:"displayName"`
	IncomingServer thunderbirdServer `xml:"incomingServer"`
	OutgoingServer thunderbirdServer `xml:"outgoingServer"`
}

type thunderbirdServer struct {
	Type           string `xml:"type,attr"`
	Hostname       string `xml:"hostname"`
	Port           int    `xml:"port"`
	SocketType     string `xml:"socketType"`
	Authentication string `xml:"authentication"`
	Username       string `xml:"username"`
}

func newServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: serverReadHeaderTimeout,
		ReadTimeout:       serverReadTimeout,
		WriteTimeout:      serverWriteTimeout,
		IdleTimeout:       serverIdleTimeout,
		MaxHeaderBytes:    serverMaxHeaderBytes,
	}
}

// certReloader loads the Let's Encrypt certificate for the primary domain and
// reloads it when certbot replaces the files.
type certReloader struct {
	fullchainPath string
	privkeyPath   string

	mu       sync.RWMutex
	cert     *tls.Certificate
	modTime  time.Time
	lastErr  error
	lastLoad time.Time
}

func newCertReloader(certDir, primaryDomain string) (*certReloader, error) {
	lineage := filepath.Join(certDir, primaryDomain)
	r := &certReloader{
		fullchainPath: filepath.Join(lineage, "fullchain.pem"),
		privkeyPath:   filepath.Join(lineage, "privkey.pem"),
	}
	if err := r.reload(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *certReloader) reload() error {
	fullchainInfo, err := os.Stat(r.fullchainPath)
	if err != nil {
		return fmt.Errorf("stat %s: %w", r.fullchainPath, err)
	}
	privkeyInfo, err := os.Stat(r.privkeyPath)
	if err != nil {
		return fmt.Errorf("stat %s: %w", r.privkeyPath, err)
	}

	cert, err := tls.LoadX509KeyPair(r.fullchainPath, r.privkeyPath)
	if err != nil {
		return fmt.Errorf("load certificate: %w", err)
	}

	modTime := fullchainInfo.ModTime()
	if privkeyInfo.ModTime().After(modTime) {
		modTime = privkeyInfo.ModTime()
	}

	r.mu.Lock()
	r.cert = &cert
	r.modTime = modTime
	r.lastErr = nil
	r.lastLoad = time.Now()
	r.mu.Unlock()
	return nil
}

func (r *certReloader) getCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	r.mu.RLock()
	modTime := r.modTime
	cert := r.cert
	lastLoad := r.lastLoad
	r.mu.RUnlock()

	// Check at most once a minute to keep handshakes cheap. A renewal that lands
	// just after a check is picked up on the next one.
	if cert == nil || time.Since(lastLoad) > time.Minute {
		if err := r.reload(); err != nil {
			r.mu.Lock()
			r.lastErr = err
			r.mu.Unlock()
			if cert == nil {
				return nil, err
			}
		} else {
			r.mu.RLock()
			cert = r.cert
			modTime = r.modTime
			r.mu.RUnlock()
		}
	}

	// Stat both files on every handshake; certbot replaces them atomically and a
	// changed mtime is a cheap signal that a reload is due.
	fullchainInfo, err := os.Stat(r.fullchainPath)
	if err == nil && fullchainInfo.ModTime().After(modTime) {
		if err := r.reload(); err != nil {
			r.mu.Lock()
			r.lastErr = err
			r.mu.Unlock()
		} else {
			r.mu.RLock()
			cert = r.cert
			r.mu.RUnlock()
		}
	}

	if cert == nil {
		r.mu.RLock()
		err := r.lastErr
		r.mu.RUnlock()
		if err == nil {
			err = errors.New("certificate is not loaded")
		}
		return nil, err
	}
	return cert, nil
}

func buildTLSConfig(certDir, primaryDomain string) (*tls.Config, error) {
	reloader, err := newCertReloader(certDir, primaryDomain)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		MinVersion:     tls.VersionTLS12,
		GetCertificate: reloader.getCertificate,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return defaultValue
	}
	value, err := strconvParseBool(raw)
	if err != nil {
		log.Printf("WARN: invalid %s=%q, using default %t", key, raw, defaultValue)
		return defaultValue
	}
	return value
}

// strconvParseBool keeps main.go's imports small and makes the accepted values
// explicit: true/false, 1/0, t/f and their case variants.
func strconvParseBool(raw string) (bool, error) {
	switch strings.ToLower(raw) {
	case "1", "t", "true":
		return true, nil
	case "0", "f", "false":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean %q", raw)
	}
}

func main() {
	log.Println("OpenGovMail Autodiscover Service (Go)")

	configFile := getEnv("CONFIG_FILE", "/etc/autodiscover/opengovmail.yaml")
	port := getEnv("PORT", "8080")
	mailHost := getEnv("MAIL_HOSTNAME", "")
	tlsEnabled := getEnvAsBool("TLS_ENABLED", false)
	tlsCertDir := getEnv("TLS_CERT_DIR", "/etc/letsencrypt/live")

	cfg, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}
	s, err := newService(cfg, mailHost)
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}
	handler, err := newHandler(cfg, mailHost)
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	addr := "0.0.0.0:" + port
	srv := newServer(addr, handler)

	if tlsEnabled {
		tlsConfig, err := buildTLSConfig(tlsCertDir, s.primary)
		if err != nil {
			log.Fatalf("TLS configuration error: %v", err)
		}
		srv.TLSConfig = tlsConfig
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("Serving %s (TLS: %t, mail host: %s, primary domain: %s)", addr, tlsEnabled, s.mailHost, s.primary)
		if tlsEnabled {
			serveErr <- srv.ListenAndServeTLS("", "")
		} else {
			serveErr <- srv.ListenAndServe()
		}
	}()

	select {
	case err := <-serveErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	case <-ctx.Done():
		log.Println("Shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}
}
