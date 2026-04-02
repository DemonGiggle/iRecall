# Quotes Sharing Design

## Goal

Add a quote-sharing feature so users can share selected quotes with other users, keep a stable world-wide identity for each shared quote, propagate later updates to prior recipients of the same quote, and show quote origin in the UI when a quote came from someone else.

This document is a product + technical design for the first real version of the feature. It is intentionally biased toward a simple, defensible model that fits the current codebase.

## Product Model

### Core behavior

1. A user creates a quote locally.
2. The user shares one or more selected quotes to another user.
3. Every shared quote has a globally unique quote ID that remains stable across devices and recipients.
4. If the original author later edits that same quote and shares it again to a recipient who already has it, the recipient should receive an update instead of a duplicate.
5. Quotes that are not authored by the current user should show the source user name in the UI.
6. On first launch, if the local user has not set a profile name yet, the app must block normal usage with a modal prompt asking for the name and explaining why it is needed.

### Non-goals for v1

1. Multi-author collaborative editing on the same quote.
2. Merge conflict resolution between two independent editors.
3. Per-field patch sync.
4. Background network sync without an explicit share/import action.
5. Authentication stronger than a user identity/profile name plus transport-level trust.

## Recommended Ownership Model

Use a single-writer model per quote.

1. Every quote has one canonical author.
2. Only the canonical author can publish authoritative updates for that quote ID.
3. Recipients may view, search, and re-share the quote, but if they edit it locally that edit must fork into a new quote with a new global ID unless they are the author.

This avoids ambiguous semantics for “the same quotes being modified could also be shared again”. Without this rule, the system needs operational transforms or conflict-free replicated data types, which is far beyond the current app.

## Identity Model

### User identity

Add a persisted local profile:

```go
type UserProfile struct {
    UserID      string // UUID
    DisplayName string
    CreatedAt   time.Time
}
```

Rules:

1. `UserID` is a stable UUID generated once on first-run.
2. `DisplayName` is required for sharing.
3. `DisplayName` is user-editable later from Settings.
4. The startup modal is shown if `DisplayName` is empty.

Why both fields:

1. `DisplayName` is user-facing and changeable.
2. `UserID` is the durable machine identity used in quote ownership and update routing.

## Quote Identity Model

The current `Quote.ID int64` is only a local database row ID. It cannot serve as a cross-user identity.

Add a global quote identity layer:

```go
type Quote struct {
    ID             int64
    GlobalID       string // UUID
    AuthorUserID   string // UUID
    AuthorName     string
    SourceUserID   string // UUID
    SourceName     string
    Content        string
    Tags           []string
    Version        int64
    IsOwnedByMe    bool
    CreatedAt      time.Time
    UpdatedAt      time.Time
    SharedAt       *time.Time
}
```

Field semantics:

1. `ID`: local SQLite row ID only.
2. `GlobalID`: stable quote UUID shared across all copies of the same quote.
3. `AuthorUserID`: canonical owner of the quote and must be a UUID.
4. `AuthorName`: latest known display name for the author, denormalized for convenience.
5. `SourceUserID`: the user who directly sent this copy to me and must be a UUID.
6. `SourceName`: user-facing label shown in the UI.
7. `Version`: monotonically increasing author-controlled version number.
8. `IsOwnedByMe`: derived from `AuthorUserID == local profile UserID`.
9. `SharedAt`: when this local copy was most recently imported/shared.

Important distinction:

1. `Author*` answers “who owns this quote?”
2. `Source*` answers “who did I receive it from?”

For v1, if sharing is always direct author-to-recipient, `Author*` and `Source*` will usually match. Keeping both now prevents repainting the model later.

## Sync / Sharing Envelope

Do not share raw DB rows. Share a transport envelope:

```go
type SharedQuoteEnvelope struct {
    SchemaVersion int                `json:"schema_version"`
    ExportedAt    time.Time          `json:"exported_at"`
    Quotes        []SharedQuoteEntry `json:"quotes"`
}

type SharedQuoteEntry struct {
    GlobalID       string    `json:"global_id"`      // UUID
    AuthorUserID   string    `json:"author_user_id"` // UUID
    AuthorName     string    `json:"author_name"`
    SourceUserID   string    `json:"source_user_id"` // UUID
    SourceName     string    `json:"source_name"`
    Version        int64     `json:"version"`
    Content        string    `json:"content"`
    Tags           []string  `json:"tags"`
    CreatedAtUTC   time.Time `json:"created_at_utc"`
    UpdatedAtUTC   time.Time `json:"updated_at_utc"`
}
```

`SchemaVersion` is required so the share format can evolve independently of the quote `Version` field.

### Schema compatibility rules

1. `quote.version` tracks changes to a quote lineage.
2. `schema_version` tracks changes to the export/import file format.
3. These two version numbers are unrelated and must never be conflated.

Importer policy:

1. If `schema_version` is supported, import normally.
2. If `schema_version` is newer than this app supports, reject the payload with a clear error.
3. If `schema_version` is older, route through a version-specific importer if still supported.
4. Unknown JSON fields should be ignored when the schema version itself is supported.

Recommended implementation pattern:

```go
switch env.SchemaVersion {
case 1:
    return importV1(env)
default:
    return fmt.Errorf("unsupported share schema version: %d", env.SchemaVersion)
}
```

## Update Semantics

### Share/import rule

When importing a shared quote:

1. Look up by `global_id`.
2. If not found:
   create a new local row.
3. If found:
   compare versions.

Apply this policy:

1. Incoming version greater than local version:
   overwrite content, tags, author/source metadata, updated time, and version.
2. Incoming version equal to local version:
   treat as duplicate/no-op, but refresh source metadata if needed.
3. Incoming version less than local version:
   ignore as stale.

### Author edit rule

When the current user edits a quote:

1. If `IsOwnedByMe` is true:
   update the quote in place and increment `Version`.
2. If `IsOwnedByMe` is false:
   do not modify the shared quote in place.
   create a fork:
   new `GlobalID`
   `AuthorUserID = me`
   `AuthorName = my display name`
   `Source* = previous quote author/source or me depending on UX choice`
   `Version = 1`

This is the cleanest way to avoid unauthorized overwrites of somebody else’s quote lineage.

## Persistence Design

### New tables

Recommended migration:

```sql
ALTER TABLE quotes ADD COLUMN global_id TEXT;
ALTER TABLE quotes ADD COLUMN author_user_id TEXT;
ALTER TABLE quotes ADD COLUMN author_name TEXT;
ALTER TABLE quotes ADD COLUMN source_user_id TEXT;
ALTER TABLE quotes ADD COLUMN source_name TEXT;
ALTER TABLE quotes ADD COLUMN version INTEGER NOT NULL DEFAULT 1;
ALTER TABLE quotes ADD COLUMN shared_at INTEGER;

CREATE UNIQUE INDEX IF NOT EXISTS idx_quotes_global_id ON quotes(global_id);

CREATE TABLE IF NOT EXISTS user_profile (
    user_id      TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    created_at   INTEGER NOT NULL,
    updated_at   INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS quote_share_history (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    quote_global_id     TEXT NOT NULL,
    shared_by_user_id   TEXT NOT NULL,
    shared_to_user_id   TEXT NOT NULL,
    shared_to_name      TEXT NOT NULL,
    shared_version      INTEGER NOT NULL,
    created_at          INTEGER NOT NULL
);
```

### Backfill strategy

For all existing quotes during migration:

1. generate `global_id = uuidv7` or UUIDv4
2. assign `author_user_id = local profile user_id`
3. assign `author_name = local profile display_name` if available, otherwise empty
4. assign `source_user_id = author_user_id`
5. assign `source_name = author_name`
6. assign `version = 1`

If the user profile does not exist yet at migration time:

1. create `user_id`
2. leave `display_name` empty
3. block app with first-run name modal until name is provided

## UI Design

### Startup name modal

Trigger:

1. App launch
2. `user_profile.display_name == ""`

Behavior:

1. Modal blocks all other screens.
2. Single input field: user name.
3. Description text:
   “Your name is attached to quotes you share with others and shown as the source when someone receives your quotes.”
4. Save action:
   enabled only when non-empty after trim.
5. Escape should not dismiss the modal.

### Quotes list / Recall references

For any quote not owned by me:

1. show a source line:
   `From: Alice`

For my own quotes:

1. optionally show:
   `From: me`
2. or omit origin text to reduce noise.

Recommended rule:

1. show origin only for non-owned quotes in normal browsing
2. show detailed metadata in a quote details/share dialog

### Share action

Add share entry points where quote selection already exists:

1. `Quotes` page for bulk share from selected quotes
2. `Recall` reference panel for share of selected retrieved quotes

Suggested hotkey:

1. `s` = share selected quotes

Share dialog fields:

1. recipient identifier
2. recipient display name if needed
3. transfer preview with quote count and versions

This design intentionally does not hard-code transport. The dialog should call a share service interface.

## Service Layer Design

Add a dedicated sharing service abstraction instead of putting all logic in `Engine`:

```go
type ShareService interface {
    ExportQuotes(ctx context.Context, ids []int64, source UserProfile, recipient Recipient) ([]SharedQuoteEnvelope, error)
    ImportQuotes(ctx context.Context, envelopes []SharedQuoteEnvelope) (ImportResult, error)
}
```

And extend `Engine` with thin orchestration methods:

```go
func (e *Engine) LoadUserProfile(ctx context.Context) (*UserProfile, error)
func (e *Engine) SaveUserProfile(ctx context.Context, profile *UserProfile) error
func (e *Engine) ShareQuotes(ctx context.Context, ids []int64, recipient Recipient) error
func (e *Engine) ImportSharedQuotes(ctx context.Context, payload []byte) (ImportResult, error)
```

This keeps transport and import rules isolated from recall/search behavior.

## Transport Options

The product statement says “share to other user” but does not define transport. The internal data model should not depend on one transport.

Recommended order:

### v1

Manual import/export.

1. Export selected quotes as JSON payload or file.
2. User sends file/text out-of-band.
3. Recipient imports the file.

Why:

1. no server required
2. validates the identity/update model first
3. easiest to test

### v2

Direct online delivery through a server or relay.

Server responsibilities:

1. user lookup
2. delivery queue
3. authentication
4. recipient inbox

The same `SharedQuoteEnvelope` format should still be used.

## Search / Recall Behavior

Shared quotes should participate in search exactly like local quotes.

No recall-model changes are needed beyond:

1. index shared quotes in FTS as normal
2. optionally include source label in the rendered reference list

Do not inject source metadata into FTS ranking for v1 unless there is a user need to search by source.

## Conflict and Safety Rules

### Duplicate import

Handled by `global_id` uniqueness.

### Stale update

Ignored by version check.

### Same display name for different users

Allowed, because the actual identity key is `user_id`.

### User renames themselves

Future shares carry the new `author_name`.
Existing received quotes may keep old denormalized names until refreshed by a later import or explicit metadata sync.

### Recipient edits shared quote

Must fork.

This rule should be visible in the UI:

“Editing a shared quote creates your own copy.”

## Migration / Rollout Plan

### Phase 1

Data model and local profile

1. add user profile persistence
2. add quote identity columns
3. add startup name modal
4. backfill existing quotes

### Phase 2

Manual sharing

1. share/export selected quotes to JSON
2. import quotes from JSON
3. apply version-based update semantics
4. show origin in Quotes and Recall UIs

### Phase 3

Refinements

1. share history
2. better duplicate/update messaging
3. quote detail metadata
4. optional fork-on-edit prompt for non-owned quotes

### Phase 4

Online delivery

1. recipient addressing
2. inbox/outbox
3. auth
4. background fetch

## Recommended API / DB Decisions

These are the main design decisions I recommend locking in now:

1. Add a stable local `user_id` plus editable `display_name`.
2. Add `global_id` and `version` to quotes.
3. Treat quote sharing as single-writer, author-owned data.
4. Fork when a recipient edits a foreign quote.
5. Store both author and source metadata on every quote.
6. Start with manual import/export before any network sharing.
7. Block app startup with a required name modal if `display_name` is unset.

## Open Questions

These need product decisions before implementation:

1. Is v1 sharing manual import/export, or do you already want a network transport?
2. Should recipients be allowed to re-share somebody else’s quote while preserving original author attribution?
3. When a non-owner edits a received quote, should the app auto-fork or ask first?
4. Should author/source names be editable in-place when imported metadata changes?
5. Do we want per-user share permissions later, or is “once received, can re-share” acceptable?

## Suggested Next Implementation Slice

The smallest safe implementation slice is:

1. add `user_profile` persistence and startup name modal
2. add `global_id`, `author_*`, `source_*`, and `version` to quotes
3. backfill existing quotes
4. render source attribution for non-owned quotes
5. add JSON export/import for selected quotes
6. implement import-by-`global_id` with version overwrite

That slice proves the core model before any network work.
