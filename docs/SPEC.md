# iRecall Technical Specification

## Overview

iRecall is a Bubble Tea TUI backed by a small Go core. The core owns persistence, prompt orchestration, provider communication, and settings. The TUI owns navigation, form state, viewport rendering, and asynchronous command wiring.

The current application flow is:

1. Store a note in SQLite.
2. Ask the configured LLM for 6 to 12 tags.
3. Save tags and update the FTS index.
4. For recall, ask the LLM for 3 to 6 search keywords.
5. Search the note corpus with SQLite FTS5.
6. Stream a grounded answer using only the retrieved notes.

## Repository Structure

```text
iRecall/
├── cmd/
│   └── irecall/
│       └── main.go
├── config/
│   └── config.go
├── core/
│   ├── engine.go
│   ├── models.go
│   ├── db/
│   │   ├── migrations.go
│   │   └── store.go
│   └── llm/
│       ├── client.go
│       └── provider.go
├── tui/
│   ├── app.go
│   ├── pages/
│   │   ├── addquote.go
│   │   ├── quotes.go
│   │   ├── recall.go
│   │   └── settings.go
│   └── styles/
│       └── theme.go
├── Makefile
├── PLAN.md
├── README.md
└── SPEC.md
```

## Runtime and Entry Point

`cmd/irecall/main.go` is responsible for:

- parsing `--debug` and `--version`
- creating XDG-style directories with `config.EnsureDirs()`
- configuring structured JSON logging to `~/.local/state/irecall/irecall.log`
- opening SQLite with migrations
- creating `core.Engine`
- loading saved settings from the database
- starting the Bubble Tea program in the alternate screen

The binary version string is injected with linker flags and defaults to `dev`.

## Data Model

### Quote

```go
type Quote struct {
    ID        int64
    Content   string
    Tags      []string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Settings

```go
type Settings struct {
    Provider ProviderConfig
    Search   SearchConfig
}

type SearchConfig struct {
    MaxResults   int
    MinRelevance float64
}
```

### Provider Configuration

`core.ProviderConfig` is a re-export of `core/llm.ProviderConfig`.

```go
type ProviderConfig struct {
    Host   string
    Port   int
    HTTPS  bool
    APIKey string
    Model  string
}
```

`BaseURL()` renders `<scheme>://<host>:<port>/v1`.

### Default Settings

The default configuration is:

```go
ProviderConfig{
    Host:  "localhost",
    Port:  11434,
    HTTPS: false,
    Model: "",
}

SearchConfig{
    MaxResults:   5,
    MinRelevance: 0.0,
}
```

## Persistence

### File System Layout

`config/config.go` creates and exposes:

- data dir: `~/.local/share/irecall`
- config dir: `~/.config/irecall`
- state dir: `~/.local/state/irecall`

Concrete files currently used:

- database: `~/.local/share/irecall/irecall.db`
- log file: `~/.local/state/irecall/irecall.log`

The config directory exists but does not currently contain a persisted config file.

### SQLite Pragmas

The DB layer enables:

```sql
PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;
PRAGMA busy_timeout = 5000;
```

### Schema

Migration version `1` creates:

```sql
CREATE TABLE IF NOT EXISTS quotes (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    content    TEXT    NOT NULL,
    created_at INTEGER NOT NULL,
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

CREATE VIRTUAL TABLE IF NOT EXISTS quotes_fts USING fts5(
    content,
    tags,
    content='quotes',
    content_rowid='id',
    tokenize='porter unicode61'
);

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

CREATE TABLE IF NOT EXISTS settings (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS schema_version (
    version INTEGER NOT NULL
);
```

### Settings Storage

Settings are persisted as a JSON blob under the `settings` key in the `settings` table.

## Core Engine

`core.Engine` owns:

- the SQLite store
- the active LLM client
- the in-memory settings snapshot

Public behavior currently implemented:

```go
func New(store *db.Store, cfg *Settings) *Engine
func (e *Engine) Close() error

func (e *Engine) UpdateProvider(cfg ProviderConfig)
func (e *Engine) UpdateSettings(s *Settings)

func (e *Engine) AddQuote(ctx context.Context, content string) (*Quote, error)
func (e *Engine) ListQuotes(ctx context.Context) ([]Quote, error)
func (e *Engine) DeleteQuote(ctx context.Context, id int64) error

func (e *Engine) ExtractTags(ctx context.Context, text string) ([]string, error)
func (e *Engine) ExtractKeywords(ctx context.Context, question string) ([]string, error)
func (e *Engine) SearchQuotes(ctx context.Context, keywords []string) ([]Quote, error)
func (e *Engine) GenerateResponse(ctx context.Context, question string, candidates []Quote, tokenCh chan<- string) error

func (e *Engine) FetchModels(ctx context.Context, cfg ProviderConfig) ([]string, error)
func (e *Engine) TestProvider(ctx context.Context, cfg ProviderConfig) error

func (e *Engine) LoadSettings(ctx context.Context) (*Settings, error)
func (e *Engine) SaveSettings(ctx context.Context, s *Settings) error
```

### Add Quote Flow

`AddQuote` performs:

1. trim and validate content
2. insert the quote row
3. call `ExtractTags`
4. on tag success:
   - upsert tags
   - create `quote_tags` associations
   - rewrite the FTS row with tag text included
5. return a `Quote`

If tag extraction fails, the note is still saved without tags.

### Recall Flow

`RecallPage` drives the following engine sequence:

1. `ExtractKeywords(question)`
2. `SearchQuotes(keywords)`
3. render reference notes immediately
4. `GenerateResponse(question, quotes, tokenCh)`
5. drain the token channel recursively into Bubble Tea messages

The response prompt explicitly instructs the model to:

- use only the reference notes
- cite note numbers like `[1]`
- return a fixed insufficiency sentence if the notes are not enough
- stay brief

### Keyword and Tag Parsing

Both extraction methods expect a JSON string array from the model.

`parseJSONStringArray`:

- trims the response
- isolates the first bracketed array if extra text is present
- tries `json.Unmarshal`
- falls back to comma splitting if JSON parsing fails
- lowercases fallback values

## DB Store Behavior

`core/db/store.go` currently implements:

- `Open(path)`
- `Close()`
- `InsertQuote(content)`
- `UpdateQuoteFTS(id, tags)`
- `DeleteQuote(id)`
- `ListQuotes()`
- `SearchQuotes(keywords, limit)`
- `UpsertTags(names)`
- `InsertQuoteTags(quoteID, tagIDs)`
- `GetSetting(key)`
- `SetSetting(key, value)`

### Search Query Behavior

Search currently:

- trims each keyword
- escapes embedded double quotes
- wraps each keyword in FTS phrase quotes
- joins them with `OR`
- orders by `fts.rank`
- limits results by `SearchConfig.MaxResults`

Example generated match expression:

```text
"flash memory" OR "partition" OR "offset"
```

Important current limitation:

- `SearchConfig.MinRelevance` is collected and persisted but is not yet used in `SearchQuotes`.

## LLM Client

`core/llm/client.go` is a minimal OpenAI-compatible HTTP client.

### Supported Calls

- `POST /v1/chat/completions`
- `GET /v1/models`

### Chat Behavior

`Chat(...)`:

- sends `model`, `messages`, and `stream`
- optionally includes `temperature` and `max_tokens`
- adds a bearer token only when `APIKey` is non-empty
- uses a 30-second timeout for non-streaming requests
- disables client timeout for streaming requests

### Streaming

When `tokenCh` is non-nil:

- the request is sent with `stream: true`
- a goroutine parses SSE lines from the response body
- `data: [DONE]` ends the stream
- each `choices[0].delta.content` fragment is emitted as one token
- the token channel is closed when parsing ends

### Model Fetching

`FetchModels`:

- wraps the request in a 10-second context timeout
- decodes `data[].id`
- sorts the model IDs lexicographically before returning them

## TUI Structure

### Root App

`tui/app.go` owns:

- active page state
- overlay state
- the `Recall`, `Quotes`, and `Settings` pages
- the `Add Quote` modal
- global window sizing and page routing

Page cycle:

```text
Recall -> Quotes -> Settings -> Recall
```

Overlay behavior:

- the add-quote modal blocks page routing while open
- `Ctrl+C` still exits globally

### Recall Page

`tui/pages/recall.go` contains:

- single-line `textinput` for the question
- response `viewport`
- reference quotes `viewport`
- spinner-driven busy state
- keyword display
- recursive message pattern for streaming tokens

Recall-specific messages:

- `TokenMsg`
- `RecallDoneMsg`
- `QuotesReadyMsg`
- `KeywordsReadyMsg`
- `OpenAddQuoteMsg`
- internal `quotesAndStreamMsg`
- internal `tokenWithChannel`

### Quotes Page

`tui/pages/quotes.go` provides:

- a scrollable viewport of all saved quotes
- rendered tags per quote
- reload support with `r`
- helpful empty-state messaging

### Settings Page

`tui/pages/settings.go` manages:

- host
- port
- HTTPS toggle
- API key
- fetch-models action
- selected model
- max results
- min relevance

The page validates:

- port: `1..65535`
- max results: `1..20`
- min relevance: decimal number between `0.0` and `1.0`

Model selection behavior:

- before fetch, the page shows the initial saved model string
- after fetch, models are cycled with left and right arrows

### Add Quote Modal

`tui/pages/addquote.go` provides:

- multi-line `textarea`
- spinner while saving
- transient status messaging
- auto-close after a successful save

## Styling

`tui/styles/theme.go` defines the current visual system:

- violet primary and accent colors
- rounded panel borders
- double-border modal
- shared help, status, button, label, and quote styles

## Build and Tooling

The repository currently includes these `Makefile` targets:

- `build`
- `run`
- `test`
- `lint`
- `tidy`
- `install`
- `clean`
- `build-linux-amd64`
- `build-linux-arm64`
- `build-darwin-amd64`
- `build-darwin-arm64`
- `build-windows-amd64`
- `build-all`

The build output defaults to `bin/irecall`.

## Current Gaps

The following are visible in the codebase today:

- no automated tests are checked in
- no web or native UI exists yet
- `MinRelevance` is not enforced in search
- quotes can be deleted through the engine but there is no TUI delete flow
- the config directory is reserved but unused
