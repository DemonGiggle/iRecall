# Desktop Client

This directory contains the Wails desktop client for iRecall.

## Structure

```text
desktop/
├── backend/          # Go application layer exposed to Wails
├── frontend/         # Wails frontend shell
├── main_wails.go     # Wails entrypoint (build-tagged)
├── wails.json        # Wails project configuration
└── README.md
```

## Design source

The desktop client must follow:

1. [docs/UI_DESIGN.md](/home/gigo/workspace/iRecall/docs/UI_DESIGN.md)
2. [docs/WAILS_DESKTOP.md](/home/gigo/workspace/iRecall/docs/WAILS_DESKTOP.md)

The backend in `desktop/backend` is intentionally separate from `tui/` so the desktop app can share the same product logic without inheriting terminal-specific presentation details.

## Current scope

The desktop app currently provides:

1. the same three top-level surfaces as the TUI: `Recall`, `Quotes`, and `Settings`
2. required startup name gating for quote attribution
3. quote add, edit, refine, delete, import, and export flows
4. grounded recall execution with reference quote actions
5. provider settings save and model fetch actions
6. Wails-native open/save dialogs for import and export

## Build note

`main_wails.go` is behind the `wails` build tag so the core repo can continue to build and test without requiring the Wails runtime for non-desktop work.

To build the desktop client directly from the repo:

1. install frontend dependencies in `desktop/frontend` with `npm install`
2. build the frontend bundle with `npm run build`
3. build the desktop target from the repo root with `go build -tags wails ./desktop/...`

You can also use the Wails CLI once it is installed in your environment.
