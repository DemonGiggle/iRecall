# iRecall Implementation Status and Forward Plan

## Current State

The repository is no longer at the scaffold stage. The following are implemented and working at the code level:

- Go module initialized as `github.com/gigol/irecall`
- Bubble Tea TUI with `Recall`, `Quotes`, and `Settings` pages
- Add Quote modal overlay with asynchronous save flow
- SQLite persistence with migrations, FTS5, and tag joins
- OpenAI-compatible HTTP client for chat completions and model listing
- Streaming response handling in the TUI
- XDG-style directory setup for data, config, and state
- File-based structured logging
- `Makefile` build, run, test, lint, install, tidy, clean, and cross-build targets

## What the App Does Today

### End-user flow

1. Configure an LLM endpoint and model on the Settings page.
2. Add notes from the Recall page with `Ctrl+N`.
3. Let the model extract tags for each note.
4. Ask a question on the Recall page.
5. Let the model extract search keywords.
6. Search notes through SQLite FTS5.
7. Stream a grounded answer alongside the matched notes.

### Core architecture

- `core/` remains independent from the TUI and contains the engine, DB layer, and LLM client
- `tui/` is a presentation layer that turns asynchronous engine work into Bubble Tea messages
- settings are persisted in SQLite, not a standalone config file

## Completed Work by Area

### Platform and packaging

- [x] Module initialized
- [x] Single-binary build flow via `make build`
- [x] Cross-compilation targets in `Makefile`
- [x] Version injection via linker flags

### Persistence

- [x] SQLite store
- [x] Migration runner
- [x] Schema version tracking
- [x] `quotes`, `tags`, `quote_tags`, `quotes_fts`, and `settings` tables
- [x] FTS triggers for quote insert, update, and delete
- [x] Explicit FTS refresh with tag text after tag association writes

### LLM integration

- [x] OpenAI-compatible provider configuration
- [x] `/v1/chat/completions` support
- [x] `/v1/models` support
- [x] SSE token streaming
- [x] API key support

### TUI

- [x] Root app with page routing
- [x] Recall page
- [x] Quotes page
- [x] Settings page
- [x] Add Quote modal
- [x] Responsive resizing through `tea.WindowSizeMsg`

## Known Gaps

These are the main areas where the implementation still lags the intended product shape:

- [ ] `SearchConfig.MinRelevance` is collected and saved but not applied in search queries
- [ ] No automated tests are present even though `make test` is wired
- [ ] No UI flow exists for deleting quotes, despite engine support
- [ ] The config directory is created but not used for stored configuration
- [ ] No CI workflow is present in the repository
- [ ] No web or native client exists yet

## Recommended Next Steps

### Short-term

- [ ] Apply `MinRelevance` to the FTS query so saved search settings affect retrieval
- [ ] Add core tests for DB migrations, FTS search, settings persistence, and LLM parsing fallbacks
- [ ] Add provider validation in the UI for empty host and empty model
- [ ] Surface quote timestamps on the Quotes page

### Medium-term

- [ ] Add quote deletion and maybe editing from the Quotes page
- [ ] Add cancellation for in-flight recall requests when the user asks a new question or exits
- [ ] Add integration tests around the full add-then-recall path with a mock provider
- [ ] Add a GitHub Actions workflow for build and test

### Longer-term

- [ ] Expose the same `core/` engine through a web UI
- [ ] Consider import and export flows for notes
- [ ] Add better retrieval controls such as AND/OR search strategies, tag filters, or score thresholds

## Non-goals for the Current Codebase

The repository does not currently include:

- background sync
- multi-user support
- remote storage
- embeddings or vector search
- a REST API
- any GUI outside the terminal
