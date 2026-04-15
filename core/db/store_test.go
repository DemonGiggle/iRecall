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

	summaries, err := store.ListRecallHistory()
	if err != nil {
		t.Fatalf("list recall history: %v", err)
	}
	if len(summaries) != 1 {
		t.Fatalf("history count = %d, want 1", len(summaries))
	}
	if summaries[0].ID != historyID || summaries[0].Question != "How do goroutines coordinate?" {
		t.Fatalf("history summary = %+v", summaries[0])
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

	if err := store.DeleteRecallHistory([]int64{historyID}); err != nil {
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
