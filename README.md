# iRecall

iRecall is a local-first quote and note recall tool written in Go. It stores notes in SQLite, uses an OpenAI-compatible model to extract tags and recall keywords, and generates answers grounded only in the quotes it retrieves.

The project currently ships with:

- a Bubble Tea terminal client
- an HTTP web UI server
- a Wails-based desktop client
- a shared Go core for persistence, retrieval, import/export, and provider integration

## Features

- Create and edit free-form quotes and notes
- Auto-tag notes with an OpenAI-compatible chat-completions endpoint
- Search the local quote corpus with SQLite FTS5
- Filter weaker matches with a configurable relevance threshold
- Generate grounded answers from retrieved quotes
- Save Recall question/response pairs as quotes with generated tags
- Review Recall history and promote past sessions into quotes
- Export and import quotes through a versioned JSON share format
- Preserve author and source provenance on imported content
- Protect the web UI with a locally managed password and browser session
- Configure provider settings, retrieval settings, theme, web port, and mock-LLM debug mode
- Run isolated local instances with a custom data root
- Switch the active storage root from the TUI Settings page or with `-data-path`

## Getting Started

### Requirements

- Go
- a compatible OpenAI-style API endpoint for chat completions
- optional: Node.js and npm for the desktop frontend build

### Build the terminal client

```bash
make build
./bin/irecall
```

### Build the web UI server

```bash
make build-web
./bin/irecall-web
```

Useful flags:

```bash
./bin/irecall-web --debug
./bin/irecall-web -host 0.0.0.0
./bin/irecall-web -port 9527
./bin/irecall-web -data-path /tmp/irecall-web-dev
```

The persisted web-port default is `9527`. Port `95270` is not usable because TCP ports must be in the range `1..65535`.

Useful flags:

```bash
./bin/irecall --version
./bin/irecall --debug
./bin/irecall -data-path /tmp/irecall-dev
```

### First run

1. Save your display name in the startup prompt.
2. Open `Settings`.
3. Configure the provider host, port, HTTPS setting, API key if required, and model.
4. Optionally fetch available models from `/v1/models`.
5. If you are using the web UI, the first launch prompts for the web password in the terminal before the server starts listening. Use `Settings` to change it later.
6. Save the settings and start adding quotes.

## Provider Compatibility

iRecall expects an OpenAI-compatible API with:

- `POST /v1/chat/completions`
- `GET /v1/models`

Typical setups:

| Provider | Host | Port | HTTPS | API Key |
| --- | --- | --- | --- | --- |
| Ollama | `localhost` | `11434` | off | not required |
| LM Studio | `localhost` | `1234` | off | not required |
| Hosted OpenAI-compatible endpoint | provider host | provider port | usually on | provider-specific |

## Usage

The shipped clients currently expose four primary surfaces:

- `Recall`: ask questions and review the retrieved quotes used to answer them
- `History`: review saved recall sessions and reopen their reference quotes
- `Quotes`: browse, edit, delete, import, and export stored quotes
- `Settings`: configure provider, retrieval, and theme options

Notable workflows:

- add quotes from the TUI and refine drafts before saving
- save the current Recall question/response as a quote
- reopen a past History entry and save it as a quote later
- export selected quotes to a JSON payload
- import shared quotes back into another instance
- tune retrieval with `MaxResults` and `MinRelevance`

`MinRelevance` uses a normalized `0.0..1.0` scale:

- `0.0` disables filtering
- `0.3` to `0.7` is a good practical range
- `1.0` is very strict

## Data Storage

On Linux iRecall now consolidates defaults under the data directory so everything lives together by default. The effective defaults (when XDG vars are not set) are:

| Item | Default Path |
| --- | --- |
| All app files (data/config/state) | `~/.local/share/irecall/` |
| SQLite database | `~/.local/share/irecall/irecall.db` |
| Log file | `~/.local/share/irecall/irecall.log` |
| Persisted preferred root file | `~/.local/share/irecall/root-path` (used only if you choose a custom root)

Notes:

- Config and state now fall back to the data directory by default so the app uses a single per-user directory unless XDG_* environment variables override individual locations.
- To run an isolated instance with its own storage root, pass a root path. The root will contain `data/`, `config/`, and `state/` subdirectories:

```bash
./bin/irecall -data-path /path/to/instance
```

- You can also set the storage root from the TUI `Settings` page. When you save a new root, iRecall:
  - closes the current runtime (to avoid DB locks),
  - if the target root is empty, copies current `data/`, `config/`, and `state/` into `<root>/...` (copy, not move),
  - if the target already contains iRecall data, attaches to it without overwriting,
  - persists the chosen root so it is applied on future launches,
  - attempts to reopen the previous runtime on failure (partially-copied files are not automatically removed).

Leaving the setting blank returns to the platform default behavior (XDG paths or the consolidated data dir fallback described above).

Desktop and web currently show the active storage paths in `Settings`, but they do not yet expose a root-path editor.
## Project Layout

```text
iRecall/
├── cmd/irecall/      # terminal entry point
├── config/           # XDG path helpers
├── core/             # engine, models, DB layer, LLM client
│   ├── db/
│   └── llm/
├── app/              # Shared desktop/web application orchestration
├── desktop/          # Wails desktop runtime
├── frontend/         # Shared frontend assets and source
├── web/              # HTTP web UI runtime
├── docs/             # roadmap, specs, design docs, plans
├── tools/            # auxiliary tools such as Redmine export
├── tui/              # Bubble Tea application and pages
│   ├── pages/
│   └── styles/
├── Makefile
└── README.md
```

## Development

Common targets:

```bash
make build
make run
make test
make lint
make build-all
```

Desktop build:

```bash
make build-desktop
```

Web UI build:

```bash
make build-web
```

Frontend dependencies:

```bash
make frontend-install
make frontend-build
```

## Documentation

- [docs/PLAN.md](docs/PLAN.md): roadmap and planning index
- [docs/SPEC.md](docs/SPEC.md): technical specification
- [docs/KEYWORD_MATCHING.md](docs/KEYWORD_MATCHING.md): how keyword extraction and quote matching work
- [docs/schema.md](docs/schema.md): quote, share, and provenance schema guide
- [docs/UI_DESIGN.md](docs/UI_DESIGN.md): shared UI contract
- [docs/WAILS_DESKTOP.md](docs/WAILS_DESKTOP.md): desktop mapping
- [docs/QUOTES_SHARING_DESIGN.md](docs/QUOTES_SHARING_DESIGN.md): sharing model
- [docs/WEB_REST_API.md](docs/WEB_REST_API.md): current web REST endpoint reference
- [docs/WEB_API_TOKEN_AUTH.md](docs/WEB_API_TOKEN_AUTH.md): REST API bearer-token design
- [tools/redmine_export/README.md](tools/redmine_export/README.md): Redmine export tool usage
