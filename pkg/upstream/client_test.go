package upstream

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c := NewClient("https://example.com", "key", 10)
	if c.endpoint != "https://example.com" {
		t.Fatal("endpoint mismatch")
	}
	if c.apiKey != "key" {
		t.Fatal("key mismatch")
	}
}

func TestHealth_Up(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := NewClient(srv.URL, "", 5)
	res := c.Health(context.Background())
	if res["healthy"] != true {
		t.Fatalf("expected healthy, got %v", res)
	}
}

func TestHealth_Down(t *testing.T) {
	c := NewClient("http://localhost:1", "", 1)
	res := c.Health(context.Background())
	if res["healthy"] != false {
		t.Fatalf("expected not healthy, got %v", res)
	}
}

func TestListModels(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/models", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{map[string]interface{}{"id": "model-1"}}})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := NewClient(srv.URL, "", 5)
	models, err := c.ListModels(context.Background())
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(models))
	}
}

func TestInvokeRaw(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := NewClient(srv.URL, "key", 5)
	res, err := c.InvokeRaw(context.Background(), "GET", "/test", nil)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res["ok"] != true {
		t.Fatalf("expected ok, got %v", res)
	}
}

func TestInvokeRaw_Retry_SuccessAfterFailure(t *testing.T) {
	var count atomic.Int32
	mux := http.NewServeMux()
	mux.HandleFunc("/flaky", func(w http.ResponseWriter, r *http.Request) {
		if count.Add(1) <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "temporarily unavailable"})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := NewClient(srv.URL, "key", 5)
	start := time.Now()
	res, err := c.InvokeRaw(context.Background(), "GET", "/flaky", nil)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("expected success after retry, got err: %v", err)
	}
	if res["ok"] != true {
		t.Fatalf("expected ok, got %v", res)
	}
	if count.Load() != 3 {
		t.Fatalf("expected 3 attempts, got %d", count.Load())
	}
	if elapsed < 250*time.Millisecond {
		t.Fatalf("expected backoff delay, got %v", elapsed)
	}
}

func TestInvokeRaw_Retry_Exhausted(t *testing.T) {
	var count atomic.Int32
	mux := http.NewServeMux()
	mux.HandleFunc("/alwaysfail", func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "down"})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := NewClient(srv.URL, "key", 5)
	_, err := c.InvokeRaw(context.Background(), "GET", "/alwaysfail", nil)
	if err == nil {
		t.Fatal("expected error after retries exhausted")
	}
	if count.Load() != 4 { // initial + 3 retries
		t.Fatalf("expected 4 attempts, got %d", count.Load())
	}
}

func TestInvokeRaw_NoRetry_4xx(t *testing.T) {
	var count atomic.Int32
	mux := http.NewServeMux()
	mux.HandleFunc("/badrequest", func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "bad"})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := NewClient(srv.URL, "key", 5)
	_, err := c.InvokeRaw(context.Background(), "GET", "/badrequest", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if count.Load() != 1 {
		t.Fatalf("expected 1 attempt for 4xx, got %d", count.Load())
	}
}
