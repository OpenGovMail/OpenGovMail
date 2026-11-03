package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const testMaxRequestBodyBytes = 64 << 10 // 64 KiB

func TestLoadConfigReadsDomains(t *testing.T) {
	path := writeTestConfig(t, `config-url: "https://github.com/OpenGovMail/config"
domains:
  - domain: example.com
    dkim-selector: mail
    dkim-key-size: 2048
  - domain: example.org
`)

	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}
	if len(cfg.Domains) != 2 {
		t.Fatalf("len(Domains) = %d, want 2", len(cfg.Domains))
	}
	if cfg.Domains[0].Domain != "example.com" {
		t.Errorf("Domains[0].Domain = %q, want example.com", cfg.Domains[0].Domain)
	}
	if cfg.Domains[1].Domain != "example.org" {
		t.Errorf("Domains[1].Domain = %q, want example.org", cfg.Domains[1].Domain)
	}
}

func TestLoadConfigRejectsNoDomains(t *testing.T) {
	path := writeTestConfig(t, "config-url: \"https://github.com/OpenGovMail/config\"\ndomains: []\n")

	if _, err := loadConfig(path); err == nil {
		t.Fatal("loadConfig() error = nil, want an error for an empty domains list")
	}
}

func TestLoadConfigRejectsEmptyDomainName(t *testing.T) {
	path := writeTestConfig(t, "domains:\n  - domain: \"\"\n")

	if _, err := loadConfig(path); err == nil {
		t.Fatal("loadConfig() error = nil, want an error for an empty domain")
	}
}

func TestNewHandlerDefaultsMailHostToPrimaryDomain(t *testing.T) {
	h := mustHandler(t, testConfig("example.com", "example.org"), "")

	rec := doRequest(h, http.MethodGet, "/mail/config-v1.1.xml?emailaddress=user@example.org", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	cfg := decodeThunderbirdConfig(t, rec.Body.Bytes())
	if got := cfg.Provider.IncomingServer.Hostname; got != "mail.example.com" {
		t.Errorf("incoming hostname = %q, want mail.example.com", got)
	}
}

func TestNewHandlerRejectsNoDomains(t *testing.T) {
	if _, err := newHandler(&Config{}, ""); err == nil {
		t.Fatal("newHandler() error = nil, want an error for an empty domains list")
	}
}

func TestThunderbirdEndpointReturnsKnownDomainSettings(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")

	rec := doRequest(h, http.MethodGet, "/.well-known/autoconfig/mail/config-v1.1.xml?emailaddress=user@example.com", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/xml") {
		t.Errorf("Content-Type = %q, want text/xml", ct)
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", got)
	}

	cfg := decodeThunderbirdConfig(t, rec.Body.Bytes())
	if cfg.Version != "1.1" {
		t.Errorf("clientConfig version = %q, want 1.1", cfg.Version)
	}
	if cfg.Provider.ID != "example.com" {
		t.Errorf("emailProvider id = %q, want example.com", cfg.Provider.ID)
	}
	if cfg.Provider.IncomingServer.Type != "imap" {
		t.Errorf("incoming type = %q, want imap", cfg.Provider.IncomingServer.Type)
	}
	if cfg.Provider.IncomingServer.Port != 993 {
		t.Errorf("incoming port = %d, want 993", cfg.Provider.IncomingServer.Port)
	}
	if cfg.Provider.IncomingServer.SocketType != "SSL" {
		t.Errorf("incoming socketType = %q, want SSL", cfg.Provider.IncomingServer.SocketType)
	}
	if cfg.Provider.IncomingServer.Username != "%EMAILADDRESS%" {
		t.Errorf("incoming username = %q, want %%EMAILADDRESS%%", cfg.Provider.IncomingServer.Username)
	}
	if cfg.Provider.OutgoingServer.Type != "smtp" {
		t.Errorf("outgoing type = %q, want smtp", cfg.Provider.OutgoingServer.Type)
	}
	if cfg.Provider.OutgoingServer.Port != 587 {
		t.Errorf("outgoing port = %d, want 587", cfg.Provider.OutgoingServer.Port)
	}
	if cfg.Provider.OutgoingServer.SocketType != "STARTTLS" {
		t.Errorf("outgoing socketType = %q, want STARTTLS", cfg.Provider.OutgoingServer.SocketType)
	}
	if cfg.Provider.OutgoingServer.Username != "%EMAILADDRESS%" {
		t.Errorf("outgoing username = %q, want %%EMAILADDRESS%%", cfg.Provider.OutgoingServer.Username)
	}
}

func TestThunderbirdEndpointServesBothPaths(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")

	for _, path := range []string{
		"/mail/config-v1.1.xml?emailaddress=user@example.com",
		"/.well-known/autoconfig/mail/config-v1.1.xml?emailaddress=user@example.com",
	} {
		rec := doRequest(h, http.MethodGet, path, "")
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s status = %d, want %d", path, rec.Code, http.StatusOK)
		}
	}
}

func TestThunderbirdEndpointRejectsUnknownDomain(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")

	rec := doRequest(h, http.MethodGet, "/mail/config-v1.1.xml?emailaddress=user@unknown.test", "")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestThunderbirdEndpointRejectsMissingOrInvalidEmail(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")

	for _, path := range []string{
		"/mail/config-v1.1.xml",
		"/mail/config-v1.1.xml?emailaddress=not-an-email",
		"/mail/config-v1.1.xml?emailaddress=user@",
	} {
		rec := doRequest(h, http.MethodGet, path, "")
		if rec.Code != http.StatusBadRequest {
			t.Errorf("GET %s status = %d, want %d", path, rec.Code, http.StatusBadRequest)
		}
	}
}

func TestThunderbirdEndpointRejectsPost(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")

	rec := doRequest(h, http.MethodPost, "/mail/config-v1.1.xml?emailaddress=user@example.com", "")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestOutlookEndpointReturnsIMAPAndSMTPSettings(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")

	rec := doRequest(h, http.MethodPost, "/autodiscover/autodiscover.xml", outlookRequest("user@example.com"))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/xml") {
		t.Errorf("Content-Type = %q, want text/xml", ct)
	}

	resp := decodeOutlookResponse(t, rec.Body.Bytes())
	if resp.Response.Account.AccountType != "email" {
		t.Errorf("AccountType = %q, want email", resp.Response.Account.AccountType)
	}
	if resp.Response.Account.Action != "settings" {
		t.Errorf("Action = %q, want settings", resp.Response.Account.Action)
	}
	if len(resp.Response.Account.Protocols) != 2 {
		t.Fatalf("len(Protocols) = %d, want 2", len(resp.Response.Account.Protocols))
	}

	imap := findProtocol(t, resp, "IMAP")
	if imap.Server != "mail.example.com" {
		t.Errorf("IMAP Server = %q, want mail.example.com", imap.Server)
	}
	if imap.Port != 993 {
		t.Errorf("IMAP Port = %d, want 993", imap.Port)
	}
	if imap.LoginName != "user@example.com" {
		t.Errorf("IMAP LoginName = %q, want user@example.com", imap.LoginName)
	}
	if imap.SSL != "on" {
		t.Errorf("IMAP SSL = %q, want on", imap.SSL)
	}
	if imap.Encryption != "SSL" {
		t.Errorf("IMAP Encryption = %q, want SSL", imap.Encryption)
	}

	smtp := findProtocol(t, resp, "SMTP")
	if smtp.Server != "mail.example.com" {
		t.Errorf("SMTP Server = %q, want mail.example.com", smtp.Server)
	}
	if smtp.Port != 587 {
		t.Errorf("SMTP Port = %d, want 587", smtp.Port)
	}
	if smtp.LoginName != "user@example.com" {
		t.Errorf("SMTP LoginName = %q, want user@example.com", smtp.LoginName)
	}
	if smtp.SSL != "off" {
		t.Errorf("SMTP SSL = %q, want off", smtp.SSL)
	}
	if smtp.Encryption != "TLS" {
		t.Errorf("SMTP Encryption = %q, want TLS (STARTTLS)", smtp.Encryption)
	}
}

func TestOutlookEndpointServesBothCapitalizations(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")

	for _, path := range []string{"/autodiscover/autodiscover.xml", "/Autodiscover/Autodiscover.xml"} {
		rec := doRequest(h, http.MethodPost, path, outlookRequest("user@example.com"))
		if rec.Code != http.StatusOK {
			t.Errorf("POST %s status = %d, want %d", path, rec.Code, http.StatusOK)
		}
	}
}

func TestOutlookEndpointUsesPrimaryMailHostForSecondaryDomains(t *testing.T) {
	h := mustHandler(t, testConfig("example.com", "example.org"), "")

	rec := doRequest(h, http.MethodPost, "/autodiscover/autodiscover.xml", outlookRequest("user@example.org"))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	resp := decodeOutlookResponse(t, rec.Body.Bytes())
	imap := findProtocol(t, resp, "IMAP")
	if imap.Server != "mail.example.com" {
		t.Errorf("IMAP Server = %q, want mail.example.com", imap.Server)
	}
}

func TestOutlookEndpointRejectsMalformedXML(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")

	rec := doRequest(h, http.MethodPost, "/autodiscover/autodiscover.xml", "<Autodiscover>")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	assertOutlookError(t, rec.Body.Bytes(), "600")
}

func TestOutlookEndpointRejectsOversizedBody(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")
	body := strings.Repeat("a", testMaxRequestBodyBytes+1)

	rec := doRequest(h, http.MethodPost, "/autodiscover/autodiscover.xml", body)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestOutlookEndpointRejectsUnknownDomain(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")

	rec := doRequest(h, http.MethodPost, "/autodiscover/autodiscover.xml", outlookRequest("user@unknown.test"))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	assertOutlookError(t, rec.Body.Bytes(), "600")
}

func TestOutlookEndpointRejectsMissingEmailAddress(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")

	rec := doRequest(h, http.MethodPost, "/autodiscover/autodiscover.xml", `<?xml version="1.0"?><Autodiscover xmlns="http://schemas.microsoft.com/exchange/autodiscover/outlook/requestschema/2006"><Request></Request></Autodiscover>`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestOutlookEndpointRejectsGet(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")

	rec := doRequest(h, http.MethodGet, "/autodiscover/autodiscover.xml", "")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestHealthz(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")

	rec := doRequest(h, http.MethodGet, "/healthz", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	if !strings.Contains(rec.Body.String(), `"status":"healthy"`) {
		t.Errorf("body = %s, want status healthy", rec.Body.String())
	}
}

func TestUnknownPathReturnsNotFound(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")

	rec := doRequest(h, http.MethodGet, "/not-a-real-endpoint", "")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestNewServerConfiguresTimeoutsAndLimits(t *testing.T) {
	srv := newServer("0.0.0.0:8080", http.NewServeMux())

	if srv.Addr != "0.0.0.0:8080" {
		t.Errorf("Addr = %q, want 0.0.0.0:8080", srv.Addr)
	}
	if srv.Handler == nil {
		t.Error("Handler is nil, want the configured handler")
	}
	if srv.ReadHeaderTimeout != 5*time.Second {
		t.Errorf("ReadHeaderTimeout = %v, want %v", srv.ReadHeaderTimeout, 5*time.Second)
	}
	if srv.ReadTimeout != 10*time.Second {
		t.Errorf("ReadTimeout = %v, want %v", srv.ReadTimeout, 10*time.Second)
	}
	if srv.WriteTimeout != 15*time.Second {
		t.Errorf("WriteTimeout = %v, want %v", srv.WriteTimeout, 15*time.Second)
	}
	if srv.IdleTimeout != 60*time.Second {
		t.Errorf("IdleTimeout = %v, want %v", srv.IdleTimeout, 60*time.Second)
	}
	if srv.MaxHeaderBytes != 1<<14 {
		t.Errorf("MaxHeaderBytes = %d, want %d", srv.MaxHeaderBytes, 1<<14)
	}
}

func TestServerServesHealthOverTCP(t *testing.T) {
	h := mustHandler(t, testConfig("example.com"), "")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	srv := newServer(ln.Addr().String(), h)
	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(ln) }()
	defer srv.Close()

	resp, err := http.Get("http://" + ln.Addr().String() + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestBuildTLSConfigLoadsCertificate(t *testing.T) {
	certDir := t.TempDir()
	writeTestCertificate(t, certDir, "example.com", 1)

	tlsConfig, err := buildTLSConfig(certDir, "example.com")
	if err != nil {
		t.Fatalf("buildTLSConfig() error = %v", err)
	}
	if tlsConfig.MinVersion != tls.VersionTLS12 {
		t.Errorf("MinVersion = %d, want TLS 1.2", tlsConfig.MinVersion)
	}

	cert, err := tlsConfig.GetCertificate(&tls.ClientHelloInfo{})
	if err != nil {
		t.Fatalf("GetCertificate() error = %v", err)
	}
	if got := certificateSerial(t, cert); got != 1 {
		t.Errorf("serial = %d, want 1", got)
	}
}

func TestBuildTLSConfigRejectsMissingCertificate(t *testing.T) {
	if _, err := buildTLSConfig(t.TempDir(), "example.com"); err == nil {
		t.Fatal("buildTLSConfig() error = nil, want an error for missing files")
	}
}

func TestCertificateReloaderPicksUpRenewedCertificate(t *testing.T) {
	certDir := t.TempDir()
	writeTestCertificate(t, certDir, "example.com", 1)

	tlsConfig, err := buildTLSConfig(certDir, "example.com")
	if err != nil {
		t.Fatalf("buildTLSConfig() error = %v", err)
	}
	cert, err := tlsConfig.GetCertificate(&tls.ClientHelloInfo{})
	if err != nil {
		t.Fatalf("GetCertificate() error = %v", err)
	}
	if got := certificateSerial(t, cert); got != 1 {
		t.Fatalf("serial = %d, want 1", got)
	}

	writeTestCertificate(t, certDir, "example.com", 2)
	future := time.Now().Add(time.Minute)
	for _, name := range []string{"fullchain.pem", "privkey.pem"} {
		path := filepath.Join(certDir, "example.com", name)
		if err := os.Chtimes(path, future, future); err != nil {
			t.Fatalf("Chtimes(%s): %v", path, err)
		}
	}

	cert, err = tlsConfig.GetCertificate(&tls.ClientHelloInfo{})
	if err != nil {
		t.Fatalf("GetCertificate() after renewal error = %v", err)
	}
	if got := certificateSerial(t, cert); got != 2 {
		t.Errorf("serial after renewal = %d, want 2", got)
	}
}

func TestEmailDomain(t *testing.T) {
	tests := []struct {
		address string
		want    string
		wantErr bool
	}{
		{address: "user@example.com", want: "example.com"},
		{address: "USER+tag@Example.COM", want: "example.com"},
		{address: "user@sub.example.com", want: "sub.example.com"},
		{address: "not-an-email", wantErr: true},
		{address: "user@", wantErr: true},
		{address: "@example.com", wantErr: true},
		{address: "user@example.com extra", wantErr: true},
	}
	for _, tt := range tests {
		got, err := emailDomain(tt.address)
		if tt.wantErr {
			if err == nil {
				t.Errorf("emailDomain(%q) error = nil, want error", tt.address)
			}
			continue
		}
		if err != nil {
			t.Errorf("emailDomain(%q) error = %v", tt.address, err)
			continue
		}
		if got != tt.want {
			t.Errorf("emailDomain(%q) = %q, want %q", tt.address, got, tt.want)
		}
	}
}

func TestDisplayNameUsesDomain(t *testing.T) {
	if got := displayName("example.com"); got != "example.com" {
		t.Errorf("displayName() = %q, want example.com", got)
	}
}

type thunderbirdConfig struct {
	XMLName  xml.Name `xml:"clientConfig"`
	Version  string   `xml:"version,attr"`
	Provider struct {
		ID             string `xml:"id,attr"`
		Domain         string `xml:"domain"`
		IncomingServer struct {
			Type       string `xml:"type,attr"`
			Hostname   string `xml:"hostname"`
			Port       int    `xml:"port"`
			SocketType string `xml:"socketType"`
			Username   string `xml:"username"`
		} `xml:"incomingServer"`
		OutgoingServer struct {
			Type       string `xml:"type,attr"`
			Hostname   string `xml:"hostname"`
			Port       int    `xml:"port"`
			SocketType string `xml:"socketType"`
			Username   string `xml:"username"`
		} `xml:"outgoingServer"`
	} `xml:"emailProvider"`
}

type outlookProtocol struct {
	Type       string `xml:"Type"`
	Server     string `xml:"Server"`
	Port       int    `xml:"Port"`
	LoginName  string `xml:"LoginName"`
	SSL        string `xml:"SSL"`
	Encryption string `xml:"Encryption"`
	Auth       string `xml:"AuthRequired"`
}

type outlookResponse struct {
	XMLName  xml.Name `xml:"Autodiscover"`
	Response struct {
		Account struct {
			AccountType string            `xml:"AccountType"`
			Action      string            `xml:"Action"`
			Protocols   []outlookProtocol `xml:"Protocol"`
		} `xml:"Account"`
		Error struct {
			ErrorCode string `xml:"ErrorCode"`
			Message   string `xml:"Message"`
		} `xml:"Error"`
	} `xml:"Response"`
}

func writeTestConfig(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "opengovmail.yaml")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func testConfig(domains ...string) *Config {
	cfg := &Config{}
	for _, domain := range domains {
		cfg.Domains = append(cfg.Domains, DomainConfig{Domain: domain})
	}
	return cfg
}

func mustHandler(t *testing.T, cfg *Config, mailHost string) http.Handler {
	t.Helper()
	h, err := newHandler(cfg, mailHost)
	if err != nil {
		t.Fatalf("newHandler() error = %v", err)
	}
	return h
}

func doRequest(h http.Handler, method, target, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, bytes.NewReader([]byte(body)))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decodeThunderbirdConfig(t *testing.T, body []byte) thunderbirdConfig {
	t.Helper()
	var cfg thunderbirdConfig
	if err := xml.Unmarshal(body, &cfg); err != nil {
		t.Fatalf("response is not valid XML: %v\n%s", err, body)
	}
	return cfg
}

func decodeOutlookResponse(t *testing.T, body []byte) outlookResponse {
	t.Helper()
	var resp outlookResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		t.Fatalf("response is not valid XML: %v\n%s", err, body)
	}
	return resp
}

func findProtocol(t *testing.T, resp outlookResponse, protocolType string) outlookProtocol {
	t.Helper()
	for _, protocol := range resp.Response.Account.Protocols {
		if protocol.Type == protocolType {
			return protocol
		}
	}
	t.Fatalf("protocol %s not found", protocolType)
	return outlookProtocol{}
}

func assertOutlookError(t *testing.T, body []byte, code string) {
	t.Helper()
	resp := decodeOutlookResponse(t, body)
	if resp.Response.Error.ErrorCode != code {
		t.Errorf("ErrorCode = %q, want %q\n%s", resp.Response.Error.ErrorCode, code, body)
	}
	if resp.Response.Error.Message == "" {
		t.Errorf("Error Message is empty\n%s", body)
	}
}

func outlookRequest(email string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<Autodiscover xmlns="http://schemas.microsoft.com/exchange/autodiscover/outlook/requestschema/2006">
  <Request>
    <EMailAddress>%s</EMailAddress>
    <AcceptableResponseSchema>http://schemas.microsoft.com/exchange/autodiscover/outlook/responseschema/2006a</AcceptableResponseSchema>
  </Request>
</Autodiscover>`, email)
}

func writeTestCertificate(t *testing.T, certDir, domain string, serial int64) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(serial),
		Subject:      pkix.Name{CommonName: domain},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{domain},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}

	lineage := filepath.Join(certDir, domain)
	if err := os.MkdirAll(lineage, 0o700); err != nil {
		t.Fatalf("create lineage: %v", err)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	if err := os.WriteFile(filepath.Join(lineage, "fullchain.pem"), certPEM, 0o600); err != nil {
		t.Fatalf("write certificate: %v", err)
	}
	if err := os.WriteFile(filepath.Join(lineage, "privkey.pem"), keyPEM, 0o600); err != nil {
		t.Fatalf("write key: %v", err)
	}
}

func certificateSerial(t *testing.T, cert *tls.Certificate) int64 {
	t.Helper()
	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatalf("parse certificate: %v", err)
	}
	return leaf.SerialNumber.Int64()
}
