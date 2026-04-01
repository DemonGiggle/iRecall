package db

import (
	"path/filepath"
	"testing"
)

func TestStoreQuoteLifecycleAndSearch(t *testing.T) {
	t.Parallel()

	store := openTestStore(t)

	quoteID, err := store.InsertQuote("Go channels coordinate concurrent goroutines.")
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
