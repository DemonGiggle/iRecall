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
- Include a Wails-oriented desktop scaffold that reuses the same Go core and UI contract

## Interface

The current TUI has three pages plus a modal overlay:

- `Recall`: ask questions, see extracted keywords, streamed answers, and matched reference notes
- `Quotes`: browse every stored quote and its tags
- `Settings`: configure host, port, HTTPS, API key, model, and search limits
- `Add Quote` modal: open from the Recall page with `ctrl+n`

Global navigation:

| Key | Action |
| --- | --- |
| `Tab` | Cycle `Recall -> Quotes -> Settings -> Recall` |
| `Shift+Tab` | Cycle `Recall <- Quotes <- Settings <- Recall` |
| `ctrl+c` | Quit |

Recall page:

| Key | Action |
| --- | --- |
| `enter` | Run recall workflow for the current question |
| `ctrl+n` | Open the Add Quote modal |
| `ctrl+j` | Jump focus between the input and the Reference Quotes panel |

When the Reference Quotes panel is focused:

| Key | Action |
| --- | --- |
| `up` / `down` | Move between retrieved quotes |
| `x` | Select or unselect the current quote |
| `e` | Edit the current quote |
| `d` | Delete the current selection |
| `s` | Export the current or selected quotes in the Share Quotes modal |

Quotes page:

| Key | Action |
| --- | --- |
| `ctrl+n` | Open the Add Quote modal |
| `i` | Open the Import Quotes modal |
| `r` | Reload stored quotes |
| `up` / `down` | Scroll |
| `x` | Select or unselect the current quote |
| `e` | Edit the current quote |
| `d` | Delete the current selection |
| `s` | Export the current or selected quotes in the Share Quotes modal |
| `pgup` / `pgdn` | Page scroll |

Settings page:

| Key | Action |
| --- | --- |
| `up` / `down` | Move focus between fields |
| `space` | Toggle HTTPS when focused |
| `enter` | Fetch models when the button is focused |
| `left` / `right` | Cycle fetched models when the model field is focused |
| `ctrl+s` | Save settings |

Add Quote modal:

| Key | Action |
| --- | --- |
| `ctrl+r` | Ask the LLM to refine the current draft, then preview the suggestion |
| `ctrl+s` | Save the note |
| `esc` | Close the modal if no save is in progress |

Share Quotes modal:

| Key | Action |
| --- | --- |
| `ctrl+s` / `enter` | Save the exported JSON payload to the file path in the modal |
| `esc` | Close the modal |

Import Quotes modal:

| Key | Action |
| --- | --- |
| `ctrl+s` / `enter` | Import the JSON payload from the file path in the modal |
| `esc` | Close the modal |

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

1. Open `Settings` with `tab`.
2. Enter the provider host, port, HTTPS preference, API key if needed, and a model name.
3. Optionally use `Fetch Models` to populate the model selector from `/v1/models`.
4. Save with `ctrl+s`.

After that:

1. Go back to `Recall`.
2. Press `ctrl+n` to add notes.
3. Ask questions with `enter`.

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
- Use `irecall -data-path /path/to/instance` to run an isolated local instance. That root will contain `data/irecall.db`, `config/`, and `state/irecall.log`.

## Project Layout

```text
iRecall/
├── cmd/irecall/      # CLI entry point and TUI startup
├── config/           # XDG directory helpers
├── core/             # Engine, data models, DB layer, LLM client
│   ├── db/
│   └── llm/
├── desktop/          # Wails-oriented desktop scaffold
│   ├── backend/
│   └── frontend/
├── docs/             # Shared UI and platform design references
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
- `docs/UI_DESIGN.md` is the shared UI contract for future clients.
- `desktop/` contains the current Wails-oriented desktop scaffold and backend service layer.
