# iRecall

A personal knowledge recall system powered by AI. Capture anything — notes, quotes, learnings — and retrieve them conversationally. The LLM extracts tags on ingestion and synthesizes answers from your own knowledge base when you ask questions.

## How It Works

1. **Add a quote** — paste any note or text snippet. The LLM automatically extracts keyword tags.
2. **Ask a question** — type anything in natural language. The LLM extracts keywords, finds matching quotes via full-text search, and synthesizes a response grounded in your own notes.
3. **See your sources** — matched reference quotes are always shown alongside the response so you know exactly where the answer came from.

## Features

- AI-driven tagging and retrieval via any OpenAI-compatible provider (Ollama, LM Studio, OpenAI, etc.)
- Full-text search with BM25 ranking (SQLite FTS5)
- Streaming LLM responses
- Clean separation of engine and UI — same core works with TUI, web, or native UI
- TUI interface built with Bubbletea

## Quick Start

```bash
# Build
go build -o irecall ./cmd/irecall

# Run
./irecall
```

On first launch, go to the **Settings** page (`Tab`) and configure your LLM provider.

## UI Overview

```
Main Page (Recall)         Settings Page
┌──────────────────┐       ┌──────────────────┐
│ [Question input] │       │ Host / Port      │
│ ─── Response ──  │  Tab  │ API Key          │
│ Streamed answer  │ ───── │ [Fetch Models]   │
│ ─── References ─ │       │ Model selector   │
│ [1] quote...     │       │ Search tuning    │
│ [2] quote...     │       └──────────────────┘
└──────────────────┘
Ctrl+N → Add Quote popup
```

## Keyboard Shortcuts

| Key | Action |
|---|---|
| `Enter` | Submit question |
| `Ctrl+N` | Open Add Quote popup |
| `Tab` | Switch between Recall / Settings |
| `Ctrl+S` | Save (in popup or settings) |
| `Esc` | Close popup / cancel |
| `Q` / `Ctrl+C` | Quit |

## LLM Provider Setup

iRecall works with any OpenAI-compatible API endpoint.

| Provider | Host | Port | Notes |
|---|---|---|---|
| Ollama | localhost | 11434 | No API key needed |
| LM Studio | localhost | 1234 | No API key needed |
| OpenAI | api.openai.com | 443 | Enable HTTPS + API key |
| Any compatible | your host | your port | |

## Data Storage

Follows XDG Base Directory spec:

| Item | Default Path |
|---|---|
| Database | `~/.local/share/irecall/irecall.db` |
| Config | `~/.config/irecall/config.json` |
| Logs | `~/.local/state/irecall/irecall.log` |

## Project Structure

```
iRecall/
├── core/          # Engine — pure business logic, no UI
│   ├── db/        # SQLite + FTS5 layer
│   ├── llm/       # OpenAI-compatible client
│   ├── engine.go  # Public engine API
│   └── models.go  # Shared data types
├── tui/           # Bubbletea TUI (depends only on core)
│   ├── pages/     # Recall, Settings, AddQuote
│   └── styles/    # Lipgloss theme
├── cmd/irecall/   # Entry point
└── config/        # XDG config helpers
```

## Roadmap

- v1.0 — TUI with recall, add quote, settings
- v2.0 — Web UI over the same core engine
- v3.0 — Native OS UI

## License

MIT
