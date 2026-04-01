# iRecall — Technical Specification

## 1. Project Structure

```
iRecall/
├── cmd/
│   └── irecall/
│       └── main.go              # Entry point; wires config, engine, TUI
├── config/
│   └── config.go                # XDG path resolution, config file I/O
├── core/
│   ├── models.go                # Shared data structs (Quote, Tag, Settings, ProviderConfig)
│   ├── engine.go                # Engine struct; public API surface
│   ├── db/
│   │   ├── store.go             # SQLite connection, query methods
│   │   └── migrations.go        # Schema version tracking and migration runner
│   └── llm/
│       ├── client.go            # HTTP client for OpenAI-compatible chat/completions
│       └── provider.go          # ProviderConfig, BaseURL, model fetch
├── tui/
│   ├── app.go                   # Root Bubbletea model; page routing; global keys
│   ├── styles/
│   │   └── theme.go             # Lipgloss color palette, borders, layout constants
│   └── pages/
│       ├── recall.go            # Main Q&A page
│       ├── addquote.go          # Add Quote modal overlay
│       └── settings.go          # LLM provider + search settings
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── PLAN.md
└── SPEC.md
```

---

## 2. Data Models

```go
// core/models.go

type Quote struct {
    ID        int64
    Content   string
    Tags      []string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Tag struct {
    ID   int64
    Name string
}

type Settings struct {
    Provider ProviderConfig
    Search   SearchConfig
}

type SearchConfig struct {
    MaxResults       int     // default: 5
    MinRelevance     float64 // default: 0.0 (FTS5 rank threshold; 0 = no filter)
}
```

```go
// core/llm/provider.go

type ProviderConfig struct {
    Host   string // hostname or IP, no scheme
    Port   int
    HTTPS  bool
    APIKey string // empty = no Authorization header sent
    Model  string
}

func (p ProviderConfig) BaseURL() string {
    scheme := "http"
    if p.HTTPS {
        scheme = "https"
    }
    return fmt.Sprintf("%s://%s:%d/v1", scheme, p.Host, p.Port)
}
```

---

## 3. Database Schema

### 3.1 Tables

```sql
CREATE TABLE IF NOT EXISTS quotes (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    content    TEXT    NOT NULL,
    created_at INTEGER NOT NULL,   -- Unix seconds
    updated_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS tags (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT    NOT NULL UNIQUE COLLATE NOCASE
);

CREATE TABLE IF NOT EXISTS quote_tags (
    quote_id INTEGER NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    tag_id   INTEGER NOT NULL REFERENCES tags(id)   ON DELETE CASCADE,
    PRIMARY KEY (quote_id, tag_id)
);

-- FTS5 virtual table
-- 'tags' column stores space-separated tag names for combined search
CREATE VIRTUAL TABLE IF NOT EXISTS quotes_fts USING fts5(
    content,
    tags,
    content='quotes',
    content_rowid='id',
    tokenize='porter unicode61'
);

-- Settings as a key-value store
CREATE TABLE IF NOT EXISTS settings (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

-- Schema versioning
CREATE TABLE IF NOT EXISTS schema_version (
    version INTEGER NOT NULL
);
```

### 3.2 FTS Sync Triggers

```sql
-- Keep quotes_fts in sync with the quotes table
CREATE TRIGGER IF NOT EXISTS quotes_ai AFTER INSERT ON quotes BEGIN
    INSERT INTO quotes_fts(rowid, content, tags)
    VALUES (new.id, new.content, '');
END;

CREATE TRIGGER IF NOT EXISTS quotes_ad AFTER DELETE ON quotes BEGIN
    INSERT INTO quotes_fts(quotes_fts, rowid, content, tags)
    VALUES ('delete', old.id, old.content, '');
END;

CREATE TRIGGER IF NOT EXISTS quotes_au AFTER UPDATE ON quotes BEGIN
    INSERT INTO quotes_fts(quotes_fts, rowid, content, tags)
    VALUES ('delete', old.id, old.content, '');
    INSERT INTO quotes_fts(rowid, content, tags)
    VALUES (new.id, new.content, '');
END;
```

Note: tag updates call `db.UpdateQuoteFTS(id, tags)` which issues the `'delete'`
then re-insert pair directly (not through the trigger) to include updated tag text.

### 3.3 Search Query

```sql
-- keywords joined as: "golang OR channels OR concurrency"
SELECT q.id, q.content, q.created_at, q.updated_at,
       fts.rank
FROM quotes_fts AS fts
JOIN quotes AS q ON q.id = fts.rowid
WHERE quotes_fts MATCH ?
ORDER BY fts.rank          -- BM25: lower value = better match
LIMIT ?;
```

### 3.4 Connection Settings

```go
// Applied after opening the DB connection
PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;
PRAGMA busy_timeout = 5000;
```

---

## 4. Core Engine API

```go
// core/engine.go

type Engine struct {
    store  *db.Store
    llm    *llm.Client
    cfg    *config.Config
}

func New(cfg *config.Config) (*Engine, error)
func (e *Engine) Close() error

// Quote management
func (e *Engine) AddQuote(ctx context.Context, content string) (*Quote, error)
func (e *Engine) ListQuotes(ctx context.Context) ([]Quote, error)
func (e *Engine) DeleteQuote(ctx context.Context, id int64) error

// Recall workflow
func (e *Engine) ExtractKeywords(ctx context.Context, text string) ([]string, error)
func (e *Engine) SearchQuotes(ctx context.Context, keywords []string) ([]Quote, error)
func (e *Engine) GenerateResponse(
    ctx    context.Context,
    question  string,
    candidates []Quote,
    tokenCh   chan<- string,   // receives streamed tokens; closed when done
) error

// Provider management
func (e *Engine) FetchModels(ctx context.Context, p llm.ProviderConfig) ([]string, error)
func (e *Engine) TestProvider(ctx context.Context, p llm.ProviderConfig) error

// Settings
func (e *Engine) LoadSettings(ctx context.Context) (*Settings, error)
func (e *Engine) SaveSettings(ctx context.Context, s *Settings) error
```

### 4.1 AddQuote Flow

```
AddQuote(content)
    │
    ├─ db.InsertQuote(content)           → quote.ID
    │
    ├─ llm.ExtractTags(content)          → []string tags
    │   └─ prompt: tag-extraction (§6.1)
    │
    ├─ db.UpsertTags(tags)               → []int64 tagIDs
    ├─ db.InsertQuoteTags(quoteID, tagIDs)
    ├─ db.UpdateQuoteFTS(quoteID, tags)  → FTS delete+reinsert with tag text
    │
    └─ return &Quote{...}
```

### 4.2 Recall Flow

```
ExtractKeywords(question)              → []string keywords
    └─ prompt: keyword-extraction (§6.2)

SearchQuotes(keywords)                 → []Quote (ranked, limit MaxResults)
    └─ FTS query: "kw1 OR kw2 OR ..."

GenerateResponse(question, quotes, ch)
    ├─ build prompt with numbered quote list (§6.3)
    ├─ stream chat completion
    └─ send each token → ch; close ch when done
```

---

## 5. LLM Client

```go
// core/llm/client.go

type Client struct {
    cfg        ProviderConfig
    httpClient *http.Client
}

func NewClient(cfg ProviderConfig) *Client

// Chat with optional streaming
func (c *Client) Chat(
    ctx     context.Context,
    msgs    []Message,
    stream  bool,
    tokenCh chan<- string,   // nil when stream=false
) (string, error)

// List available models via GET /v1/models
func (c *Client) FetchModels(ctx context.Context) ([]string, error)

type Message struct {
    Role    string // "system" | "user" | "assistant"
    Content string
}
```

**Streaming:** SSE parsing of `data: {...}` lines; extract `choices[0].delta.content`;
send each non-empty token to `tokenCh`; stop on `data: [DONE]`.

**Authorization:** If `APIKey != ""`, set `Authorization: Bearer <key>` header.

**Timeouts:**
- Non-streaming: 30s total
- Streaming: 5s connect, no read deadline (stream may be long)
- `FetchModels`: 10s total

---

## 6. LLM Prompts

### 6.1 Tag Extraction (AddQuote)

```
System:
You are a keyword extractor. Given a piece of text, return a JSON array of
3 to 8 short, lowercase keyword tags that best represent the core concepts.
Return ONLY the JSON array with no explanation.
Example: ["machine learning", "neural networks", "backpropagation"]

User:
<quote content>
```

Response parsing: unmarshal JSON array of strings; if parsing fails, fall back
to splitting on commas as a best-effort recovery.

### 6.2 Keyword Extraction (Recall)

```
System:
You are a search keyword extractor. Given a question, return a JSON array of
3 to 6 lowercase keywords or short phrases most useful for searching a
personal knowledge base. Return ONLY the JSON array.

User:
<question>
```

### 6.3 Response Synthesis (Recall)

```
System:
You are a personal knowledge assistant. Answer the user's question using ONLY
the reference notes provided below. Cite the notes by their number, e.g. [1].
If the notes do not contain enough information, say so clearly.
Be concise and direct.

Reference notes:
[1] <quote 1 content>
[2] <quote 2 content>
[3] <quote 3 content>
...

User:
<question>
```

---

## 7. Configuration

### 7.1 File Format (`config.json`)

```json
{
  "provider": {
    "host": "localhost",
    "port": 11434,
    "https": false,
    "api_key": "",
    "model": "llama3.2:latest"
  },
  "search": {
    "max_results": 5,
    "min_relevance": 0.0
  }
}
```

### 7.2 Paths (XDG)

```go
// config/config.go

func DataDir() string    // $XDG_DATA_HOME/irecall  or ~/.local/share/irecall
func ConfigDir() string  // $XDG_CONFIG_HOME/irecall or ~/.config/irecall
func StateDir() string   // $XDG_STATE_HOME/irecall  or ~/.local/state/irecall

func DBPath() string     // DataDir()/irecall.db
func ConfigPath() string // ConfigDir()/config.json
func LogPath() string    // StateDir()/irecall.log
```

All directories are created with `os.MkdirAll(..., 0700)` on first use.

---

## 8. TUI Specification

Framework: **Bubbletea** (Elm architecture). Each page is a self-contained model.
The root `app` model owns routing and the global `Engine` reference.

### 8.1 Page Routing

```go
type page int
const (
    pageRecall page = iota
    pageSettings
)

type overlayType int
const (
    overlayNone overlayType = iota
    overlayAddQuote
)

type App struct {
    engine   *core.Engine
    page     page
    overlay  overlayType
    recall   RecallPage
    settings SettingsPage
    addQuote AddQuotePage
    width    int
    height   int
}
```

### 8.2 Recall Page Layout

```
┌─ iRecall ─────────────────────────────────────────── [Recall | Settings] ─┐
│                                                                             │
│  Question                                                                   │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ What did I learn about Go channels?                                  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│  Enter: Ask   Ctrl+N: Add Quote   Tab: Settings   Q: Quit                  │
│                                                                             │
│  ─── Response ─────────────────────────────────────────────────────── ▲ ──│
│  Based on your notes, Go channels are...                               │   │
│  Note [1] describes buffered vs unbuffered clearly.                    │   │
│                                                                        ▼   │
│  ─── Reference Quotes ──────────────────────────────────────────────────── │
│  [1] "Unbuffered channels block the sender until a receiver is ready..."   │
│  [2] "Use select{} to multiplex across multiple channels without..."       │
│  [3] "chan<- is send-only; <-chan is receive-only..."                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Components:**
- `textinput.Model` — question input (single line)
- `viewport.Model` — response area (scrollable, streaming)
- `viewport.Model` — reference quotes (scrollable)
- `spinner.Model` — shown during any LLM call
- Status bar line showing key bindings

**Streaming:** `GenerateResponse` runs in a `tea.Cmd`. A goroutine reads from
`tokenCh` and wraps each token in a `TokenMsg`. The Bubbletea `Update` appends
tokens to the response viewport and calls `viewport.SetContent`.

```go
type TokenMsg struct{ Token string }
type RecallDoneMsg struct{ Err error }
type QuotesReadyMsg struct{ Quotes []core.Quote }
```

### 8.3 Add Quote Overlay Layout

```
╔═ Add Quote ═══════════════════════════════════════════════════════════════╗
║                                                                           ║
║  ┌───────────────────────────────────────────────────────────────────┐   ║
║  │ Type or paste your note here.                                     │   ║
║  │ Multi-line input supported.                                       │   ║
║  │                                                                   │   ║
║  │                                                                   │   ║
║  └───────────────────────────────────────────────────────────────────┘   ║
║                                                                           ║
║  Tags will be extracted automatically by the LLM.                        ║
║                                                                           ║
║  [Ctrl+S: Save]  [Esc: Cancel]                                            ║
╚═══════════════════════════════════════════════════════════════════════════╝
```

**Components:**
- `textarea.Model` — multi-line input, 6 visible rows, soft-wrap
- Spinner shown while saving
- Status line: `"Saved." | "Error: <msg>"` shown for 2s then cleared

**Dismiss:** `Esc` discards content and returns to recall page.

### 8.4 Settings Page Layout

```
┌─ iRecall ─────────────────────────────────────────── [Recall | Settings] ─┐
│                                                                             │
│  ─── LLM Provider ──────────────────────────────────────────────────────── │
│                                                                             │
│  Host / IP   [ localhost                              ]                     │
│  Port        [ 11434          ]                                             │
│  HTTPS       [ ] off                                                        │
│  API Key     [ ••••••••••••••••••••••••••••••••••••• ]  (leave blank if    │
│                                                          not required)      │
│                                                                             │
│  [ Fetch Models ]                                                           │
│                                                                             │
│  Model       [ llama3.2:latest                        ▼ ]                   │
│                                                                             │
│  ─── Search ─────────────────────────────────────────────────────────────  │
│                                                                             │
│  Max reference quotes  [ 5   ]  (1–20)                                     │
│  Min relevance score   [ 0.0 ]  (0.0 = no filter)                          │
│                                                                             │
│  Tab: Back   Ctrl+S: Save                                                   │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Components:**
- `textinput.Model` per field (Host, Port, API Key, Max Results, Min Relevance)
- HTTPS: boolean toggle (`Space` to flip)
- Model selector: `list.Model` populated after Fetch Models
- `[Fetch Models]` button: focusable item; `Enter` triggers `Engine.FetchModels`
- Spinner while fetching models
- Status bar: save confirmation or error

**Focus:** `Tab` / `Shift+Tab` cycles through form fields.

### 8.5 Global Key Bindings

| Key | Scope | Action |
|---|---|---|
| `Tab` | Global | Switch between Recall and Settings pages |
| `Q` | Global (no input focused) | Quit |
| `Ctrl+C` | Global | Quit |
| `Ctrl+N` | Recall page | Open Add Quote overlay |
| `Enter` | Question input | Submit recall query |
| `Enter` | Settings — Fetch Models button | Trigger model fetch |
| `Ctrl+S` | Add Quote overlay | Save quote |
| `Ctrl+S` | Settings page | Save settings |
| `Esc` | Add Quote overlay | Dismiss overlay |
| `↑ / ↓` | Viewports | Scroll |
| `PgUp / PgDn` | Viewports | Scroll by page |

---

## 9. Error Handling

- All engine errors surface to the TUI as a status bar message (never a crash).
- LLM errors during recall: show error in response panel; still show any
  already-retrieved reference quotes.
- LLM errors during tag extraction (AddQuote): store the quote with empty tags
  and surface a warning. The quote is still searchable by content via FTS.
- Settings save errors: shown inline; not applied to running engine.
- DB errors: logged to file and surfaced as status messages.
- All long-running operations respect `context.Context` cancellation so that
  quitting the TUI cleanly cancels in-flight requests.

---

## 10. Logging

- Log destination: file only (`StateDir()/irecall.log`) — stdout/stderr would
  corrupt the TUI.
- Format: structured JSON lines via `log/slog`.
- Level: `INFO` in production, `DEBUG` via `--debug` flag.
- Logged events: DB queries (debug), LLM requests/responses (info), errors (error).
- No sensitive data (API keys, quote content) logged at INFO or above.

---

## 11. Build & Distribution

```makefile
# Makefile (excerpt)

VERSION  := $(shell git describe --tags --always --dirty)
LDFLAGS  := -ldflags "-X main.version=$(VERSION) -s -w"

build:
    go build $(LDFLAGS) -o bin/irecall ./cmd/irecall

test:
    go test ./...

lint:
    go vet ./...
    staticcheck ./...

install:
    go install $(LDFLAGS) ./cmd/irecall
```

Cross-compilation targets (no CGO required):
- `linux/amd64`, `linux/arm64`
- `darwin/amd64`, `darwin/arm64`
- `windows/amd64`

---

## 12. Future Considerations (Out of Scope for v1)

- **Web UI:** HTTP server in a `web/` package; same `core.Engine` API.
- **Native UI:** macOS/Windows app via a `native/` package (Wails or similar).
- **Quote editing:** Update quote content and re-extract tags.
- **Tag management UI:** View, rename, merge tags.
- **Export/Import:** JSON or Markdown dump of all quotes and tags.
- **Multiple knowledge bases:** Separate DB files per "vault".
- **Embedding-based search:** Hybrid FTS + vector similarity for semantic recall.
