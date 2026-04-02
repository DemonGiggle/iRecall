# iRecall

iRecall is a terminal-first personal knowledge recall tool written in Go. It stores your notes in SQLite, asks an OpenAI-compatible model to extract tags and search keywords, and then answers questions strictly from the notes it retrieved.

## Current Capabilities

- Add free-form notes through a modal composer in the TUI
- Auto-tag notes with an OpenAI-compatible chat-completions endpoint
- Search notes with SQLite FTS5 and BM25 ranking
- Generate grounded answers with streamed LLM output
- Browse all stored notes on a dedicated Quotes page
- Configure provider connection details and search settings from the Settings page
- Persist notes and settings locally with XDG-style data/state directories

## Interface

The current TUI has three pages plus a modal overlay:

- `Recall`: ask questions, see extracted keywords, streamed answers, and matched reference notes
- `Quotes`: browse every stored quote and its tags
- `Settings`: configure host, port, HTTPS, API key, model, and search limits
- `Add Quote` modal: open from the Recall page with `Ctrl+N`

Global navigation:

| Key | Action |
| --- | --- |
| `Tab` | Cycle `Recall -> Quotes -> Settings -> Recall` |
| `Shift+Tab` | Cycle `Recall <- Quotes <- Settings <- Recall` |
| `Ctrl+C` | Quit |

Recall page:

| Key | Action |
| --- | --- |
| `Enter` | Run recall workflow for the current question |
| `Ctrl+N` | Open the Add Quote modal |

Quotes page:

| Key | Action |
| --- | --- |
| `R` | Reload stored quotes |
| `Up` / `Down` | Scroll |
| `PgUp` / `PgDn` | Page scroll |

Settings page:

| Key | Action |
| --- | --- |
| `Up` / `Down` | Move focus between fields |
| `Space` | Toggle HTTPS when focused |
| `Enter` | Fetch models when the button is focused |
| `Left` / `Right` | Cycle fetched models when the model field is focused |
| `Ctrl+S` | Save settings |

Add Quote modal:

| Key | Action |
| --- | --- |
| `Ctrl+R` | Ask the LLM to refine the current draft, then preview the suggestion |
| `Ctrl+S` | Save the note |
| `Esc` | Close the modal if no save is in progress |

## Quick Start

```bash
make build
./bin/irecall
```

Useful flags:

```bash
./bin/irecall --version
./bin/irecall --debug
```

On first launch:

1. Open `Settings` with `Tab`.
2. Enter the provider host, port, HTTPS preference, API key if needed, and a model name.
3. Optionally use `Fetch Models` to populate the model selector from `/v1/models`.
4. Save with `Ctrl+S`.

After that:

1. Go back to `Recall`.
2. Press `Ctrl+N` to add notes.
3. Ask questions with `Enter`.

## Provider Compatibility

iRecall expects an OpenAI-compatible API with:

- `POST /v1/chat/completions`
- `GET /v1/models`

Typical setups:

| Provider | Host | Port | HTTPS | API Key |
| --- | --- | --- | --- | --- |
| Ollama | `localhost` | `11434` | off | not required |
| LM Studio | `localhost` | `1234` | off | not required |
| OpenAI-compatible hosted endpoint | provider host | provider port | usually on | provider-specific |

## Data and Logs

iRecall follows XDG-style directories and stores everything locally:

| Item | Default Path |
| --- | --- |
| SQLite database | `~/.local/share/irecall/irecall.db` |
| Log file | `~/.local/state/irecall/irecall.log` |
| Reserved config directory | `~/.config/irecall/` |

Notes:

- Settings are currently stored in the SQLite `settings` table, not in a JSON config file.
- The config directory is created up front but is not yet used for persisted configuration.

## Project Layout

```text
iRecall/
├── cmd/irecall/      # CLI entry point and TUI startup
├── config/           # XDG directory helpers
├── core/             # Engine, data models, DB layer, LLM client
│   ├── db/
│   └── llm/
├── tui/              # Bubble Tea application and pages
│   ├── pages/
│   └── styles/
├── README.md
├── SPEC.md
├── PLAN.md
└── Makefile
```

## Development

```bash
make build
make run
make test
make lint
make build-all
```

Current implementation notes:

- Search uses `MaxResults` today.
- `MinRelevance` is captured and persisted in settings, but it is not yet applied to the FTS query.
- The app is TUI-only at the moment; there is no web or native UI in the repository.
