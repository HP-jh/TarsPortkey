package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Setenv("TARS_MCP_ENDPOINT", "https://example.com")
	os.Setenv("TARS_MCP_API_KEY", "test-key")
	defer os.Unsetenv("TARS_MCP_ENDPOINT")
	defer os.Unsetenv("TARS_MCP_API_KEY")
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if cfg.Endpoint != "https://example.com" {
		t.Fatalf("expected endpoint, got %s", cfg.Endpoint)
	}
	if cfg.APIKey != "test-key" {
		t.Fatalf("expected key, got %s", cfg.APIKey)
	}
	if cfg.Timeout != 30 {
		t.Fatalf("expected timeout 30, got %d", cfg.Timeout)
	}
}

func TestLoad_MissingEndpoint(t *testing.T) {
	os.Unsetenv("TARS_MCP_ENDPOINT")
	os.Unsetenv("PORTKEY_ENDPOINT")
	_, err := Load("")
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestLoad_FromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{"endpoint":"https://cfg.example.com","api_key":"cfg-key","timeout":45}`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if cfg.Endpoint != "https://cfg.example.com" {
		t.Fatalf("expected cfg endpoint, got %s", cfg.Endpoint)
	}
	if cfg.APIKey != "cfg-key" {
		t.Fatalf("expected cfg key, got %s", cfg.APIKey)
	}
	if cfg.Timeout != 45 {
		t.Fatalf("expected timeout 45, got %d", cfg.Timeout)
	}
}

func TestLoad_FromFile_MissingEndpoint(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	content := `{"api_key":"key"}`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing endpoint in file")
	}
}

func TestLoad_FromFile_DefaultTimeout(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{"endpoint":"https://example.com"}`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if cfg.Timeout != 30 {
		t.Fatalf("expected default timeout 30, got %d", cfg.Timeout)
	}
}
