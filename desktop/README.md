# Desktop Client

This directory contains the Wails-oriented desktop client scaffold for iRecall.

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

The scaffold currently provides:

1. bootstrap state for a shell with `Recall`, `Quotes`, and `Settings`
2. quote CRUD wrappers
3. quote import/export wrappers
4. user profile and settings wrappers
5. a non-streaming recall wrapper for desktop integration
6. a minimal frontend shell that mirrors the current product structure

## Build note

`main_wails.go` is behind the `wails` build tag so the core repo can continue to build and test without the Wails module installed.

When you are ready to build the desktop client with the Wails toolchain, install the Wails dependencies and build the app from this scaffold.
