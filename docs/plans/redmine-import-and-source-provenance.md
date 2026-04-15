# Redmine Import and Source Provenance

## Goal

Build a script in `tools/` that connects to a Redmine PostgreSQL database and exports Redmine ticket data into iRecall's import quote format.

The script should:

- accept PostgreSQL connection inputs such as host/IP, port, username, password, and database name
- read Redmine data directly from the source database
- transform Redmine records into iRecall's import quote payload
- preserve as much provenance as possible from the original Redmine data
- include the Redmine user information relevant to each exported quote

## Status

This plan is partially implemented.

Implemented:

- generalized source provenance fields in the quote model and share envelope
- migration support and provenance backfill
- Redmine export tool under `tools/redmine_export`

Still open:

- broader validation against real Redmine data
- final UX/documentation decisions for routine import flows
- any future direct integration beyond the current file-based handoff

## Current Findings

### Existing iRecall import format

The current import/export format is `core.SharedQuoteEnvelope` with `core.SharedQuoteEntry`.

The field definitions and provenance semantics are documented in [docs/schema.md](../schema.md).

The short version is:

- existing fields already cover quote identity, authorship, person-level source, content, tags, and import versioning
- they do not fully describe system-level origin such as Redmine issue IDs, journal IDs, LAN node identities, or repository identities

### Historical schema gap

Before the current provenance work landed, the model did not have a complete system-level provenance model for external or network-discovered sources such as:

- source backend: `redmine`, `lan`, `node`, `file_import`
- source scope: a remote system, host, node, or repository
- source object type: `issue`, `journal`, `quote`, `shared_quote`
- source object ID: Redmine issue ID / journal ID / remote quote ID
- source URL or canonical reference

That is why the current provenance fields were added. The person-level source fields alone are still not sufficient for durable external-source identity or future filtering/discovery features.

## Current recommendation

Keep the flexible source metadata model as the foundation for any further import work, and build the remaining Redmine workflow on top of it.

Minimum recommended fields on quotes:

- `source_backend`
- `source_namespace`
- `source_entity_type`
- `source_entity_id`
- `source_label`
- `source_url`

The detailed meaning of these fields is documented in [docs/schema.md](../schema.md).

This source model is meant to support:

- Redmine imports
- LAN quote discovery
- quote repositories or remote nodes
- future filtering by source

## Provenance Design Principles

The source model should be designed as a future filter key, not just descriptive metadata.

Recommended principles:

1. Keep person provenance and system provenance separate.
2. Prefer normalized machine-stable fields over one overloaded text blob.
3. Store enough information to reconstruct a canonical source identity.
4. Allow multiple source families without schema churn.
5. Make exact filtering practical in SQL.

Recommended canonical source identity:

- `(source_backend, source_namespace, source_entity_type, source_entity_id)`

This tuple should uniquely identify an externally sourced quote origin.

Recommended indexing direction for future work:

- index on `source_backend`
- index on `source_namespace`
- composite index on `(source_backend, source_namespace, source_entity_type, source_entity_id)`

This gives us a clean path to future filtering such as:

- all quotes from Redmine
- all quotes from a LAN-discovered node
- all quotes from one specific quote repository
- all quotes sourced from one specific Redmine issue

## Redmine Data Candidates

Based on `tools/ref/redmine-schema/schema.rb`, the likely source tables are:

- `issues`
- `journals`
- `users`
- `projects`
- `trackers`
- `issue_statuses`

Likely exportable quote sources:

1. issue descriptions
2. issue journal notes

Possible mappings:

- quote content:
  - issue description, or
  - journal notes/comments
- author:
  - Redmine `users` row referenced by `issues.author_id` or `journals.user_id`
- timestamps:
  - `issues.created_on` / `issues.updated_on`
  - `journals.created_on`
- tags:
  - project name/identifier
  - tracker
  - status
  - Redmine-specific marker like `redmine`

## Proposed Output Mapping

### For issue descriptions

- `content`: issue subject + description combined into one quote body
- `author_user_id`: `redmine:user:<author_id>`
- `author_name`: resolved Redmine display name
- `source_user_id`: same as author for the initial importer path
- `source_name`: resolved Redmine display name
- `source_backend`: `redmine`
- `source_namespace`: `redmine:<database>`
- `source_entity_type`: `issue_description`
- `source_entity_id`: `<issue_id>`
- `source_label`: `Redmine issue #<issue_id>`
- `created_at_utc`: issue created timestamp
- `updated_at_utc`: issue updated timestamp

### For journal notes

- `content`: journal notes, optionally prefixed with issue context
- `author_user_id`: `redmine:user:<journal_user_id>`
- `author_name`: resolved Redmine display name
- `source_user_id`: same as author for the initial importer path
- `source_name`: resolved Redmine display name
- `source_backend`: `redmine`
- `source_namespace`: `redmine:<database>`
- `source_entity_type`: `issue_journal`
- `source_entity_id`: `<journal_id>`
- `source_label`: `Redmine journal #<journal_id>`
- `created_at_utc`: journal created timestamp
- `updated_at_utc`: journal created timestamp unless a better source exists

## Implementation Phases

### Phase 1: Extend iRecall quote provenance model

- add generalized source metadata fields to `core.Quote`
- add generalized source metadata fields to `core.SharedQuoteEntry`
- update SQLite schema and migrations
- update store read/write paths
- update export/import logic
- add source indexes suitable for future filtering
- update tests for share/import round-trips

Exit criteria:

- iRecall can persist, export, and import external-source metadata

### Phase 2: Define Redmine-to-iRecall mapping contract

- decide exactly which Redmine records become quotes
- decide default content formatting for issues and journals
- decide global ID strategy for imported Redmine items
- decide tag generation rules
- decide whether issue descriptions and journal notes are both enabled by default

Recommended global ID strategy:

- issue description: `redmine:<database>:issue:<id>:description`
- journal note: `redmine:<database>:journal:<id>`

Exit criteria:

- mapping rules are documented and stable

### Phase 3: Build the exporter script in `tools/`

- create a Go script or small Go command under `tools/`
- accept PostgreSQL connection flags
- query Redmine tables
- transform records into `SharedQuoteEnvelope`
- write JSON output file compatible with iRecall import

Recommended CLI shape:

```text
go run ./tools/redmine_export \
  --host 10.0.0.10 \
  --port 5432 \
  --user redmine \
  --password '...' \
  --database redmine_production \
  --output /tmp/redmine-import.json
```

Optional filters for first usable version:

- `--project`
- `--issue-id`
- `--include-journals`
- `--include-issues`

Exit criteria:

- script produces valid iRecall import JSON from a live Redmine PostgreSQL database

### Phase 4: Validate end-to-end import

- import generated JSON through existing iRecall import flow
- verify authors, timestamps, tags, and source metadata
- verify duplicate handling via stable global IDs
- add automated tests around payload generation where feasible

Exit criteria:

- exported Redmine payload imports cleanly into iRecall with correct provenance

## Open Decisions

These need to be confirmed before or during implementation:

1. Should each issue description become a quote, each journal note become a quote, or both?
2. Should blank descriptions or blank journal notes be skipped?
3. Should tags be strictly structured from Redmine metadata, or should we also run the existing tag extractor later inside iRecall?
4. Should `source_url` be populated in the initial Redmine implementation when a base Redmine URL is supplied?
5. Should `source_user_*` continue to mean person-level source while the new source metadata captures system provenance? Current recommendation: yes.
6. Do we want source metadata exposed in the UI immediately, or only persisted and filterable for now?

## Recommended Next Task

Start with Phase 2 and Phase 3.

Phase 1 is already implemented in the codebase.

Specifically:

- lock the exact Redmine mapping rules
- finalize the query set against the Redmine PostgreSQL schema
- implement the exporter under `tools/`
