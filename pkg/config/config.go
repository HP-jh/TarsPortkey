package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Endpoint string `json:"endpoint"`
	APIKey   string `json:"api_key"`
	Timeout  int    `json:"timeout"`
}

func Load(path string) (*Config, error) {
	if path != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("读取配置文件失败：%w", err)
		}
		var cfg Config
		if err := json.Unmarshal(b, &cfg); err != nil {
			return nil, fmt.Errorf("解析配置文件失败：%w", err)
		}
		if cfg.Endpoint == "" {
			return nil, fmt.Errorf("配置文件中 endpoint 不能为空")
		}
		if cfg.Timeout <= 0 {
			cfg.Timeout = 30
		}
		cfg.Endpoint = strings.TrimRight(cfg.Endpoint, "/")
		return &cfg, nil
	}
	return loadFromEnv()
}

func loadFromEnv() (*Config, error) {
	endpoint := os.Getenv("TARS_MCP_ENDPOINT")
	if endpoint == "" {
		endpoint = os.Getenv("PORTKEY_ENDPOINT")
	}
	apiKey := os.Getenv("TARS_MCP_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("PORTKEY_API_KEY")
	}
	timeout := 30
	if v := os.Getenv("TARS_MCP_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			timeout = n
		}
	}
	if endpoint == "" {
		return nil, fmt.Errorf("环境变量 TARS_MCP_ENDPOINT 或 PORTKEY_ENDPOINT 未设置")
	}
	return &Config{
		Endpoint: strings.TrimRight(endpoint, "/"),
		APIKey:   apiKey,
		Timeout:  timeout,
	}, nil
}
