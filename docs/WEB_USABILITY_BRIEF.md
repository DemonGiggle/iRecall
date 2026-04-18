# iRecall Web Usability Brief

## Purpose

This brief translates the current iRecall product into a simpler web experience.

Primary goal:

1. make the web version easy enough that an elementary school student can use it

Secondary goal:

1. preserve the real product behavior implemented by the shared Go core

This document is meant to support a future Figma design pass and frontend refactor.

## Current Product Reality

iRecall is not a generic chat app. It is a local-first note recall tool with four product surfaces:

1. `Recall`
2. `History`
3. `Quotes`
4. `Settings`

Key user flow:

1. save notes or quotes
2. ask a question
3. extract keywords
4. search local notes
5. generate an answer grounded only in retrieved notes
6. optionally save that question/answer as a quote

Core references:

1. `core/engine.go`
2. `app/app.go`
3. `web/server.go`
4. `frontend/src/app.ts`
5. `docs/UI_DESIGN.md`

## Codebase Design-System Reality

The current frontend is intentionally lightweight and does not yet have a formal component library.

### Frameworks and runtime

1. frontend runtime: TypeScript + Vite
2. frontend rendering approach: manual DOM rendering from `frontend/src/app.ts`
3. backend: Go
4. web transport: JSON HTTP endpoints exposed by `web/server.go`

Relevant files:

1. `frontend/package.json`
2. `frontend/src/main.ts`
3. `frontend/src/app.ts`
4. `frontend/src/styles.css`
5. `frontend/src/theme.ts`

### Tokens

There is a small token-like theme layer in `frontend/src/theme.ts`.

Current semantic tokens:

1. `bg`
2. `bgStrong`
3. `panel`
4. `panel2`
5. `border`
6. `borderStrong`
7. `primary`
8. `accent`
9. `muted`
10. `fg`
11. `ok`
12. `error`
13. `shadow`

Themes:

1. `violet`
2. `forest`
3. `sunset`
4. `ocean`
5. `paper`

This is useful because the product already thinks in semantic roles, but it is not yet a full design token system.

### Components

There is no standalone component directory.

Instead, the current UI is built from:

1. render helper functions in `frontend/src/app.ts`
2. shared CSS classes in `frontend/src/styles.css`

Current reusable UI patterns are mostly:

1. tabs
2. panels
3. subpanels
4. toolbars
5. buttons
6. text inputs
7. text areas
8. quote cards
9. modals
10. status banners

### Styling approach

The frontend uses one global stylesheet: `frontend/src/styles.css`.

Styling characteristics:

1. CSS custom properties for theme values
2. shared utility-like classes
3. responsive behavior through a few media queries
4. no CSS modules
5. no Tailwind
6. no CSS-in-JS

### Assets

The built frontend is bundled into Go via `frontend/assets.go`.

Important implication:

1. the web client is currently a self-contained app shell
2. a redesign can stay frontend-only as long as the API contract remains stable

### Icons and illustration

There is no dedicated icon system in the current repo.

For the simplified web redesign, that means we should prefer:

1. a very small icon set
2. simple visual metaphors
3. large labeled buttons instead of icon-only controls

## Product Constraints We Must Keep

The redesign should simplify presentation, not change the underlying product contract.

Must keep:

1. explicit name capture before normal use
2. password login for web
3. quote ownership and source attribution
4. explicit add, edit, delete, import, and export flows
5. recall answers grounded only in retrieved notes
6. history of past recalls
7. advanced provider and search settings somewhere in the product

Can simplify heavily:

1. navigation labels
2. information density
3. number of simultaneous actions shown on screen
4. wording
5. visual hierarchy
6. first-run onboarding

## Usability Direction

### Core principle

The current UI feels like a tool for technical users.

The new web version should feel like:

1. ask a question
2. see the answer
3. check the notes that helped
4. add more notes when needed

### Mental model

Use child-friendly wording.

Suggested label translations:

1. `Recall` -> `Ask`
2. `Quotes` -> `My Notes`
3. `History` -> `Past Questions`
4. `Settings` -> `Setup`
5. `Save as Quote` -> `Save Answer`
6. `Import Quotes` -> `Bring In Notes`
7. `Export Quotes` -> `Share Notes`
8. `Min Relevance` -> `How close the match should be`

### Navigation

The current four-surface architecture can remain, but the visible emphasis should change.

Recommended top-level order:

1. `Ask`
2. `My Notes`
3. `Past Questions`
4. `Setup`

Recommended default landing screen:

1. `Ask`

### Interaction simplification

The new web UI should:

1. show one primary action per screen
2. hide advanced actions behind an `More` or secondary section
3. reduce bulk-selection emphasis in the default view
4. make the happy path pointer-first
5. still preserve keyboard accessibility

## Recommended Screen Model

### 1. Ask

This should become the hero screen.

Primary content:

1. one large question box
2. one clear `Ask` button
3. answer card
4. notes used card list

Child-friendly behavior:

1. the question box should encourage natural language
2. helper text should use examples
3. evidence cards should look like friendly note cards

Suggested helper copy:

1. `Ask something like: "What did I learn about SQLite?"`

### 2. My Notes

This should feel like a simple notebook.

Primary content:

1. add note button
2. search or filter box later if needed
3. large readable note cards
4. tags shown softly, not as the star of the page

Default actions on each card:

1. open
2. edit
3. delete

Secondary action:

1. share

### 3. Past Questions

This should feel like a timeline of previous learning.

Primary content:

1. question cards
2. short answer preview
3. open to see full answer and the notes that supported it

### 4. Setup

This should be split into:

1. everyday settings
2. advanced setup

Everyday settings:

1. display name
2. theme
3. password change

Advanced setup:

1. provider host
2. provider port
3. HTTPS
4. API key
5. model selection
6. search tuning
7. web port

## Visual Direction

### Tone

The UI should feel:

1. friendly
2. calm
3. obvious
4. safe
5. encouraging

It should not feel:

1. dense
2. terminal-like
3. admin-heavy
4. overly clever

### Visual recommendations

1. use larger cards and buttons than the current UI
2. keep strong contrast
3. use fewer simultaneous panels
4. use more whitespace between major tasks
5. use a warm, approachable palette by default
6. keep error colors reserved for true problems
7. use success states generously for completed actions

### Typography

The current app uses `IBM Plex Sans` and works well for clarity.

For the redesign:

1. keep a highly legible sans serif
2. increase body and button size
3. avoid tiny muted helper text

## Suggested Figma Frames

When the writable Figma workflow is available, start with these desktop web frames:

1. `Ask - Empty`
2. `Ask - Answer Ready`
3. `My Notes - List`
4. `My Notes - Add Note`
5. `Past Questions - List`
6. `Past Questions - Detail`
7. `Setup - Basic`
8. `Setup - Advanced`
9. `Login`
10. `First Run - Set Your Name`

## Suggested Figma Components

1. primary button
2. secondary button
3. note card
4. answer card
5. question history card
6. top navigation tab
7. empty state card
8. modal dialog
9. status banner
10. simple tag chip

## First Figma Pass Priorities

If only one workflow gets polished first, it should be:

1. login
2. set name
3. ask a question
4. read the answer
5. inspect the supporting notes
6. save the answer

That flow is the clearest expression of the product.

## Implementation Notes For Later

When converting the redesign back into code, prefer:

1. preserving the existing API endpoints
2. splitting `frontend/src/app.ts` into smaller view modules
3. keeping semantic theme tokens
4. introducing a tiny component layer before a full framework rewrite

The current frontend already has enough structure to support a redesign without changing the Go backend.
