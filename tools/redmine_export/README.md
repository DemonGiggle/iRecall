# Redmine Export

This tool exports Redmine issue data from PostgreSQL into iRecall's import quote format.

The command emits a JSON file compatible with the current iRecall quote import flow.

## What It Exports

The exporter currently supports:

- issue descriptions
- issue journal notes

Each exported quote includes:

- quote content derived from Redmine
- Redmine author mapped into `author_*`
- source provenance fields such as `source_backend`, `source_namespace`, `source_entity_type`, and `source_entity_id`
- timestamps from Redmine
- structured tags derived from project, tracker, and status metadata

## Requirements

- PostgreSQL client tools installed, especially `psql`
- network or local access to the Redmine PostgreSQL server
- a PostgreSQL user that can read the Redmine database

This tool currently shells out to `psql` rather than using a Go PostgreSQL driver.

## Usage

From the repository root:

```bash
go run ./tools/redmine_export \
  --host 127.0.0.1 \
  --port 5432 \
  --user postgres \
  --password 'your-password' \
  --database redmine \
  --base-url 'https://redmine.example.com' \
  --output /tmp/redmine-import.json
```

## Required Flags

- `--host`
  optional PostgreSQL host or IP; omit it to use the local Unix socket
- `--user`
  PostgreSQL user
- `--database`
  PostgreSQL database name
- `--output`
  output file path for the iRecall import JSON

## Optional Flags

- `--port`
  PostgreSQL port, default `5432`
- `--password`
  PostgreSQL password
- `--base-url`
  optional Redmine web base URL used to build `source_url`
- `--sslmode`
  PostgreSQL SSL mode, default `disable`
- `--include-issues`
  include issue descriptions, default `true`
- `--include-journals`
  include journal notes, default `true`
- `--include-private-notes`
  include private Redmine journal notes, default `false`
- `--project`
  filter by Redmine project identifier, repeatable
- `--issue-id`
  filter by Redmine issue ID, repeatable

## Examples

Export everything from the main Redmine database:

```bash
go run ./tools/redmine_export \
  --user postgres \
  --database redmine \
  --output /tmp/redmine-import.json
```

Export only one project:

```bash
go run ./tools/redmine_export \
  --host 127.0.0.1 \
  --user postgres \
  --database redmine \
  --project my-project \
  --output /tmp/redmine-import.json
```

Export only specific issues:

```bash
go run ./tools/redmine_export \
  --host 127.0.0.1 \
  --user postgres \
  --database redmine \
  --issue-id 123 \
  --issue-id 456 \
  --output /tmp/redmine-import.json
```

Export only issue descriptions, no journal notes:

```bash
go run ./tools/redmine_export \
  --host 127.0.0.1 \
  --user postgres \
  --database redmine \
  --include-issues=true \
  --include-journals=false \
  --output /tmp/redmine-import.json
```

## Output Shape

The output is a `core.SharedQuoteEnvelope` JSON payload using the current iRecall share schema version.

Each entry has:

- `global_id`
- `author_user_id`
- `author_name`
- `source_user_id`
- `source_name`
- `source_backend`
- `source_namespace`
- `source_entity_type`
- `source_entity_id`
- `source_label`
- `source_url`
- `version`
- `content`
- `tags`
- `created_at_utc`
- `updated_at_utc`

## Current Redmine Mapping

Issue descriptions are exported as:

- `source_backend = "redmine"`
- `source_namespace = "redmine:<database>"`
- `source_entity_type = "issue_description"`
- `source_entity_id = "<issue_id>"`

Journal notes are exported as:

- `source_backend = "redmine"`
- `source_namespace = "redmine:<database>"`
- `source_entity_type = "issue_journal"`
- `source_entity_id = "<journal_id>"`

Author identity is mapped from Redmine users as:

- `author_user_id = "redmine:user:<user_id>"`
- `author_name = "<firstname> <lastname>"`, falling back to login if needed

## Import Into iRecall

After generating the JSON file, import it using the existing iRecall import flow.

If you are using the TUI:

1. open the `Quotes` page
2. press `i`
3. enter the JSON file path
4. confirm import

## Local Restore Note

If you restored the local PostgreSQL backups under `tools/ref/redmine-schema/db/`, a typical local export command is:

```bash
go run ./tools/redmine_export \
  --host 127.0.0.1 \
  --user postgres \
  --database redmine \
  --output /tmp/redmine-import.json
```

You can verify the source database first with:

```bash
psql -d redmine -Atqc "SELECT COUNT(*) FROM issues;"
psql -d redmine -Atqc "SELECT COUNT(*) FROM journals;"
psql -d redmine -Atqc "SELECT COUNT(*) FROM users;"
```

## Limitations

- the tool depends on `psql` being available in `PATH`
- it does not yet support incremental sync
- it does not yet support custom Redmine field export
- it currently derives tags only from selected Redmine metadata, not from iRecall's LLM tag extraction flow
