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
TARS_MCP_ENDPOINT=https://api.example.com TARS_MCP_API_KEY=xxx ./tarsportkey
```

## Environment Variables

- `TARS_MCP_ENDPOINT` — upstream API endpoint
- `TARS_MCP_API_KEY` — upstream API key
- `TARS_MCP_MODE` — `stdio` (default) or `http`

## License

MIT
