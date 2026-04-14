# iRecall Schema Guide

## Purpose

This document explains the application-level quote and share schema used by iRecall.

Use it for:

- field definitions
- import/export semantics
- provenance semantics
- guidance for future schema changes

This is not a full database dump. It is the conceptual schema guide for the fields iRecall treats as part of its durable data model.

## Quote and Share Fields

The quote model and share envelope currently revolve around these fields:

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

## Existing Identity and Person Fields

### `global_id`

What it means:

- the stable identity of the quote inside iRecall share/import flows

What it is for:

- deduplication across imports
- version comparison during import
- preserving one logical quote across multiple instances

Important note:

- this identifies the quote record in iRecall, not the upstream external system record

### `author_user_id`

What it means:

- the user ID of the person considered the author of the quote content

What it is for:

- authorship tracking
- ownership checks such as `IsOwnedByMe`
- preserving who wrote the content

Examples:

- local quote authored by local user
- shared quote originally written by another iRecall user
- Redmine import mapped to the Redmine author

### `author_name`

What it means:

- the display name of the author

What it is for:

- UI display
- export/import readability

Important note:

- this is presentation data, not a canonical identifier

### `source_user_id`

What it means today:

- the user-level source from which this quote was obtained in iRecall's person-level sharing/import model

What it is for:

- recording who provided the quote
- preserving the person-level source during quote sharing

Important note:

- this is a person field, not a system field
- for local quotes it is usually the same as `author_user_id`

### `source_name`

What it means today:

- the display name associated with `source_user_id`

What it is for:

- UI display
- human-readable share/import provenance

Important note:

- this is not enough for external-system provenance such as Redmine issue IDs or LAN node identities

## Source Provenance Fields

These fields describe where the quote came from as a system or external object.

### `source_backend`

What it means:

- the top-level integration family or transport by which the quote entered iRecall

What it is for:

- coarse filtering
- routing logic in future sync/discovery code
- analytics and debugging

How stable it should be:

- very stable
- this should only change if the quote is considered to come from a different source family entirely

Examples:

- `local`
- `redmine`
- `shared_import`
- `lan`
- `node`

### `source_namespace`

What it means:

- the source scope inside one backend
- this identifies the system, node, host, repository, or database that owns the source object

What it is for:

- grouping quotes by one external system
- distinguishing two different Redmine databases or two different LAN nodes

How stable it should be:

- stable for as long as the upstream source remains the same

Examples:

- `local:user-123`
- `redmine:production`
- `redmine:staging`
- `lan:192.168.1.20`
- `node:quote-repo-a`

### `source_entity_type`

What it means:

- the logical kind of source record inside the namespace

What it is for:

- distinguishing multiple record shapes from the same backend
- allowing type-specific handling later

How stable it should be:

- stable once chosen for a given imported record kind
- should come from a controlled vocabulary in code, not arbitrary free text

Examples:

- `quote`
- `shared_quote`
- `issue_description`
- `issue_journal`

### `source_entity_id`

What it means:

- the stable identifier of the upstream record inside the namespace

What it is for:

- deduplication
- reconstructing the origin
- linking updates from the same external object back to one iRecall quote

How stable it should be:

- as stable as possible
- this should be the real upstream identifier, not a display name

Examples:

- local quote: `<global_id>`
- shared quote: `<global_id>`
- Redmine issue description: `<issue_id>`
- Redmine journal: `<journal_id>`
- node quote: `<remote_quote_id>`

### `source_label`

What it means:

- a human-readable description of the source object

What it is for:

- UI display
- user-facing filtering
- debugging and operator visibility

How stable it should be:

- allowed to change when display wording improves
- should not be used as the canonical machine identity

Examples:

- `Local quote`
- `Shared quote from Alice`
- `Redmine issue #123`
- `Redmine journal #456`
- `Node quote repo-a/quote-77`

### `source_url`

What it means:

- the canonical URL or deep link to the upstream record, when one exists

What it is for:

- opening the original source
- traceability
- debugging

How stable it should be:

- stable if the external system has stable URLs
- empty if no canonical URL exists

Examples:

- `https://redmine.example.com/issues/123`
- `https://redmine.example.com/journals/456`
- empty for local-only or LAN-only records

## Full Source Identity

The intended machine identity is:

- `(source_backend, source_namespace, source_entity_type, source_entity_id)`

Interpretation:

- `source_backend`: which kind of system this came from
- `source_namespace`: which concrete system/node/repository inside that backend
- `source_entity_type`: what kind of record it is
- `source_entity_id`: which exact record

Examples:

- `("local", "local:user-1", "quote", "7a8b...")`
  means a locally authored quote owned by local user `user-1`
- `("shared_import", "share:user-2", "shared_quote", "7a8b...")`
  means a quote imported from another iRecall user's shared payload
- `("redmine", "redmine:production", "issue_description", "123")`
  means the quote originated from Redmine production issue 123
- `("lan", "node:quote-repo-a", "quote", "remote-77")`
  means the quote came from a LAN-discovered node named `quote-repo-a`

## Relationship Between Person Fields And Source Fields

The source schema does not replace person provenance.

Keep these meanings separate:

- `author_user_id` and `author_name`
  who originally wrote or authored the content
- `source_user_id` and `source_name`
  the person-level source associated with the import/share flow
- `source_backend`, `source_namespace`, `source_entity_type`, `source_entity_id`, `source_label`, `source_url`
  the system-level provenance of the quote itself

Example:

- a Redmine journal note written by Alice would have:
  - `author_*` = Alice
  - `source_user_*` = Alice for the initial importer design
  - `source_backend` = `redmine`
  - `source_namespace` = the Redmine database identity
  - `source_entity_type` = `issue_journal`
  - `source_entity_id` = the Redmine journal ID

## Versioning Fields

### `version`

What it means:

- the logical revision number of the quote record inside iRecall import/export flows

What it is for:

- deciding whether an imported quote is newer, duplicate, or stale

Important note:

- this is an iRecall quote version, not necessarily an upstream Redmine version number

### `created_at_utc`

What it means:

- the quote creation timestamp stored in the share envelope

What it is for:

- preserving the original quote timeline across imports

### `updated_at_utc`

What it means:

- the quote update timestamp stored in the share envelope

What it is for:

- preserving modification time
- helping imports retain correct record history

## Content Fields

### `content`

What it means:

- the actual quote text stored in iRecall

What it is for:

- recall/search
- display
- share/import/export

### `tags`

What it means:

- the iRecall tag list attached to the quote

What it is for:

- FTS enrichment
- browsing and filtering

Important note:

- these are iRecall tags, not a direct mirror of upstream source metadata
