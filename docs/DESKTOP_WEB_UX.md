# Desktop + Web UX Redesign

## Purpose

This document replaces the old "desktop/web should behave like the TUI" assumption.

It defines the target UX for the Wails desktop client and the HTTP web client.

The shared core product remains the same:

1. ask grounded recall questions
2. manage a local quote library
3. review prior recall sessions
4. configure the local model connection and retrieval rules

The interaction model changes substantially for desktop and web:

1. less modal churn
2. less selection-first behavior
3. more split-view and detail-pane work
4. clearer separation between primary and secondary actions
5. progressive disclosure for advanced configuration

## Current UX Problems

The current frontend is functional but structurally too close to the TUI.

### 1. The app is navigation-heavy before it is task-oriented

`Recall`, `History`, `Quotes`, and `Settings` are presented as equal siblings even though the actual primary job is asking and refining recall answers.

Effect:

1. the user has no obvious "home" workflow
2. `History` and `Quotes` feel like independent apps instead of support surfaces
3. first-time users must infer where to start

### 2. Recall mixes generation and library management

The Recall page lets users ask a question, then immediately exposes quote add, edit, delete, and share in the same top toolbar.

Effect:

1. the primary task loses hierarchy
2. destructive actions compete visually with the ask flow
3. the user is encouraged to "manage quotes" while they are trying to interpret an answer

### 3. History is list-first and context-switches too hard

History behaves like a separate browse page and then flips into a full detail state with a Back button.

Effect:

1. the user loses browsing context when opening one entry
2. comparing multiple past recalls is awkward
3. `History` feels slower than necessary for a desktop product

### 4. Quotes is bulk-action-first instead of reading-first

The Quotes page emphasizes selection, toolbar actions, and checkbox state before it supports search, scan, filtering, or quick editing.

Effect:

1. routine library browsing feels operational rather than thoughtful
2. the user is pushed toward multi-select even when they only want to inspect one quote
3. scale will break this flow as the library grows

### 5. Modals are doing normal-page work

Quote editing, import, export, success notices, and onboarding all use blocking overlays.

Effect:

1. the app repeatedly interrupts itself
2. desktop and web lose the advantage of larger persistent layouts
3. import/export feel like technical operations instead of simple product tasks

### 6. Settings is flat, technical, and over-exposed

Connection settings, retrieval rules, theme, web port, password change, and storage paths all appear with similar weight.

Effect:

1. novice users see too much too early
2. advanced fields look equally mandatory
3. the page reads like a debug form rather than product preferences

## Design Direction

Desktop and web should feel like a local knowledge workspace, not a terminal transported into a browser window.

### Principles

1. `Ask first`
   The main route should revolve around asking, reviewing evidence, and deciding what to keep.
2. `Read before manage`
   Quote and history views should privilege scanning, understanding, and detail inspection before bulk actions.
3. `Inline over modal`
   Use panes, drawers, and expandable sections for routine work. Reserve modals for destructive confirmation and first-run gating.
4. `Progressive disclosure`
   Keep advanced controls available but visually secondary.
5. `One clear next step`
   Each screen should make the next useful action obvious.

## Revised Information Architecture

### Primary navigation

Desktop/web should use four destinations, but not as equal tab peers in the current shape:

1. `Ask`
2. `Library`
3. `Activity`
4. `Settings`

Recommended shell:

1. left sidebar or compact top nav with the four destinations
2. persistent global action for `New Quote`
3. secondary profile/theme utilities in the chrome, not in page content

Renaming guidance:

1. `Recall` becomes `Ask`
2. `Quotes` becomes `Library`
3. `History` becomes `Activity`

This uses language that better matches user intent.

## Key Flows

### First-run flow

Do not drop new users into the full app with a blocking name prompt and a blank shell.

Use a lightweight setup flow:

1. step 1: set display name
2. step 2: connect model endpoint
3. step 3: test connection and choose a model
4. step 4: optional import or add first quote
5. finish on `Ask`

Rules:

1. keep it skippable only for optional steps
2. do not show advanced retrieval settings during onboarding
3. frame the setup around "make your first answer work"

### Ask flow

`Ask` is the primary workspace.

#### Layout

Use a two-column workspace on desktop and a stacked flow on narrow web widths.

Left column:

1. question composer
2. answer output
3. answer actions

Right column:

1. evidence summary
2. reference quote list
3. quick detail preview for the selected quote

#### Composer behavior

The question composer should be more than a bare input:

1. multiline field with room for natural prompts
2. visible primary CTA: `Ask`
3. optional secondary CTA: `Use previous question`
4. helper text that explains grounding in local quotes

#### During generation

Show a clear run state:

1. `Extracting keywords`
2. `Finding evidence`
3. `Generating grounded answer`

Do not expose quote edit/delete/share in the main Ask toolbar during generation.

#### After generation

Answer actions should be:

1. `Save as Quote`
2. `Open in Activity`
3. `Copy Answer`

Evidence actions should be contextual on each quote card or in the quote detail pane:

1. `Open in Library`
2. `Edit`
3. `Share`

Delete should never be a prominent Ask-surface action.

### Library flow

`Library` is a browse-and-curate experience, not a bulk-operations screen.

#### Layout

Use list-detail.

Left rail:

1. search field
2. filter chips for tags, ownership, source, and recent updates
3. optional sort menu

Center list:

1. quote cards with stronger content preview
2. author/source metadata
3. tags
4. updated date

Right detail pane:

1. full quote
2. provenance
3. tags
4. actions

#### Action hierarchy

Primary actions:

1. `New Quote`
2. `Import`

Secondary actions in detail pane:

1. `Edit`
2. `Share`
3. `Duplicate`

Destructive action:

1. `Delete` only inside detail pane or explicit overflow menu

Bulk selection should exist, but only after entering an explicit `Select` mode.

That avoids the default checkbox-heavy experience.

### Activity flow

`Activity` should be a timeline of past recall sessions with persistent detail, not a list that turns into another page.

#### Layout

Use split view.

Left list:

1. recall question
2. timestamp
3. answer preview
4. evidence count

Right detail:

1. full question
2. full answer
3. evidence used
4. actions

Actions:

1. `Save as Quote`
2. `Repeat Question`
3. `Open Evidence in Library`

Delete belongs in an overflow menu or a low-emphasis danger zone.

Benefits:

1. no back-button state switch
2. easier comparison across sessions
3. history becomes meaningfully reviewable

### Settings flow

Split Settings into clear groups with different visual weight.

#### Group 1: Connection

Visible by default:

1. host
2. port
3. HTTPS
4. API key
5. model
6. `Test Connection`
7. `Fetch Models`

`Test Connection` should be more important than raw save.

#### Group 2: Retrieval

Visible by default:

1. max reference quotes
2. minimum relevance
3. concise explanation of quality vs recall tradeoff

#### Group 3: Personalization

Visible by default:

1. display name
2. theme

#### Group 4: Security

Collapsed by default:

1. change password

#### Group 5: Advanced

Collapsed by default:

1. web port
2. local storage paths
3. low-level runtime details

This keeps the page useful for normal users without hiding power-user controls.

## Modal Strategy

### Keep modals only for:

1. first-run name/setup gate
2. destructive confirmation
3. rare blocking errors

### Replace current modals with inline patterns

Replace quote editor modal with:

1. side drawer on desktop
2. full-page sheet on smaller screens

Replace import/export modal with:

1. guided panel or sheet
2. optional advanced disclosure for raw JSON preview

Replace success notices with:

1. toast
2. inline success banner near the originating action

## Import and Export UX

Import/export currently feels technical because the raw payload is front-and-center.

Target flow:

### Export

1. user selects one or more quotes
2. user chooses `Export`
3. app shows a small summary:
   count, authors, last updated
4. user saves/downloads the file
5. optional `Show JSON` disclosure for advanced users

### Import

1. user chooses file
2. app validates and previews what will happen
3. app shows summary:
   insert, update, duplicate, stale
4. user confirms import

Do not default to exposing the raw payload blob.

## Visual and Interaction Recommendations

### Shell

1. make the app shell calmer and more product-like
2. keep the current theme system
3. use stronger spatial hierarchy between navigation, page intro, content, and detail panes

### Lists

1. increase contrast between card title/content and metadata
2. stop rendering list rows as pseudo-terminal entries with `[x]` and `[1]`
3. use native checkbox visuals only in select mode

### Feedback

1. inline validation for forms
2. toasts for completed actions
3. persistent status only when the state still matters

### Empty states

Each primary destination should tell the user what to do next:

1. Ask: explain grounded recall and suggest adding/importing quotes
2. Library: prompt import or first quote creation
3. Activity: explain that past ask sessions will appear here

## Responsive Behavior

### Desktop

Use:

1. sidebar navigation
2. split panes
3. persistent detail panes

### Web

Use:

1. top nav or compact sidebar depending on width
2. stacked layout below tablet width
3. drawers instead of tiny modal dialogs

The web client should not feel like a constrained desktop clone.

## Implementation Priorities

Recommended order:

1. rename and reframe destinations: `Ask`, `Library`, `Activity`, `Settings`
2. redesign `Ask` to separate answer actions from quote-management actions
3. redesign `Library` as search + list-detail
4. redesign `Activity` as split-view timeline + detail
5. reduce modal usage
6. restructure `Settings` into basic and advanced groups

## Non-goals

This redesign does not change:

1. the local-first product model
2. manual file-based quote sharing
3. the shared Go core APIs as a first step

It changes presentation, interaction structure, and page flow.
