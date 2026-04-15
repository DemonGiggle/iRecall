# iRecall Technical Specification

## Overview

iRecall is a local-first quote recall application with two clients built on one Go core:

- a Bubble Tea TUI launched from `cmd/irecall`
- a Wails desktop client under `desktop/`

The shared core owns:

- quote persistence in SQLite
- migration management
- user profile and quote identity bootstrapping
- quote import/export
- recall history persistence
- OpenAI-compatible provider access
- keyword extraction, retrieval, and grounded response generation

The clients own presentation, input handling, and local workflow orchestration.

## Repository Structure

```text
iRecall/
├── cmd/irecall/              # terminal entrypoint
├── config/                   # XDG-style path helpers
├── core/                     # shared domain logic
│   ├── db/                   # SQLite store + migrations
│   └── llm/                  # OpenAI-compatible client
├── desktop/                  # Wails desktop app
│   ├── backend/
│   └── frontend/
├── docs/                     # roadmap, plans, specs, design references
├── tools/
│   └── redmine_export/       # Redmine -> iRecall share payload exporter
├── tui/                      # Bubble Tea app and pages
│   ├── pages/
│   └── styles/
├── Makefile
└── README.md
```

## Runtime Entry Points

### TUI binary

`cmd/irecall/main.go` is responsible for:

- parsing `--debug`, `--version`, and `--data-path`
- initializing XDG-style directories via `config.EnsureDirs()`
- configuring structured JSON logging
- opening SQLite and running migrations
- loading persisted settings
- loading or bootstrapping the local user profile and quote identity
- starting the Bubble Tea app in the alternate screen

### Desktop app

`desktop/backend/app.go` wraps the same core engine for the Wails frontend.

It provides:

- bootstrap state for the frontend shell
- quote CRUD
- quote import/export helpers
- recall-history CRUD
- user-profile save/load
- settings save/load
- recall execution
- save-recall-as-quote actions

## Data Model

The current application-level schema is defined in `core/models.go`. Field-level semantics live in [schema.md](/home/gigo/workspace/iRecall/docs/schema.md).

### Quote

```go
type Quote struct {
    ID               int64
    GlobalID         string
    AuthorUserID     string
    AuthorName       string
    SourceUserID     string
    SourceName       string
    SourceBackend    string
    SourceNamespace  string
    SourceEntityType string
    SourceEntityID   string
    SourceLabel      string
    SourceURL        string
    Content          string
    Tags             []string
    Version          int64
    IsOwnedByMe      bool
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
```

Key points:

- `ID` is the local SQLite row ID.
- `GlobalID` is the durable iRecall quote identity used by import/export deduplication.
- `Author*` captures authorship.
- `SourceUser*` captures person-level source information.
- `SourceBackend` through `SourceURL` capture system-level provenance.
- `IsOwnedByMe` is derived at load time from the local profile.

### User profile

```go
type UserProfile struct {
    UserID      string
    DisplayName string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

The local profile is required for durable quote identity and sharing metadata. If the display name is missing, the clients block normal usage with a prompt.

### Share envelope

```go
const ShareSchemaVersion = 2

type SharedQuoteEnvelope struct {
    SchemaVersion int
    ExportedAt    time.Time
    Quotes        []SharedQuoteEntry
}

type SharedQuoteEntry struct {
    GlobalID         string
    AuthorUserID     string
    AuthorName       string
    SourceUserID     string
    SourceName       string
    SourceBackend    string
    SourceNamespace  string
    SourceEntityType string
    SourceEntityID   string
    SourceLabel      string
    SourceURL        string
    Version          int64
    Content          string
    Tags             []string
    CreatedAtUTC     time.Time
    UpdatedAtUTC     time.Time
}
```

Current compatibility behavior:

- schema version `2` is the current export format
- schema version `1` is still accepted on import
- older payloads are normalized to `shared_import` provenance on import

### Settings

```go
type Settings struct {
    Provider ProviderConfig
    Search   SearchConfig
    Theme    string
}

type SearchConfig struct {
    MaxResults   int
    MinRelevance float64
}
```

Current defaults:

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

Theme: "violet"
```

`MinRelevance` is a normalized `0.0..1.0` threshold:

- `0.0` disables filtering
- `0.3` to `0.7` is the practical recommended range
- `1.0` keeps only very strong keyword coverage matches

## Persistence

### Local paths

`config/config.go` exposes:

- data dir
- config dir
- state dir

By default those resolve to XDG-style locations under `~/.local/share`, `~/.config`, and `~/.local/state`.

Concrete files currently used:

- SQLite database: `data/irecall.db`
- log file: `state/irecall.log`

When `--data-path` is provided, iRecall uses that directory as the root for `data/`, `config/`, and `state/`.

### SQLite configuration

The DB layer enables:

```sql
PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;
PRAGMA busy_timeout = 5000;
```

### Migrations

The migration runner uses `schema_migrations` as the authoritative table and still writes the legacy `schema_version` table for backward compatibility.

Current migration set:

1. `initial_schema`
   - `quotes`, `tags`, `quote_tags`, `quotes_fts`, `settings`
   - FTS triggers for insert, update, delete
2. `quote_identity_and_user_profile`
   - quote identity columns such as `global_id`, `author_*`, `source_user_*`, `version`
   - `user_profile`
3. `quote_source_provenance`
   - `source_backend`, `source_namespace`, `source_entity_type`, `source_entity_id`, `source_label`, `source_url`
   - provenance backfill for existing rows
   - provenance indexes

The schema guide in [schema.md](/home/gigo/workspace/iRecall/docs/schema.md) is the canonical field reference.

## Core Engine

`core.Engine` owns:

- the SQLite store
- the active provider client
- the in-memory settings snapshot
- the active local user profile

Public behavior currently implemented:

```go
func New(store *db.Store, cfg *Settings) *Engine
func (e *Engine) Close() error

func (e *Engine) UpdateProvider(cfg ProviderConfig)
func (e *Engine) UpdateSettings(s *Settings)
func (e *Engine) UpdateUserProfile(profile *UserProfile)

func (e *Engine) AddQuote(ctx context.Context, content string) (*Quote, error)
func (e *Engine) ListQuotes(ctx context.Context) ([]Quote, error)
func (e *Engine) DeleteQuote(ctx context.Context, id int64) error
func (e *Engine) DeleteQuotes(ctx context.Context, ids []int64) error
func (e *Engine) UpdateQuote(ctx context.Context, id int64, content string) (*Quote, error)
func (e *Engine) RefineQuoteDraft(ctx context.Context, content string) (string, error)

func (e *Engine) ExtractTags(ctx context.Context, text string) ([]string, error)
func (e *Engine) ExtractKeywords(ctx context.Context, question string) ([]string, error)
func (e *Engine) SearchQuotes(ctx context.Context, keywords []string) ([]Quote, error)
func (e *Engine) GenerateResponse(ctx context.Context, question string, candidates []Quote, tokenCh chan<- string) error

func (e *Engine) ExportQuotes(ctx context.Context, ids []int64) ([]byte, error)
func (e *Engine) ImportSharedQuotes(ctx context.Context, payload []byte) (ImportResult, error)

func (e *Engine) FetchModels(ctx context.Context, cfg ProviderConfig) ([]string, error)
func (e *Engine) TestProvider(ctx context.Context, cfg ProviderConfig) error

func (e *Engine) LoadSettings(ctx context.Context) (*Settings, error)
func (e *Engine) SaveSettings(ctx context.Context, s *Settings) error

func (e *Engine) LoadUserProfile(ctx context.Context) (*UserProfile, error)
func (e *Engine) SaveUserProfile(ctx context.Context, profile *UserProfile) error
func (e *Engine) BootstrapQuoteIdentity(ctx context.Context) error
```

### Add / update quote flow

`AddQuote` and `UpdateQuote` both:

1. validate content
2. persist the row
3. ask the configured provider for tags
4. upsert tags and associations when tag extraction succeeds
5. rewrite the FTS row so content and tags stay in sync

If tag extraction fails, the quote still persists.

### Recall flow

The current recall pipeline is:

1. trim and validate the question
2. call `ExtractKeywords`
3. run SQLite FTS5 search through `SearchQuotes`
4. apply `MinRelevance` filtering in the engine
5. stream a grounded answer through `GenerateResponse`

Retrieval behavior:

- FTS5 candidate retrieval is still keyword-based
- `MaxResults` controls final returned quote count
- when `MinRelevance > 0`, the engine widens the candidate fetch, filters by normalized keyword coverage, then trims back to `MaxResults`

After a recall completes, the engine can also persist the question/response pair as a normal quote:

1. format the saved content as `Question: ...` plus `Response: ...`
2. reuse the recall keywords when available
3. merge those keywords with freshly extracted tags from the saved content
4. persist the saved quote through the normal quote/tag/FTS path

### Recall history flow

The engine persists completed recall sessions as history entries containing:

1. the original question
2. the grounded response
3. the exact reference quotes used for that response
4. the history entry creation timestamp

Clients can:

- list history summaries
- open one history entry in full detail
- delete selected history entries
- save a history entry as a quote, regenerating recall keywords when necessary

### Import / export behavior

Export:

- selected quotes are serialized into `SharedQuoteEnvelope`
- current schema version is written
- source provenance is preserved

Import:

- payload is validated
- rows are matched by `global_id`
- newer versions overwrite older ones
- equal versions count as duplicates
- older incoming versions count as stale
- schema version `1` payloads are normalized to current provenance defaults

## TUI Contract

The TUI shell owns:

- page routing between `Recall`, `History`, `Quotes`, and `Settings`
- blocking overlays:
  - user-profile prompt
  - quote editor
  - delete confirmation
  - share/export
  - import

Important behaviors:

- `Tab` and `Shift+Tab` cycle pages
- the first run is gated by the user-profile prompt until a display name is saved
- quote add/edit uses one shared editor with refine preview support
- share/export and import are file-based
- Recall and History both support saving a question/response pair as a quote

See [UI_DESIGN.md](/home/gigo/workspace/iRecall/docs/UI_DESIGN.md) for the higher-level UI contract.

## Desktop Contract

The Wails desktop client reuses the same engine and data model.

Current desktop responsibilities:

- bootstrap frontend state
- run recall and quote CRUD through backend methods
- support import/export through backend file helpers
- expose settings, model fetch, and user-profile operations

See:

- [WAILS_DESKTOP.md](/home/gigo/workspace/iRecall/docs/WAILS_DESKTOP.md)
- [desktop/README.md](/home/gigo/workspace/iRecall/desktop/README.md)

## External Tools

### Redmine exporter

`tools/redmine_export` exports Redmine issue descriptions and journal notes into the iRecall share-envelope format.

Current characteristics:

- connects through the `psql` CLI
- emits schema version `2` payloads
- maps Redmine authors into iRecall author/source fields
- populates source provenance as `redmine` records

See [tools/redmine_export/README.md](/home/gigo/workspace/iRecall/tools/redmine_export/README.md).
