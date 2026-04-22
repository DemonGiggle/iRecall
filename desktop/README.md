# Desktop Client

This directory contains the Wails desktop client for iRecall.

## Structure

```text
desktop/
├── main_wails.go     # Wails entrypoint (build-tagged)
├── wails.json        # Wails project configuration
└── README.md
```

## Design source

The desktop client must follow:

1. [docs/UI_DESIGN.md](../docs/UI_DESIGN.md)
2. [docs/WAILS_DESKTOP.md](../docs/WAILS_DESKTOP.md)

The desktop runtime binds the shared `app/` service and shared `frontend/` assets into a native Wails shell without inheriting terminal-specific presentation details.

## Current scope

The desktop app currently provides:

1. the same four top-level surfaces as the shared frontend shell: `Recall`, `History`, `Quotes`, and `Settings`
2. required startup name gating for quote attribution
3. quote add, edit, refine, inspect, delete, import, and export flows
4. grounded recall execution with reference quote actions and save-as-quote support
5. recall-history list/detail flows, history deletion, and save-as-quote actions
6. provider settings save and model fetch actions
7. debug/theme/web-password controls plus read-only storage-path visibility in `Settings`
8. Wails-native open/save dialogs for import and export

## Build note

`main_wails.go` is behind the `wails` build tag so the core repo can continue to build and test without requiring the Wails runtime for non-desktop work.

To build the desktop client directly from the repo:

1. install frontend dependencies in `frontend` with `npm install`
2. build the frontend bundle with `npm run build`
3. build the desktop target from the repo root with `go build -tags wails ./desktop/...`

You can also use the Wails CLI once it is installed in your environment.

## Current implementation notes

- the desktop client uses the same SQLite schema and engine behavior as the TUI
- the desktop and web clients share the same TypeScript frontend shell and Go application service
- import/export is still file-based and uses the same share envelope as the terminal client
- the frontend receives bootstrap metadata including product name, greeting, storage paths, and doc references
- settings include provider configuration, `MaxResults`, normalized `MinRelevance`, `MockLLM`, theme selection, web-password change, and the persisted web port
- the current desktop UI shows active storage paths, but storage-root switching is still exposed only in the TUI
