# iRecall MCP Bridge

This document tracks the initial in-repo MCP bridge for iRecall.

## Goal

Expose a small set of iRecall capabilities to agent clients such as OpenClaw through MCP, while reusing the existing authenticated web REST API.

## Current shape

The first scaffold adds:

- `cmd/irecall-mcp/` for the `irecall-mcp` stdio binary
- `mcp/` for bridge configuration and server wiring
- `mcp/irecallapi/` for the REST client wrapper
- `mcp/tools/` for MCP tool registration

Current tool set:

- `irecall_health`
- `irecall_recall`
- `irecall_list_quotes`
- `irecall_add_quote`
- `irecall_save_recall_as_quote`

History-oriented MCP tools are intentionally deferred to a later pass.

## Runtime model

```text
OpenClaw <-> MCP stdio <-> irecall-mcp <-> HTTP REST <-> iRecall web server
```

## Required environment

- `IRECALL_API_TOKEN` — required bearer token used for authenticated REST calls
- `IRECALL_BASE_URL` — optional base URL, defaults to `http://127.0.0.1:9527`

## Current limitations

- token provisioning and rotation are still separate work
- stdio is the only bridge transport in this first shell
- tool responses are currently returned as JSON text payloads

## Next likely steps

1. Add token provisioning commands for local operators
2. Expand the MCP tool surface to history/update/delete operations
3. Refine response formatting once the bridge contract settles
