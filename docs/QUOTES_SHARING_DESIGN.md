# Quote Sharing Design

## Purpose

This document describes the current quote-sharing model in iRecall and the main design constraints around it.

It is no longer a speculative proposal. The base sharing workflow, quote identity model, user profile gating, and source provenance model are implemented.

For field definitions, use [schema.md](/home/gigo/workspace/iRecall/docs/schema.md).

## Current Product Model

iRecall uses a manual, local-first sharing flow:

1. a user creates or edits quotes locally
2. the user exports one or more selected quotes to a JSON file
3. another user imports that JSON file into their own instance
4. the importer deduplicates and updates rows by `global_id` and `version`

There is no automatic network sync, background peer discovery, or authenticated remote transport in the current product.

## Implemented Identity Model

### User identity

The local instance has a persisted user profile:

```go
type UserProfile struct {
    UserID      string
    DisplayName string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

Rules:

- `UserID` is the durable local identity
- `DisplayName` is required for normal app usage
- first-run blocks on a name prompt if the display name is empty

### Quote identity

Each quote has both a local row ID and a durable sharing identity:

```go
type Quote struct {
    ID                int64
    GlobalID          string
    Version           int64
    AuthorUserID      string
    AuthorName        string
    SourceUserID      string
    SourceName        string
    SourceBackend     string
    SourceNamespace   string
    SourceEntityType  string
    SourceEntityID    string
    SourceLabel       string
    SourceURL         string
    ...
}
```

Important distinction:

- `ID` is local SQLite identity only
- `GlobalID` is the durable identity for import/export reconciliation
- `Author*` captures authorship
- `SourceUser*` captures person-level source
- `SourceBackend` through `SourceURL` capture system-level provenance

## Share Envelope

The current transport format is:

```go
const ShareSchemaVersion = 2

type SharedQuoteEnvelope struct {
    SchemaVersion int
    ExportedAt    time.Time
    Quotes        []SharedQuoteEntry
}
```

Current compatibility rules:

- exports use schema version `2`
- imports accept schema version `1` and `2`
- version `1` payloads are normalized into current provenance defaults during import

## Import Semantics

Current importer behavior:

1. decode and validate the envelope
2. look up each quote by `global_id`
3. if not found:
   insert a new row
4. if found and incoming `version` is newer:
   update the existing row
5. if found and incoming `version` is equal:
   count it as a duplicate
6. if found and incoming `version` is older:
   count it as stale

`ImportResult` currently reports:

- inserted
- updated
- duplicates
- stale

## Current Sharing Constraints

The current implementation intentionally keeps sharing simple:

- sharing is file-based, not network-based
- provenance is preserved, but there is no trust or signature model
- versioning is row-level, not field-level
- imports are last-writer-wins based on quote version ordering

## Design Decisions That Landed

These design goals are now implemented:

- durable quote identity via `global_id`
- local user profile for quote attribution
- versioned share envelope
- import deduplication and update semantics
- person-level source metadata
- generalized source provenance metadata

## Gaps Still Open

These areas are still deliberately unresolved or deferred:

- authenticated quote exchange
- collaborative or multi-writer conflict resolution
- fork-on-edit behavior for imported quotes
- background sync and LAN discovery
- remote transport beyond manual files

Those gaps are acceptable for the current product because the app is still explicitly local-first.

## Relationship To Source Provenance

Quote sharing now sits on top of the broader provenance model:

- local quotes use `source_backend=local`
- imported quotes from older share payloads normalize to `source_backend=shared_import`
- Redmine-exported quotes can arrive with `source_backend=redmine`

That separation is intentional:

- quote-sharing metadata explains how the quote moved between iRecall instances
- provenance metadata explains where the quote content originally came from as a system record

## Future Work

Likely next layers, if the product grows beyond manual file exchange:

- peer-to-peer or LAN discovery
- source-based filtering in the UI
- authenticated share payloads
- clearer edit semantics for quotes not authored by the local user

The implementation work for external systems is tracked separately in [docs/plans/redmine-import-and-source-provenance.md](/home/gigo/workspace/iRecall/docs/plans/redmine-import-and-source-provenance.md).
