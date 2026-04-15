# Foundation and Current State

## Purpose

This plan captures the work that has already established the current iRecall baseline. It is a historical and maintenance-oriented plan rather than an active feature proposal.

## Implemented Baseline

### Platform and packaging

- Go module initialized as `github.com/gigol/irecall`
- single-binary build flow via `make build`
- cross-compilation targets in `Makefile`
- version injection via linker flags

### Persistence

- SQLite store
- migration runner
- schema version tracking
- `quotes`, `tags`, `quote_tags`, `quotes_fts`, `settings`, `user_profile`, and migration tracking tables
- FTS triggers for quote insert, update, and delete
- explicit FTS refresh with tag text after tag association writes
- quote identity backfill and generalized source provenance columns

### LLM integration

- OpenAI-compatible provider configuration
- `/v1/chat/completions` support
- `/v1/models` support
- SSE token streaming
- API key support

### Product surfaces

- Bubble Tea TUI with `Recall`, `Quotes`, and `Settings`
- startup user-profile gating
- add/edit quote modal flow
- quote sharing export/import flow
- Wails desktop client that reuses the same Go core

## Current User Flow

1. Configure an LLM endpoint and model in `Settings`.
2. Add notes from `Recall`.
3. Let the model extract tags for each note.
4. Ask a question from `Recall`.
5. Let the model extract search keywords.
6. Search notes through SQLite FTS5.
7. Stream a grounded answer alongside the matched notes.

## Architecture Notes

- `core/` remains independent from the TUI and contains the engine, DB layer, and LLM client.
- `tui/` is a presentation layer that turns asynchronous engine work into Bubble Tea messages.
- settings are currently persisted in SQLite, not in a standalone config file.

## Maintenance Follow-ups

These are not new feature plans, but they are the main places where the current baseline still needs hardening:

- add broader automated coverage around migrations, persistence, and import/export
- tighten provider validation and failure handling in the UI
- keep the storage and provenance model stable as new import/sync features are added
