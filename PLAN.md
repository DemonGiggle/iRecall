# iRecall — Implementation Plan

## Goal

Build an AI-driven personal knowledge recall system. The user captures notes/quotes;
the system tags them automatically via LLM and retrieves relevant ones when the user
asks questions, synthesizing a grounded response from their own knowledge base.

## Language & Runtime

**Go 1.22+**

Rationale:
- `bubbletea` + `lipgloss` + `bubbles` is the premier TUI stack in any language
- Compiles to a single static binary — zero install friction
- `modernc.org/sqlite` is pure Go (no CGO), making cross-compilation trivial
- Native concurrency suits streaming LLM responses cleanly
- The engine/UI separation maps naturally to Go packages

## Architecture Principle

> The `core/` package must never import anything from `tui/`. All business logic
> lives in `core/`. The `tui/` package is a thin presentation layer that calls
> the engine API and renders results.

This contract ensures future UIs (web, native) are drop-in consumers of the
same engine with no refactoring.

---

## Phases

### Phase 0 — Repository Bootstrap
- [x] Write PLAN.md, SPEC.md, README.md
- [ ] `go mod init github.com/gigol/irecall`
- [ ] Directory scaffold (see SPEC.md §Project Structure)
- [ ] Pre-commit: `go vet`, `go fmt` checks

### Phase 1 — Core: Database Layer
**Target:** `core/db/`

- [ ] SQLite connection pool with WAL mode enabled
- [ ] Migration runner (`schema_version` table)
- [ ] Schema v1: `quotes`, `tags`, `quote_tags`, `quotes_fts`, `settings`
- [ ] FTS5 sync triggers (insert / update / delete)
- [ ] `Store` struct with methods:
  - `InsertQuote(content string, tags []string) (int64, error)`
  - `UpdateQuoteFTS(id int64, tags []string) error`
  - `SearchQuotes(ftsQuery string, limit int) ([]Quote, error)`
  - `ListQuotes() ([]Quote, error)`
  - `DeleteQuote(id int64) error`
  - `GetSetting(key string) (string, error)`
  - `SetSetting(key, value string) error`
- [ ] Unit tests with in-memory SQLite (`:memory:`)

### Phase 2 — Core: LLM Client
**Target:** `core/llm/`

- [ ] `ProviderConfig` struct (Host, Port, HTTPS, APIKey, Model)
- [ ] `BaseURL()` helper
- [ ] HTTP client wrapping OpenAI `/v1/chat/completions`
- [ ] Support streaming (`stream: true`) via SSE parsing
- [ ] `FetchModels()` — calls `/v1/models`, returns sorted model IDs
- [ ] `TestConnection()` — lightweight ping / model list check
- [ ] Timeout and context cancellation propagation
- [ ] Unit tests with `httptest.Server` mock

### Phase 3 — Core: Engine API
**Target:** `core/engine.go`, `core/models.go`

- [ ] `Engine` struct (wraps `db.Store` + `llm.Client` + config)
- [ ] `AddQuote(ctx, content)` — store → extract tags → update FTS
- [ ] `ExtractKeywords(ctx, question)` — LLM call, returns `[]string`
- [ ] `SearchQuotes(ctx, keywords)` — build FTS query, ranked results
- [ ] `GenerateResponse(ctx, question, quotes, tokenCh)` — streaming synthesis
- [ ] `FetchModels(ctx, config)` — delegates to llm.Client
- [ ] `SaveSettings / LoadSettings`
- [ ] Integration test: full add-then-recall round trip

### Phase 4 — TUI: Foundation
**Target:** `tui/app.go`, `tui/styles/`

- [ ] Root Bubbletea model with page routing (`currentPage` enum)
- [ ] Global key bindings: Tab (switch page), Q/Ctrl+C (quit)
- [ ] Lipgloss theme: colors, borders, widths (adaptive to terminal size)
- [ ] `WindowSizeMsg` handler for responsive layout
- [ ] Page interface: `Update(msg) (Page, tea.Cmd)`, `View() string`

### Phase 5 — TUI: Recall Page
**Target:** `tui/pages/recall.go`

- [ ] Question input (single-line `textinput`)
- [ ] Response viewport (scrollable, streaming-aware)
- [ ] Reference quotes panel (scrollable list below response)
- [ ] `Enter` → trigger recall workflow as `tea.Cmd`
- [ ] Spinner shown during LLM calls
- [ ] Stream tokens via `tea.Msg` channel bridge
- [ ] `Ctrl+N` → emit `OpenAddQuoteMsg` to app router

### Phase 6 — TUI: Add Quote Overlay
**Target:** `tui/pages/addquote.go`

- [ ] Modal overlay rendered on top of recall page
- [ ] Multi-line `textarea` input (`bubbles/textarea`)
- [ ] `Ctrl+S` → call `Engine.AddQuote`, show spinner, close on success
- [ ] `Esc` → cancel, return to recall page
- [ ] Error display inline if add fails

### Phase 7 — TUI: Settings Page
**Target:** `tui/pages/settings.go`

- [ ] Form fields: Host, Port, HTTPS toggle, API Key (masked), Model dropdown
- [ ] `[Fetch Models]` button → `Engine.FetchModels` → populate dropdown
- [ ] Search tuning: max reference quotes (1–20), min relevance score
- [ ] `Ctrl+S` → persist via `Engine.SaveSettings`
- [ ] Load current settings on page init
- [ ] Inline success/error feedback

### Phase 8 — Polish & Packaging
- [ ] XDG path resolution with `os.UserDataDir` fallbacks
- [ ] Structured logging to file (not stdout — would corrupt TUI)
- [ ] Graceful shutdown: cancel in-flight LLM requests on quit
- [ ] Error boundary: never crash TUI; show error in status bar
- [ ] Build flags: version, commit hash embedded at link time
- [ ] `Makefile` with `build`, `test`, `lint`, `install` targets
- [ ] GitHub Actions CI (build + test on Linux/macOS/Windows)

---

## Key Dependencies

```
github.com/charmbracelet/bubbletea    v1.x   # TUI framework
github.com/charmbracelet/bubbles      v0.x   # textinput, textarea, viewport, list, spinner
github.com/charmbracelet/lipgloss     v1.x   # styling
modernc.org/sqlite                    v1.x   # pure-Go SQLite (no CGO)
github.com/sashabaranov/go-openai     v1.x   # OpenAI-compatible client
```

No other runtime dependencies. All are vendorable.

---

## Testing Strategy

| Layer | Approach |
|---|---|
| `core/db` | In-memory SQLite (`:memory:`), table-driven |
| `core/llm` | `httptest.Server` mock, test streaming |
| `core/engine` | Integration tests, real in-memory DB + mock LLM |
| `tui/` | Manual + snapshot testing via Bubbletea test helpers |

Run all tests: `go test ./...`

---

## Decision Log

| Decision | Rationale |
|---|---|
| Go over Python | Single binary, better TUI ecosystem, no packaging mess |
| SQLite over Postgres | Embedded, zero-ops, FTS5 built-in, sufficient for personal use |
| FTS5 over manual tag matching | BM25 ranking, stemming, phrase queries — all free |
| `modernc.org/sqlite` over `mattn/go-sqlite3` | No CGO = cross-compilation works out of the box |
| OpenAI-compatible API only | Covers Ollama, LM Studio, OpenAI, Groq, etc. without vendor lock-in |
| Streaming via channel bridge | Keeps engine pure (no bubbletea types); TUI side converts tokens to `tea.Msg` |
| XDG directories | Standard on Linux, reasonable fallbacks on macOS/Windows |
