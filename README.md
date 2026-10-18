# TarsPortkey

MCP integration for Portkey (AI gateway and observability).

## Tools

- `portkey_list_models` — List models via Portkey
- `portkey_chat` — Chat completion via Portkey
- `portkey_embeddings` — Generate embeddings via Portkey
- `portkey_feedback` — Submit feedback for a request
- `portkey_logs` — Query recent request logs
- `portkey_health` — Check Portkey gateway health
- `portkey_status` — Integration status

## Quick Start

```bash
go build -o tarsportkey

# Environment variables
TARS_MCP_ENDPOINT=https://api.example.com TARS_MCP_API_KEY=xxx ./tarsportkey

# Or JSON config file
./tarsportkey -config config.json
```

## Configuration

### JSON Config File

Create `config.json`:

```json
{
  "endpoint": "https://api.portkey.ai",
  "api_key": "your-api-key",
  "timeout": 30
}
```

Run with:

```bash
./tarsportkey -config config.json
```

### Environment Variables

- `TARS_MCP_ENDPOINT` — upstream API endpoint (fallback: `PORTKEY_ENDPOINT`)
- `TARS_MCP_API_KEY` — upstream API key (fallback: `PORTKEY_API_KEY`)
- `TARS_MCP_TIMEOUT` — request timeout in seconds (default: 30)
- `TARS_MCP_MODE` — `stdio` (default) or `http`

## Reliability

- **Automatic retry** with exponential backoff (3 attempts) on 5xx and network errors
- **Context timeout** propagation for all upstream calls
- **CORS** and security headers on HTTP endpoints

## License

MIT
