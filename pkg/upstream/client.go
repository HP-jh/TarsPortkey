package upstream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	endpoint string
	apiKey   string
	client   *http.Client
	maxRetry int
}

func NewClient(endpoint, apiKey string, timeout int) *Client {
	if endpoint == "" {
		endpoint = "https://api.example.com"
	}
	if timeout <= 0 {
		timeout = 30
	}
	return &Client{
		endpoint: strings.TrimRight(endpoint, "/"),
		apiKey:   apiKey,
		client:   &http.Client{Timeout: time.Duration(timeout) * time.Second},
		maxRetry: 3,
	}
}

func (c *Client) Health(ctx context.Context) map[string]interface{} {
	res, err := c.InvokeRaw(ctx, "GET", "/health", nil)
	if err != nil {
		return map[string]interface{}{"healthy": false, "error": err.Error()}
	}
	res["healthy"] = true
	return res
}

func (c *Client) ListModels(ctx context.Context) ([]interface{}, error) {
	res, err := c.InvokeRaw(ctx, "GET", "/v1/models", nil)
	if err != nil {
		return nil, err
	}
	models, _ := res["data"].([]interface{})
	if models == nil {
		models = []interface{}{}
	}
	return models, nil
}

func (c *Client) InvokeRaw(ctx context.Context, method, path string, body interface{}) (map[string]interface{}, error) {
	var lastErr error
	for attempt := 0; attempt <= c.maxRetry; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(100*(1<<(attempt-1))) * time.Millisecond
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
		res, err := c.invokeOnce(ctx, method, path, body)
		if err == nil {
			return res, nil
		}
		lastErr = err
		if !isRetryable(err) {
			break
		}
	}
	return nil, lastErr
}

func (c *Client) invokeOnce(ctx context.Context, method, path string, body interface{}) (map[string]interface{}, error) {
	var bodyR io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bodyR = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+path, bodyR)
	if err != nil {
		return nil, err
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(b))
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return map[string]interface{}{"raw": "non-JSON response"}, nil
	}
	return result, nil
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	if strings.Contains(s, "HTTP 5") {
		return true
	}
	if _, ok := err.(net.Error); ok {
		return true
	}
	if strings.Contains(s, "connection refused") || strings.Contains(s, "no such host") || strings.Contains(s, "timeout") || strings.Contains(s, "temporary failure") {
		return true
	}
	return false
}
