# Keyword Matching

This document describes how iRecall turns a recall question into a set of matching quotes today.

It is intentionally implementation-oriented. The goal is to describe the behavior that the current code actually executes, not an aspirational retrieval design.

## Scope

This applies to the current recall flow driven by:

- [core/engine.go:350](../core/engine.go)
- [core/db/store.go:264](../core/db/store.go)

It does not describe future semantic retrieval or embedding-based search. The current system is still keyword and FTS based.

## End-to-end flow

When the user asks a recall question, iRecall currently does this:

1. Send the question to the configured LLM and ask it to return `3` to `6` short lowercase search keywords.
2. Run a SQLite FTS5 search over quote content plus tag text.
3. Order the raw candidates by SQLite FTS rank.
4. Optionally apply an engine-side `MinRelevance` filter based on keyword coverage.
5. Trim the result list to `MaxResults`.
6. Pass the surviving quotes into grounded response generation.

In code:

- keyword extraction: [core/engine.go:350](../core/engine.go)
- candidate search: [core/engine.go:383](../core/engine.go)
- DB query: [core/db/store.go:264](../core/db/store.go)
- coverage score: [core/engine.go:417](../core/engine.go)

## 1. Keyword extraction

The engine asks the model for a JSON array of search keywords.

Current prompt contract:

- output must be JSON only
- output should be `3` to `6` short lowercase strings
- the keywords should be useful for searching a knowledge base

That behavior is defined in [core/engine.go:353](../core/engine.go).

Important consequences:

- Retrieval quality depends heavily on the keyword extractor.
- If the model chooses narrow or generic keywords, the search quality changes accordingly.
- iRecall does not currently expand synonyms or run alternate retrieval strategies after keyword extraction.

## 2. What is indexed

The FTS table stores:

- the quote content
- the quote tags as a single space-separated string

That update path is implemented in [core/db/store.go:125](../core/db/store.go).

Practical meaning:

- A quote can match because the keyword appears in the note body.
- A quote can also match because the keyword appears in one of the generated tags.
- Tags matter for retrieval, not just for display.

## 3. FTS query construction

The SQLite search path is in [core/db/store.go:266](../core/db/store.go).

Current behavior:

- each keyword is trimmed
- empty keywords are discarded
- each remaining keyword is wrapped in double quotes for FTS5
- embedded double quotes are escaped by doubling them
- the final MATCH expression is built as:

```text
"kw1" OR "kw2" OR "kw3"
```

This means the DB candidate set is an OR query, not an AND query.

Practical meaning:

- A quote can enter the candidate pool if it matches any one keyword.
- Broad keywords increase recall but also increase noise.
- Precision is improved later by the optional `MinRelevance` filter, not by the SQL query itself.

## 4. Candidate ordering

The SQL query orders results by `fts.rank` in SQLite FTS5:

- see [core/db/store.go:289](../core/db/store.go)

So the raw DB candidate list is ranked by FTS before any engine-side filtering happens.

This is the first-stage ranking only. If `MinRelevance` is enabled, low-coverage matches may still be removed after the DB returns them.

## 5. How `MinRelevance` works

The engine-side filtering logic is in [core/engine.go:383](../core/engine.go) and [core/engine.go:417](../core/engine.go).

When `MinRelevance == 0`:

- no coverage filtering is applied
- the FTS-ranked candidates are returned directly
- the final result count is capped by `MaxResults`

When `MinRelevance > 0`:

- the engine widens the DB fetch size to `max(MaxResults * 5, 25)`
- each returned quote receives a relevance score
- quotes below the threshold are dropped
- the survivors are trimmed back to `MaxResults`

### Relevance score formula

The current score is simple keyword coverage:

1. normalize keywords to lowercase
2. trim whitespace
3. deduplicate repeated keywords
4. build a lowercase haystack from:
   - quote content
   - quote tags
5. count how many normalized keywords appear via substring match
6. score = `matched_keywords / total_keywords`
7. round to 2 decimal places

So if the extracted keywords are:

```text
["raft", "leader", "quorum", "replication"]
```

and a quote contains `raft`, `leader`, and `replication`, then the score is:

```text
3 / 4 = 0.75
```

That quote passes a `MinRelevance` threshold of `0.7`, but fails a threshold of `0.8`.

## 6. Important behavior details

### Matching is not semantic

The current filter checks `strings.Contains(...)` on lowercase content and tags.

That means:

- it does not understand synonyms
- it does not understand paraphrases
- it does not use embeddings
- it does not use stemming beyond whatever SQLite FTS tokenizer already does for the DB stage

### The DB stage and the filter stage are different

There are two separate notions of “match”:

- DB candidate match: FTS5 `MATCH` against content + tags
- engine relevance match: substring keyword coverage across content + tags

This matters because a quote may:

- be returned by FTS
- then be filtered out by `MinRelevance`

### `MinRelevance` improves precision, not ranking sophistication

`MinRelevance` does not rerank quotes by score. It only removes quotes below a threshold.

After filtering, the remaining quotes keep their original FTS ordering.

## 7. Configuration impact

Relevant settings:

- `MaxResults`
- `MinRelevance`

Current guidance:

- `MinRelevance = 0.0` gives the broadest result set
- `0.3` to `0.7` is a practical range for cleaner matches
- `1.0` effectively requires full keyword coverage

The settings model is defined in [core/models.go:95](../core/models.go).

## 8. Current limitations

Known limitations of the current keyword matching model:

- keyword extraction quality can dominate retrieval quality
- OR-based DB matching can admit noisy candidates
- substring coverage can overcount partial textual matches
- results are not semantic and may miss conceptually related quotes
- the engine does not yet blend multiple retrieval strategies

## 9. Summary

Today, iRecall retrieval is:

- LLM-generated keyword extraction
- SQLite FTS5 candidate retrieval over content + tags
- optional normalized keyword-coverage filtering
- final truncation to `MaxResults`

That gives a simple, inspectable retrieval model, but it should be understood as keyword search with a lightweight coverage filter, not semantic search.
