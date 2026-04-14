# Retrieval and Product Gap Closure

## Purpose

This plan tracks near-term improvements to the current product baseline. It is the main backlog for maturing the existing recall workflow before broader sync or network features are added.

## Active Gaps

- `SearchConfig.MinRelevance` is collected and saved but not applied in search queries
- test coverage is present in parts of the repo, but the project still needs stronger automated coverage across persistence, search, and UI flows
- the config directory is created but not meaningfully used for stored configuration
- CI is not yet established in the repository
- the desktop client exists as a scaffold, but the project still relies mainly on the TUI surface

## Near-Term Work

### Retrieval quality

- apply `MinRelevance` to the FTS query so saved search settings affect retrieval
- add better retrieval controls such as AND/OR strategies, tag filters, or score thresholds
- verify that search ranking and filtering remain stable as provenance and import features expand

### Product usability

- surface quote timestamps on the `Quotes` page
- continue tightening edit/delete/share/import flows across TUI and desktop
- improve error visibility for provider misconfiguration and failed model calls

### Quality and delivery

- add core tests for DB migrations, FTS search, settings persistence, and parsing fallbacks
- add integration tests around the add-search-answer and import-export paths
- add a GitHub Actions workflow for build and test

## Sequencing

### First

- retrieval relevance controls
- persistence and migration tests
- provider validation

### Next

- integration test coverage
- CI
- deeper retrieval controls

### Later

- broader client polish once the data and retrieval model stabilize
