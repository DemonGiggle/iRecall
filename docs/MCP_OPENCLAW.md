# iRecall MCP Bridge

This document describes the in-repo MCP bridge for iRecall and the local operator flow for connecting it to agent clients such as OpenClaw.

## Goal

Expose a small set of iRecall capabilities through MCP while reusing the existing authenticated web REST API.

The intended runtime model is:

```text
OpenClaw <-> MCP stdio <-> irecall-mcp <-> HTTP REST <-> iRecall web server
```

This keeps the integration decoupled:

- iRecall continues to own the local web server and REST API
- the MCP bridge is a small adapter process
- the bearer token stays in the bridge process environment or a protected local credential file

## Components

- `cmd/irecall-mcp/` — stdio MCP bridge binary
- `mcp/` — bridge configuration and server wiring
- `mcp/irecallapi/` — authenticated REST client wrapper
- `mcp/tools/` — MCP tool registration
- `web auth ...` — local token provisioning commands on the web binary

## Current MCP tool set

- `irecall_health` — checks the local iRecall REST API and bearer-token auth, returning only minimal status
- `irecall_recall` — calls the recall flow
- `irecall_list_quotes` — lists stored quotes with `limit`/`offset` pagination
- `irecall_add_quote` — stores a free-form note or quote
- `irecall_save_recall_as_quote` — persists a recall question/response pair as a quote
- `irecall_update_quote` — updates an existing quote by ID
- `irecall_delete_quotes` — deletes one or more quotes by ID
- `irecall_list_history` — lists saved recall-history summaries
- `irecall_get_history` — gets one saved recall-history entry with referenced quotes
- `irecall_delete_history` — deletes one or more recall-history entries by ID

## Required environment for `irecall-mcp`

- `IRECALL_API_TOKEN` — required bearer token used for authenticated REST calls
- `IRECALL_BASE_URL` — optional base URL, defaults to `http://127.0.0.1:9527`

Example:

```bash
IRECALL_BASE_URL=http://127.0.0.1:9527 \
IRECALL_API_TOKEN="$(cat ~/.config/irecall/mcp-api-token)" \
./bin/irecall-mcp
```

## Local token provisioning

Token provisioning is intentionally a local operator action. It requires the existing web password and does not require browser automation.

Build the web binary first:

```bash
make build-web
```

Issue the first MCP token and write it to a protected file:

```bash
printf '%s\n' 'your-web-password' | \
  ./bin/irecall-web auth issue-token \
    --password-stdin \
    --write-token-file ~/.config/irecall/mcp-api-token
```

Rotate the token with the same flow:

```bash
printf '%s\n' 'your-web-password' | \
  ./bin/irecall-web auth rotate-token \
    --password-stdin \
    --write-token-file ~/.config/irecall/mcp-api-token
```

Revoke the token:

```bash
printf '%s\n' 'your-web-password' | \
  ./bin/irecall-web auth revoke-token --password-stdin
```

Check whether a token is configured without printing the token:

```bash
./bin/irecall-web auth token-status
```

For isolated instances, pass `--data-path` to each auth command:

```bash
printf '%s\n' 'your-web-password' | \
  ./bin/irecall-web auth issue-token \
    --data-path /path/to/irecall-instance \
    --password-stdin \
    --write-token-file /path/to/irecall-token
```

The token file is written with mode `0600`. Command output prints only the destination path and token prefix when `--write-token-file` is used.

## OpenClaw setup outline

1. Start `irecall-web` locally, preferably bound to `127.0.0.1`.
2. Configure the web password through the normal first-run flow if it is not already configured.
3. Use `irecall-web auth issue-token --password-stdin --write-token-file ...` to create the MCP bearer token.
4. Configure the MCP launcher to run `irecall-mcp` with `IRECALL_BASE_URL` and `IRECALL_API_TOKEN` loaded from the protected local credential file or systemd credential.
5. Verify the connection by calling `irecall_health`.

## OpenClaw MCP launcher example

The exact OpenClaw MCP configuration shape may vary by OpenClaw version, but the launcher should be equivalent to:

```json
{
  "command": "/path/to/bin/irecall-mcp",
  "env": {
    "IRECALL_BASE_URL": "http://127.0.0.1:9527",
    "IRECALL_API_TOKEN": "<contents of ~/.config/irecall/mcp-api-token>"
  }
}
```

Prefer loading `IRECALL_API_TOKEN` from a protected local file or service credential rather than hard-coding it into a shared config file.

## Current limitations

- stdio is the only bridge transport in this implementation
- tool responses are currently returned as JSON text payloads
- health intentionally omits bootstrap details such as UI pages, local paths, settings, and docs
- concrete OpenClaw config may need adjustment for the installed OpenClaw MCP config schema
- the bridge assumes the iRecall web server is already running

## Next likely steps

1. Run a real local operator bootstrap: start web, issue token, launch MCP, call `irecall_health`
2. Refine response formatting once the bridge contract settles
3. Add a version-pinned OpenClaw config snippet after confirming the target config schema
