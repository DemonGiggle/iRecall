# iRecall Roadmap

## Purpose

This file is the top-level roadmap and planning index for iRecall.

Use it for:

- project-wide priorities
- status across major initiatives
- links to detailed plans

Detailed execution plans live under [docs/plans/](./plans/README.md).

## Planning Structure

- [docs/plans/foundation-and-current-state.md](./plans/foundation-and-current-state.md)
  Historical baseline and implemented foundation.
- [docs/plans/retrieval-and-product-gap-closure.md](./plans/retrieval-and-product-gap-closure.md)
  Active backlog for improving the current recall workflow.
- [docs/plans/redmine-import-and-source-provenance.md](./plans/redmine-import-and-source-provenance.md)
  Source provenance model and Redmine import/export work.

Related design references:

- [docs/UI_DESIGN.md](./UI_DESIGN.md)
- [docs/QUOTES_SHARING_DESIGN.md](./QUOTES_SHARING_DESIGN.md)
- [docs/WAILS_DESKTOP.md](./WAILS_DESKTOP.md)
- [docs/SPEC.md](./SPEC.md)
- [docs/schema.md](./schema.md)

## Roadmap Summary

### Completed foundation

Status: established

The current codebase already has a working TUI, SQLite persistence, OpenAI-compatible provider integration, quote sharing/import, and a reusable Go core. The detailed baseline lives in [foundation-and-current-state.md](./plans/foundation-and-current-state.md).

### Product maturity

Status: active

The next layer of work is improving retrieval quality, coverage, validation, and product polish around the current note-search-answer workflow. The detailed backlog lives in [retrieval-and-product-gap-closure.md](./plans/retrieval-and-product-gap-closure.md).

### External source provenance and Redmine import

Status: in progress

The quote model now has first-class source identity fields and a Redmine export tool exists under `tools/redmine_export`. The remaining work is hardening the workflow, validating it against live data, and deciding how far import UX should go beyond the current share-envelope path. The detailed plan lives in [redmine-import-and-source-provenance.md](./plans/redmine-import-and-source-provenance.md).

## Near-Term Priorities

1. Strengthen retrieval and search correctness.
2. Improve automated coverage around persistence, imports, and recall flows.
3. Harden the generalized source provenance model as more import paths land.
4. Mature the Redmine export/import workflow on top of that source model.

## Medium-Term Priorities

1. Expand product polish across TUI and desktop surfaces.
2. Add stronger CI and delivery discipline.
3. Prepare the data model for future sync and discovery scenarios.

## Longer-Term Themes

1. External imports and synchronization.
2. Network-aware quote discovery and source filtering.
3. Additional client surfaces once the core model stabilizes.

## Non-Goals For Now

- embeddings or vector search
- a REST API
- remote-first storage
- full multi-user collaboration semantics
