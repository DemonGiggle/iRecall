package db

import (
	"path/filepath"
	"testing"
)

func TestStoreQuoteLifecycleAndSearch(t *testing.T) {
	t.Parallel()

	store := openTestStore(t)

	quoteID, err := store.InsertQuote("Go channels coordinate concurrent goroutines.", QuoteIdentity{
		GlobalID:         "quote-1",
		AuthorUserID:     "user-1",
		AuthorName:       "Alice",
		SourceUserID:     "user-1",
		SourceName:       "Alice",
		SourceBackend:    "local",
		SourceNamespace:  "local:user-1",
		SourceEntityType: "quote",
		SourceEntityID:   "quote-1",
		SourceLabel:      "Local quote",
		Version:          1,
	})
	if err != nil {
		t.Fatalf("insert quote: %v", err)
	}

	tagIDs, err := store.UpsertTags([]string{"concurrency", "golang"})
	if err != nil {
		t.Fatalf("upsert tags: %v", err)
	}
	if len(tagIDs) != 2 {
		t.Fatalf("tag id count = %d, want 2", len(tagIDs))
	}

	if err := store.InsertQuoteTags(quoteID, tagIDs); err != nil {
		t.Fatalf("insert quote tags: %v", err)
	}
	if err := store.UpdateQuoteFTS(quoteID, []string{"concurrency", "golang"}); err != nil {
		t.Fatalf("update quote fts: %v", err)
	}

	results, err := store.SearchQuotes([]string{"concurrency"}, 5)
	if err != nil {
		t.Fatalf("search quotes by tag: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("search result count = %d, want 1", len(results))
	}
	if results[0].ID != quoteID {
		t.Fatalf("search returned quote id %d, want %d", results[0].ID, quoteID)
	}

	listed, err := store.ListQuotes()
	if err != nil {
		t.Fatalf("list quotes: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("listed quote count = %d, want 1", len(listed))
	}
	if listed[0].GlobalID != "quote-1" {
		t.Fatalf("global id = %q, want quote-1", listed[0].GlobalID)
	}
	if listed[0].SourceBackend != "local" || listed[0].SourceEntityID != "quote-1" {
		t.Fatalf("source provenance = %+v, want local provenance", listed[0])
	}
	if listed[0].Tags != "concurrency,golang" && listed[0].Tags != "golang,concurrency" {
		t.Fatalf("listed tags = %q, want concurrency and golang", listed[0].Tags)
	}

	count, err := store.CountQuotes()
	if err != nil {
		t.Fatalf("count quotes: %v", err)
	}
	if count != 1 {
		t.Fatalf("quote count = %d, want 1", count)
	}

	if err := store.DeleteQuote(quoteID); err != nil {
		t.Fatalf("delete quote: %v", err)
	}

	results, err = store.SearchQuotes([]string{"concurrency"}, 5)
	if err != nil {
		t.Fatalf("search after delete: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("search result count after delete = %d, want 0", len(results))
	}

	count, err = store.CountQuotes()
	if err != nil {
		t.Fatalf("count quotes after delete: %v", err)
	}
	if count != 0 {
		t.Fatalf("quote count after delete = %d, want 0", count)
	}
}

func TestApplyQuoteTagUpdateCanTouchQuoteAtomically(t *testing.T) {
	t.Parallel()

	store := openTestStore(t)

	quoteID, err := store.InsertQuote("SQLite WAL keeps readers moving.", QuoteIdentity{
		GlobalID:         "quote-touch",
		AuthorUserID:     "user-1",
		AuthorName:       "Alice",
		SourceUserID:     "user-1",
		SourceName:       "Alice",
		SourceBackend:    "local",
		SourceNamespace:  "local:user-1",
		SourceEntityType: "quote",
		SourceEntityID:   "quote-touch",
		SourceLabel:      "Local quote",
		Version:          1,
	})
	if err != nil {
		t.Fatalf("insert quote: %v", err)
	}

	if err := store.ApplyQuoteTagUpdate(quoteID, []string{"sqlite", "wal"}, true); err != nil {
		t.Fatalf("ApplyQuoteTagUpdate() error = %v", err)
	}

	results, err := store.SearchQuotes([]string{"wal"}, 5)
	if err != nil {
		t.Fatalf("search quotes by updated tag: %v", err)
	}
	if len(results) != 1 || results[0].ID != quoteID {
		t.Fatalf("SearchQuotes() = %+v, want quote id %d", results, quoteID)
	}

	listed, err := store.ListQuotes()
	if err != nil {
		t.Fatalf("list quotes: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("listed quote count = %d, want 1", len(listed))
	}
	if listed[0].Version != 2 {
		t.Fatalf("quote version = %d, want 2", listed[0].Version)
	}
	if listed[0].Tags != "sqlite,wal" && listed[0].Tags != "wal,sqlite" {
		t.Fatalf("listed tags = %q, want sqlite and wal", listed[0].Tags)
	}
}

func TestStoreSettingsRoundTrip(t *testing.T) {
	t.Parallel()

	store := openTestStore(t)

	if err := store.SetSetting("settings", `{"provider":{"host":"localhost"}}`); err != nil {
		t.Fatalf("set setting: %v", err)
	}

	got, err := store.GetSetting("settings")
	if err != nil {
		t.Fatalf("get setting: %v", err)
	}
	if got != `{"provider":{"host":"localhost"}}` {
		t.Fatalf("setting value = %q", got)
	}
}

func TestStoreUserProfileRoundTripAndBackfill(t *testing.T) {
	t.Parallel()

	store := openTestStore(t)
	if _, err := store.InsertQuote("legacy quote", QuoteIdentity{}); err != nil {
		t.Fatalf("insert legacy quote: %v", err)
	}

	profile := UserProfileRow{
		UserID:      "user-1",
		DisplayName: "Alice",
		CreatedAt:   100,
		UpdatedAt:   100,
	}
	if err := store.SaveUserProfile(profile); err != nil {
		t.Fatalf("save user profile: %v", err)
	}

	got, err := store.GetUserProfile()
	if err != nil {
		t.Fatalf("get user profile: %v", err)
	}
	if got.UserID != profile.UserID || got.DisplayName != profile.DisplayName {
		t.Fatalf("user profile = %+v, want %+v", got, profile)
	}

	if err := store.BackfillQuoteIdentity(profile.UserID, profile.DisplayName, 200, func() string { return "uuid-1" }); err != nil {
		t.Fatalf("backfill quote identity: %v", err)
	}

	quotes, err := store.ListQuotes()
	if err != nil {
		t.Fatalf("list quotes: %v", err)
	}
	if len(quotes) != 1 {
		t.Fatalf("quote count = %d, want 1", len(quotes))
	}
	if quotes[0].GlobalID != "uuid-1" || quotes[0].AuthorName != "Alice" || quotes[0].SourceName != "Alice" {
		t.Fatalf("backfilled quote = %+v", quotes[0])
	}
}

func TestStoreRecallHistoryLifecycle(t *testing.T) {
	t.Parallel()

	store := openTestStore(t)

	firstQuoteID, err := store.InsertQuote("First history quote", QuoteIdentity{
		GlobalID:         "quote-1",
		AuthorUserID:     "user-1",
		AuthorName:       "Alice",
		SourceUserID:     "user-1",
		SourceName:       "Alice",
		SourceBackend:    "local",
		SourceNamespace:  "local:user-1",
		SourceEntityType: "quote",
		SourceEntityID:   "quote-1",
		SourceLabel:      "Local quote",
		Version:          1,
	})
	if err != nil {
		t.Fatalf("insert first quote: %v", err)
	}
	secondQuoteID, err := store.InsertQuote("Second history quote", QuoteIdentity{
		GlobalID:         "quote-2",
		AuthorUserID:     "user-1",
		AuthorName:       "Alice",
		SourceUserID:     "user-1",
		SourceName:       "Alice",
		SourceBackend:    "local",
		SourceNamespace:  "local:user-1",
		SourceEntityType: "quote",
		SourceEntityID:   "quote-2",
		SourceLabel:      "Local quote",
		Version:          1,
	})
	if err != nil {
		t.Fatalf("insert second quote: %v", err)
	}

	historyID, err := store.InsertRecallHistory("How do goroutines coordinate?", "Use channels to synchronize work.", []int64{firstQuoteID, secondQuoteID})
	if err != nil {
		t.Fatalf("insert recall history: %v", err)
	}
	secondHistoryID, err := store.InsertRecallHistory("How do I page through history?", "Use limit and offset.", []int64{secondQuoteID})
	if err != nil {
		t.Fatalf("insert second recall history: %v", err)
	}

	summaries, err := store.ListRecallHistory()
	if err != nil {
		t.Fatalf("list recall history: %v", err)
	}
	if len(summaries) != 2 {
		t.Fatalf("history count = %d, want 2", len(summaries))
	}
	if summaries[0].ID != secondHistoryID || summaries[0].Question != "How do I page through history?" {
		t.Fatalf("newest history summary = %+v", summaries[0])
	}
	if summaries[1].ID != historyID || summaries[1].Question != "How do goroutines coordinate?" {
		t.Fatalf("older history summary = %+v", summaries[1])
	}

	count, err := store.CountRecallHistory()
	if err != nil {
		t.Fatalf("count recall history: %v", err)
	}
	if count != 2 {
		t.Fatalf("recall history count = %d, want 2", count)
	}

	page, err := store.ListRecallHistoryPage(1, 0)
	if err != nil {
		t.Fatalf("list recall history page 1: %v", err)
	}
	if len(page) != 1 || page[0].ID != secondHistoryID {
		t.Fatalf("page 1 = %+v, want newest entry %d", page, secondHistoryID)
	}

	page, err = store.ListRecallHistoryPage(1, 1)
	if err != nil {
		t.Fatalf("list recall history page 2: %v", err)
	}
	if len(page) != 1 || page[0].ID != historyID {
		t.Fatalf("page 2 = %+v, want older entry %d", page, historyID)
	}

	entry, err := store.GetRecallHistory(historyID)
	if err != nil {
		t.Fatalf("get recall history: %v", err)
	}
	if entry.Response != "Use channels to synchronize work." {
		t.Fatalf("history response = %q, want exact saved response", entry.Response)
	}
	if len(entry.Quotes) != 2 {
		t.Fatalf("history quote count = %d, want 2", len(entry.Quotes))
	}
	if entry.Quotes[0].ID != firstQuoteID || entry.Quotes[1].ID != secondQuoteID {
		t.Fatalf("history quote order = [%d %d], want [%d %d]", entry.Quotes[0].ID, entry.Quotes[1].ID, firstQuoteID, secondQuoteID)
	}

	if err := store.DeleteRecallHistory([]int64{historyID, secondHistoryID}); err != nil {
		t.Fatalf("delete recall history: %v", err)
	}

	summaries, err = store.ListRecallHistory()
	if err != nil {
		t.Fatalf("list recall history after delete: %v", err)
	}
	if len(summaries) != 0 {
		t.Fatalf("history count after delete = %d, want 0", len(summaries))
	}
}

func openTestStore(t *testing.T) *Store {
	t.Helper()

	path := filepath.Join(t.TempDir(), "store.db")
	store, err := Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() {
		store.Close()
	})
	return store
}
