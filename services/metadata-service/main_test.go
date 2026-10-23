package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"
)

// testMaxRequestBodyBytes mirrors the 1 MiB request body cap the server enforces.
const testMaxRequestBodyBytes = 1 << 20

func TestReceiveSuperPlatformResultAcceptsNormalPayload(t *testing.T) {
	payload := []byte(`{"status":"ok","data":{"note":"hello"},"timestamp":"2026-01-01T00:00:00Z"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/results", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	receiveSuperPlatformResultHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var ack map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &ack); err != nil {
		t.Fatalf("acknowledgment is not valid JSON: %v", err)
	}
	if ack["success"] != true {
		t.Errorf("success = %v, want true", ack["success"])
	}
	if _, ok := ack["message"]; !ok {
		t.Errorf("acknowledgment missing message field: %v", ack)
	}
}

func TestReceiveSuperPlatformResultRejectsOversizedBody(t *testing.T) {
	blob := strings.Repeat("a", testMaxRequestBodyBytes+1)
	payload := []byte(fmt.Sprintf(`{"status":"ok","data":{"blob":%q}}`, blob))
	if len(payload) <= testMaxRequestBodyBytes {
		t.Fatalf("test payload is %d bytes, want > %d", len(payload), testMaxRequestBodyBytes)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/results", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	receiveSuperPlatformResultHandler(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestReceiveSuperPlatformResultLogsPayloadSizeNotContent(t *testing.T) {
	const secret = "SUPER_SECRET_PAYLOAD_VALUE"
	payload := []byte(fmt.Sprintf(`{"status":"ok","data":{"note":%q},"timestamp":"2026-01-01T00:00:00Z"}`, secret))

	logs := captureLogs(t, func() {
		req := httptest.NewRequest(http.MethodPost, "/api/results", bytes.NewReader(payload))
		rec := httptest.NewRecorder()
		receiveSuperPlatformResultHandler(rec, req)
	})

	if strings.Contains(logs, secret) {
		t.Errorf("log output leaked payload content: %s", logs)
	}
	if !regexp.MustCompile(`Data received: \d+ bytes`).MatchString(logs) {
		t.Errorf("log output does not report payload size in bytes: %s", logs)
	}
}

func TestNewServerConfiguresTimeoutsAndLimits(t *testing.T) {
	srv := newServer("0.0.0.0:8888", http.NewServeMux())

	if srv.Addr != "0.0.0.0:8888" {
		t.Errorf("Addr = %q, want %q", srv.Addr, "0.0.0.0:8888")
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

func TestServerServesHealthAndRejectsOversizedPost(t *testing.T) {
	router := http.NewServeMux()
	router.HandleFunc("/health", healthHandler)
	router.HandleFunc("/api/results", apiKeyAuthMiddleware(receiveSuperPlatformResultHandler))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	srv := newServer(ln.Addr().String(), router)

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(ln) }()
	defer srv.Close()

	baseURL := "http://" + ln.Addr().String()

	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /health status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	blob := strings.Repeat("a", testMaxRequestBodyBytes+1)
	resp, err = http.Post(baseURL+"/api/results", "application/json",
		bytes.NewReader([]byte(fmt.Sprintf(`{"status":"ok","data":{"blob":%q}}`, blob))))
	if err != nil {
		t.Fatalf("POST /api/results failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("oversized POST status = %d, want %d", resp.StatusCode, http.StatusRequestEntityTooLarge)
	}
}

// captureLogs redirects the standard logger for the duration of fn and returns
// everything written.
func captureLogs(t *testing.T, fn func()) string {
	t.Helper()
	var buf bytes.Buffer
	prev := log.Writer()
	flags := log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer func() {
		log.SetOutput(prev)
		log.SetFlags(flags)
	}()
	fn()
	return buf.String()
}
