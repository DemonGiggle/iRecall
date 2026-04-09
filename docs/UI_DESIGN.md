# iRecall UI Design

## Purpose

This document defines the current iRecall UI contract so future clients can implement the same product design consistently.

It is not a Bubble Tea implementation guide. It describes:

1. information architecture
2. visual language
3. screen and modal behavior
4. interaction contracts
5. state and feedback rules

The goal is that a Windows desktop app, another terminal UI, or a web client can preserve the same user experience even if the widgets and platform conventions differ.

## Product Model

iRecall is a personal recall and quote management app with three core surfaces:

1. `Recall`
2. `Quotes`
3. `Settings`

It also uses task-specific modal overlays for blocking work:

1. set user name
2. add/edit quote
3. confirm delete
4. export quotes
5. import quotes

The current product is intentionally local-first:

1. content is stored locally
2. quote sharing is manual import/export
3. settings are local
4. quote ownership and source metadata are visible in browsing contexts

## Design Principles

### 1. Single-purpose surfaces

Each top-level page has one primary job:

1. `Recall` is for asking and grounding answers in quotes
2. `Quotes` is for managing the quote library
3. `Settings` is for provider and search configuration

Do not turn the top-level pages into multi-purpose dashboards.

### 2. Keyboard-first interaction

The current TUI is keyboard-driven. Future clients may add mouse support, but keyboard access must remain a first-class path.

Important implication:

1. every interactive area needs a clear focus state
2. actions need discoverable shortcuts or command affordances
3. modal workflows should be completable without pointer input

### 3. Explicit state over hidden automation

The UI should show:

1. where the user is
2. what panel is focused
3. whether the app is busy
4. what will happen on confirm
5. whether an action succeeded or failed

Avoid silent background behavior that changes user data without obvious feedback.

### 4. Dense but calm presentation

The current UI is compact and information-dense, but avoids clutter by:

1. using bordered panels
2. placing key hints near the area they control
3. showing origin metadata only when useful
4. using modal overlays for branch flows instead of embedding everything inline

Future clients should keep this density and clarity balance.

## Visual Language

The current visual system is defined in [tui/styles/theme.go](/home/gigo/workspace/iRecall/tui/styles/theme.go).

### Palette

Current colors:

1. primary: violet `#7C3AED`
2. accent: light violet `#A78BFA`
3. muted text: gray `#6B7280`
4. success: green `#10B981`
5. error: red `#EF4444`
6. warning: amber `#F59E0B`
7. foreground: near white `#F9FAFB`
8. border: dark gray `#374151`
9. selected/active background: dark slate `#1F2937`

Cross-platform rule:

1. preserve the semantic roles above even if exact colors change
2. active state must be clearly distinct from inactive state
3. success and error colors must remain reserved for outcome feedback

### Typography and emphasis

The current design relies on text styling rather than typography families:

1. bold for titles, active tabs, and important labels
2. accent color for section headers and highlighted metadata
3. muted color for help and secondary explanation

For non-terminal clients:

1. use one strong title style
2. use one section-header style
3. use one muted helper style
4. avoid decorative typography that changes the product tone

### Containers

There are three main container types:

1. `Panel`
   standard rounded border
2. `PanelActive`
   same structure, primary-colored border for focus/active areas
3. `Modal`
   double border and centered placement for blocking overlays

Cross-platform rule:

1. top-level content should remain visually segmented into bordered or card-like regions
2. focus must be visible on the active pane or field group
3. blocking dialogs must appear centered and visually separated from the base surface

## Information Architecture

### App shell

The shell is defined in [tui/app.go](/home/gigo/workspace/iRecall/tui/app.go).

The shell has:

1. a title bar
2. a user greeting when a display name exists
3. top-level tabs for `Recall`, `Quotes`, and `Settings`
4. one active page at a time
5. at most one blocking overlay at a time

### Header contract

The header must always show:

1. product name: `iRecall`
2. active navigation tabs
3. current user greeting: `Hi! <DisplayName>` when available

The greeting is secondary to navigation and should never overpower the page tabs.

### Page model

Only one page is active at a time:

1. `Recall`
2. `Quotes`
3. `Settings`

Navigation is cyclical:

1. forward moves `Recall -> Quotes -> Settings -> Recall`
2. backward moves `Recall <- Quotes <- Settings <- Recall`

Future clients do not need to preserve `Tab` and `Shift+Tab` specifically, but should preserve this three-page model and clear page switching.

## Screen Specifications

## Recall Page

Reference implementation: [tui/pages/recall.go](/home/gigo/workspace/iRecall/tui/pages/recall.go)

### Purpose

The Recall page is the core recall workflow:

1. ask a question
2. extract keywords
3. retrieve relevant quotes
4. generate a grounded answer from those quotes

### Layout

The page has four vertical regions:

1. question input
2. extracted keywords line
3. response panel
4. reference quotes panel

### Input region

Behavior:

1. single primary question field
2. starts focused by default
3. `enter` runs recall only when the input has content and the page is not busy

Visual treatment:

1. active input uses active panel styling
2. placeholder communicates open-ended questioning

### Keywords line

Purpose:

1. show the extracted search keywords that drive retrieval

Rules:

1. show `Keywords: —` when empty
2. after extraction, show keywords in accent styling
3. keywords are informational, not editable in the current design

### Response panel

Purpose:

1. show the streamed grounded answer

Rules:

1. stream progressively
2. auto-scroll to the bottom during generation
3. preserve the current response until the next question starts
4. show the current question at the top of the response panel so the answer remains anchored to its prompt

### Reference quotes panel

Purpose:

1. show the retrieved evidence used for recall
2. allow quote-specific operations in context

Rules:

1. this panel has its own focus state
2. quote-specific actions only work when this panel is focused
3. the panel must show local help for its own commands

Current quote actions here:

1. move through results
2. select/unselect
3. edit
4. delete
5. share/export

### Focus model

The Recall page has two internal focus zones:

1. question input
2. reference quotes panel

Contract:

1. only one is focused at a time
2. the focused zone must be visually distinguished
3. quote actions are disabled unless the reference panel is focused

This focus split must be preserved in other clients, even if implemented via tabs, split panes, or explicit focus rings.

### Busy state

During recall generation:

1. page-level help is replaced by a spinner plus thinking state
2. the app is actively producing output
3. response streaming continues until complete or error

### Empty and error states

When nothing has been asked yet:

1. response panel is empty
2. keywords show placeholder
3. reference quotes panel is empty

When generation fails:

1. keep the page structure
2. show an explicit error message
3. do not hide the user’s context

## Quotes Page

Reference implementation: [tui/pages/quotes.go](/home/gigo/workspace/iRecall/tui/pages/quotes.go)

### Purpose

The Quotes page is the library management surface.

It owns:

1. browsing stored quotes
2. selection
3. add
4. edit
5. delete
6. export/share
7. import

### Layout

The page is a single large panel with:

1. section header showing quote count
2. quote list viewport
3. local help line

### Quote row design

Reference implementation: [tui/pages/quotefunctions.go](/home/gigo/workspace/iRecall/tui/pages/quotefunctions.go)

Each quote item includes:

1. cursor indicator on the active row
2. selection state `[ ]` or `[x]`
3. list index
4. quote content
5. source line for foreign quotes
6. a compact tag preview line in list contexts

List density rule:

1. the quote list is intentionally concise
2. tag previews should show only the first few tags in list view
3. when more tags exist, the preview should indicate that additional tags are hidden
4. quote content in list view should also be truncated to a compact preview so rows stay visually consistent

Detail rule:

1. pressing `enter` on the current quote opens a detail view inside the Quotes page
2. the detail view shows the full quote metadata, including all tags
3. pressing `enter` or `esc` exits detail view and returns to the compact list

Origin rule:

1. if the quote is owned by the current user, omit source metadata by default
2. if the quote is not owned by the current user, show `From: <SourceName>`

This is important and should be preserved across clients.

### Selection model

Selection is lightweight and explicit:

1. there is always a current row when the list is non-empty
2. users may select multiple rows
3. if no rows are selected, actions fall back to the current row

This behavior currently applies to:

1. delete
2. share/export

Future clients should preserve this because it keeps bulk actions efficient while still making single-item action easy.

### Empty state

When there are no quotes:

1. the page remains accessible
2. a direct action hint is shown to add a quote

The empty state should feel actionable, not dead.

## Settings Page

Reference implementation: [tui/pages/settings.go](/home/gigo/workspace/iRecall/tui/pages/settings.go)

### Purpose

Settings manages provider and search configuration.

### Structure

The page is split into two sections:

1. `LLM Provider`
2. `Search`
3. `Local Storage`

### Provider fields

Current provider controls:

1. host / IP
2. port
3. HTTPS toggle
4. API key
5. fetch models action
6. selected model

### Search fields

Current search controls:

1. max reference quotes
2. minimum relevance

### Local storage information

The Settings page also shows read-only local filesystem paths for:

1. data directory
2. config directory
3. state directory

Rules:

1. these paths are informational, not editable in the current design
2. they should reflect the active runtime root, including `-data-path` overrides

### Interaction model

This screen uses a form-navigation model:

1. move focus field-by-field
2. edit text inputs in place
3. toggle HTTPS explicitly
4. fetch models on demand
5. cycle model choices after fetch
6. save settings explicitly

### Feedback

The page must surface:

1. save success
2. save failure
3. model fetch progress
4. model fetch failure
5. empty model list result

## Modal Overlays

All overlays are blocking and centered. The base page remains visually present behind the overlay, but interaction is captured by the overlay.

### Set Your Name

Reference implementation: [tui/pages/userprofile.go](/home/gigo/workspace/iRecall/tui/pages/userprofile.go)

Purpose:

1. collect the required user display name

Rules:

1. shown on launch when display name is missing
2. blocks access to all other screens
3. cannot be dismissed with escape
4. requires non-empty trimmed input

Meaning of the field must be explained:

1. the name is used as quote-source identity in sharing flows

### Add / Edit Quote

Reference implementation: [tui/pages/addquote.go](/home/gigo/workspace/iRecall/tui/pages/addquote.go)

Purpose:

1. create a new quote
2. edit an existing quote

Rules:

1. same modal handles both add and edit
2. save is explicit
3. empty content cannot be saved
4. tags are regenerated automatically after save

### Refine draft preview

This is a sub-state of the quote editor, not a separate top-level modal.

Purpose:

1. let the user ask the LLM to improve wording before saving

Rules:

1. the user’s draft is never silently replaced
2. after refinement, show a side-by-side comparison:
   current draft
   refined draft
3. user explicitly accepts or rejects the refined draft
4. rejecting returns the user to editing

This explicit preview-and-confirm behavior should be preserved in all clients.

### Delete confirmation

Reference implementation: [tui/pages/confirmdelete.go](/home/gigo/workspace/iRecall/tui/pages/confirmdelete.go)

Purpose:

1. prevent accidental destructive delete

Rules:

1. destructive action must require confirmation
2. copy should communicate that deletion is irreversible
3. for multi-select delete, summarize the first few affected quotes

### Share / Export Quotes

Reference implementation: [tui/pages/sharequote.go](/home/gigo/workspace/iRecall/tui/pages/sharequote.go)

Purpose:

1. export selected quotes as a share payload
2. save that payload to a JSON file

Rules:

1. the modal previews the actual payload
2. the modal shows a summary of the quotes being exported
3. save path is user-controlled
4. sharing is currently manual file transfer, not direct transport

### Import Quotes

Reference implementation: [tui/pages/importquote.go](/home/gigo/workspace/iRecall/tui/pages/importquote.go)

Purpose:

1. import a previously exported quote share file

Rules:

1. user provides a file path
2. import result must show:
   inserted
   updated
   duplicates
   stale
3. quote list reloads after successful import flow closes

## Interaction Contracts

## Global app interactions

Current TUI shortcuts:

1. `tab` and `shift+tab` switch top-level pages
2. `ctrl+c` quits

Cross-platform guidance:

1. exact shortcuts may vary
2. top-level page switching must remain simple and always available
3. destructive or blocking actions should never be hidden behind deep navigation

## Local action discoverability

Each page or panel should advertise the actions relevant to itself.

Current pattern:

1. page-level help line for page actions
2. panel-level help line for panel-specific actions
3. modal-level help line for confirm/cancel/submit actions

This should be preserved across clients.

## Feedback model

There are four feedback types:

1. passive helper text
2. progress indicators
3. success status
4. error status

Rules:

1. helper text uses muted styling
2. progress replaces or overrides the normal help line during active work
3. success is visible but brief
4. errors are explicit and localized to the task that failed

## Content and copy rules

### Titles

Use short action-oriented or noun-oriented titles:

1. `Recall`
2. `Stored Quotes`
3. `Settings`
4. `Add Quote`
5. `Import Quotes`

### Helper copy

Helper copy should be:

1. short
2. operational
3. directly tied to the current task

Avoid marketing-style or decorative language.

### Metadata display

Do not surface every backend field all the time.

Current UI shows:

1. content
2. tags in list management contexts
3. source metadata for foreign quotes
4. keyword extraction results in the recall flow

Keep this selective approach in future clients.

## State and synchronization rules

These rules matter for keeping multiple UIs consistent.

### Quote identity

Every client should preserve:

1. quote global identity
2. author identity
3. source identity
4. quote version

The UI does not need to show every field everywhere, but must not lose them in editing or import/export flows.

### Ownership

The UI should distinguish:

1. quotes owned by me
2. quotes received from others

This affects:

1. source display
2. future editing/fork semantics
3. sharing expectations

### Import/export semantics

All clients must preserve the same import result meanings:

1. inserted
2. updated
3. duplicates
4. stale

Do not redefine these terms per platform.

## Cross-platform translation guidance

## What must stay the same

The following are product-level contracts and should remain stable across clients:

1. the three top-level pages
2. blocking startup name capture
3. explicit quote add/edit/delete/share/import flows
4. recall page split between answer generation and evidence review
5. quote source attribution for foreign quotes
6. explicit refine-preview-accept/reject flow
7. manual import/export share model

## What may change per platform

These may adapt to native platform conventions:

1. exact keyboard shortcuts
2. fonts
3. spacing scale
4. border rendering style
5. iconography
6. exact placement of helper hints

But the underlying workflow and UI state machine should remain equivalent.

## Windows desktop adaptation guidance

For a Windows desktop client, the recommended mapping is:

1. top-level tabs remain top-level tabs
2. Recall remains a vertically stacked workspace with an evidence pane
3. Quotes remains a list/detail-lite management surface
4. modal overlays may become native dialogs or centered in-app dialogs
5. keyboard hints may become button labels, tooltips, or command-bar actions

Recommended additions for desktop:

1. file picker for import/export instead of raw path entry
2. richer focus rings
3. optional toolbar buttons mirroring keyboard actions
4. resizable split behavior for recall answer vs evidence

These are platform improvements, not product changes.

## Future maintenance

When the UI changes, update this document if any of these change:

1. page architecture
2. modal inventory
3. ownership/source display rules
4. import/export semantics
5. focus model
6. status and feedback behavior
7. visual design tokens or semantic color roles

If a future UI diverges intentionally from this design, document that divergence explicitly instead of letting platform implementations drift independently.
