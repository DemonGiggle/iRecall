package core

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gigol/irecall/core/db"
)

func TestParseJSONStringArray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "plain json array",
			input: `["emmc","flash memory","partition"]`,
			want:  []string{"emmc", "flash memory", "partition"},
		},
		{
			name:  "markdown fenced json array",
			input: "```json\n[\"emmc\", \"flash memory\"]\n```",
			want:  []string{"emmc", "flash memory"},
		},
		{
			name:  "extra prose before array",
			input: "Here you go: [\"alpha\", \"beta\"]",
			want:  []string{"alpha", "beta"},
		},
		{
			name:  "comma fallback",
			input: `"Alpha", beta, gamma`,
			want:  []string{"alpha", "beta", "gamma"},
		},
		{
			name:    "empty response",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseJSONStringArray(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseJSONStringArray() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestRefineQuoteDraft(t *testing.T) {
	t.Parallel()

	var gotRequest struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
		Stream bool `json:"stream"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path = %q, want /v1/chat/completions", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotRequest); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{
					"message": map[string]string{
						"content": "A clearer version of the original note.",
					},
				},
			},
		})
	}))
	defer srv.Close()

	engine := newTestEngine(t, srv.Listener.Addr().String())

	refined, err := engine.RefineQuoteDraft(context.Background(), "messy draft note")
	if err != nil {
		t.Fatalf("RefineQuoteDraft() error = %v", err)
	}
	if refined != "A clearer version of the original note." {
		t.Fatalf("RefineQuoteDraft() = %q", refined)
	}
	if gotRequest.Model != "test-model" {
		t.Fatalf("model = %q, want test-model", gotRequest.Model)
	}
	if gotRequest.Stream {
		t.Fatal("stream = true, want false")
	}
	if len(gotRequest.Messages) != 2 {
		t.Fatalf("message count = %d, want 2", len(gotRequest.Messages))
	}
	if gotRequest.Messages[0].Role != "system" {
		t.Fatalf("system role = %q, want system", gotRequest.Messages[0].Role)
	}
	if gotRequest.Messages[1].Role != "user" {
		t.Fatalf("user role = %q, want user", gotRequest.Messages[1].Role)
	}
	if gotRequest.Messages[1].Content == "" || gotRequest.Messages[1].Content == "messy draft note" {
		t.Fatalf("unexpected user prompt content = %q", gotRequest.Messages[1].Content)
	}
}

func TestExtractTagsRequestsBroaderTagSet(t *testing.T) {
	t.Parallel()

	var gotRequest struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
		MaxTokens int  `json:"max_tokens"`
		Stream    bool `json:"stream"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path = %q, want /v1/chat/completions", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotRequest); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{
					"message": map[string]string{
						"content": `["distributed systems","consensus","raft","leader election","quorum","replication","fault tolerance"]`,
					},
				},
			},
		})
	}))
	defer srv.Close()

	engine := newTestEngine(t, srv.Listener.Addr().String())

	tags, err := engine.ExtractTags(context.Background(), "Raft relies on leader election, replication, and quorum to keep a distributed system consistent during failures.")
	if err != nil {
		t.Fatalf("ExtractTags() error = %v", err)
	}
	if len(tags) != 7 {
		t.Fatalf("tag count = %d, want 7", len(tags))
	}
	if gotRequest.Model != "test-model" {
		t.Fatalf("model = %q, want test-model", gotRequest.Model)
	}
	if gotRequest.Stream {
		t.Fatal("stream = true, want false")
	}
	if gotRequest.MaxTokens != 384 {
		t.Fatalf("max_tokens = %d, want 384", gotRequest.MaxTokens)
	}
	if len(gotRequest.Messages) != 2 {
		t.Fatalf("message count = %d, want 2", len(gotRequest.Messages))
	}
	if !strings.Contains(gotRequest.Messages[0].Content, "For short or simple text, return fewer tags when appropriate.") {
		t.Fatalf("system prompt = %q, want flexible short-quote guidance", gotRequest.Messages[0].Content)
	}
}

func TestExtractTagsRepairsWeakGenericResults(t *testing.T) {
	t.Parallel()

	var requestCount int
	var gotRequests []struct {
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
		MaxTokens int `json:"max_tokens"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		var req struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
			MaxTokens int `json:"max_tokens"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request %d: %v", requestCount, err)
		}
		gotRequests = append(gotRequests, req)

		w.Header().Set("Content-Type", "application/json")
		content := `["quote","note","text","raft"]`
		if requestCount == 2 {
			content = `["raft","consensus","leader election","replication","quorum","fault tolerance"]`
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{
					"message": map[string]string{
						"content": content,
					},
				},
			},
		})
	}))
	defer srv.Close()

	engine := newTestEngine(t, srv.Listener.Addr().String())

	tags, err := engine.ExtractTags(context.Background(), "Raft relies on leader election, replication, and quorum to keep a distributed system consistent during failures.")
	if err != nil {
		t.Fatalf("ExtractTags() error = %v", err)
	}
	want := []string{"raft", "consensus", "leader election", "replication", "quorum", "fault tolerance"}
	if !reflect.DeepEqual(tags, want) {
		t.Fatalf("ExtractTags() = %#v, want %#v", tags, want)
	}
	if requestCount != 2 {
		t.Fatalf("request count = %d, want 2", requestCount)
	}
	if !strings.Contains(gotRequests[1].Messages[0].Content, "repairing a JSON keyword extractor result") {
		t.Fatalf("repair prompt = %q, want repair instructions", gotRequests[1].Messages[0].Content)
	}
}

func TestSearchQuotesAppliesMinRelevance(t *testing.T) {
	t.Parallel()

	engine := newProfiledTestEngine(t, "localhost")
	engine.cfg.Search.MaxResults = 5
	engine.cfg.Search.MinRelevance = 0.6

	identity := db.QuoteIdentity{
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
	}
	id1, err := engine.store.InsertQuote("Go concurrency patterns with channels and goroutines.", identity)
	if err != nil {
		t.Fatalf("InsertQuote() first error = %v", err)
	}
	if err := engine.store.UpdateQuoteFTS(id1, []string{"go", "concurrency"}); err != nil {
		t.Fatalf("UpdateQuoteFTS() first error = %v", err)
	}

	identity.GlobalID = "quote-2"
	identity.SourceEntityID = "quote-2"
	id2, err := engine.store.InsertQuote("Go testing basics.", identity)
	if err != nil {
		t.Fatalf("InsertQuote() second error = %v", err)
	}
	if err := engine.store.UpdateQuoteFTS(id2, []string{"go", "testing"}); err != nil {
		t.Fatalf("UpdateQuoteFTS() second error = %v", err)
	}

	got, err := engine.SearchQuotes(context.Background(), []string{"go", "concurrency"})
	if err != nil {
		t.Fatalf("SearchQuotes() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("SearchQuotes() count = %d, want 1", len(got))
	}
	if got[0].GlobalID != "quote-1" {
		t.Fatalf("SearchQuotes() returned %q, want quote-1", got[0].GlobalID)
	}
}

func TestNormalizeTagsFiltersNoiseAndCapsCount(t *testing.T) {
	t.Parallel()

	input := []string{
		"Quote", "raft", "RAFT", "text", "leader election", " ", "a",
		"consensus", "replication", "quorum", "fault tolerance",
		"distributed systems", "log replication", "state machine", "term", "commit index",
		"heartbeat", "candidate", "follower", "leader", "safety", "liveness", "cluster",
		"election timeout", "append entries", "majority", "durability", "availability",
		"partition tolerance", "recovery", "failover", "membership changes", "snapshotting",
	}

	got, stats := normalizeTags(input)
	if len(got) > maxExtractedTags {
		t.Fatalf("tag count = %d, want <= %d", len(got), maxExtractedTags)
	}
	if stats.RemovedGeneric == 0 {
		t.Fatal("expected generic tags to be removed")
	}
	if stats.RemovedDupes == 0 {
		t.Fatal("expected duplicate tags to be removed")
	}
	if got[0] != "raft" {
		t.Fatalf("first tag = %q, want raft", got[0])
	}
	if slices.Contains(got, "quote") || slices.Contains(got, "text") {
		t.Fatalf("normalizeTags() kept generic tags: %#v", got)
	}
}

func TestExportQuotesIncludesSchemaVersion(t *testing.T) {
	t.Parallel()

	engine := newProfiledTestEngine(t, "")
	q, err := engine.AddQuote(context.Background(), "share me")
	if err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}

	payload, err := engine.ExportQuotes(context.Background(), []int64{q.ID})
	if err != nil {
		t.Fatalf("ExportQuotes() error = %v", err)
	}

	var env SharedQuoteEnvelope
	if err := json.Unmarshal(payload, &env); err != nil {
		t.Fatalf("unmarshal export: %v", err)
	}
	if env.SchemaVersion != ShareSchemaVersion {
		t.Fatalf("schema version = %d, want %d", env.SchemaVersion, ShareSchemaVersion)
	}
	if len(env.Quotes) != 1 {
		t.Fatalf("export quote count = %d, want 1", len(env.Quotes))
	}
	if env.Quotes[0].GlobalID == "" {
		t.Fatal("exported quote missing global id")
	}
	if env.Quotes[0].SourceBackend != "local" || env.Quotes[0].SourceEntityID == "" {
		t.Fatalf("exported source provenance = %+v, want local provenance fields", env.Quotes[0])
	}
}

func TestImportSharedQuotesRejectsUnsupportedSchema(t *testing.T) {
	t.Parallel()

	engine := newProfiledTestEngine(t, "")
	payload, err := json.Marshal(SharedQuoteEnvelope{
		SchemaVersion: ShareSchemaVersion + 1,
		ExportedAt:    time.Now().UTC(),
		Quotes: []SharedQuoteEntry{
			{
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
				Content:          "hello",
				CreatedAtUTC:     time.Now().UTC(),
				UpdatedAtUTC:     time.Now().UTC(),
			},
		},
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	if _, err := engine.ImportSharedQuotes(context.Background(), payload); err == nil {
		t.Fatal("expected unsupported schema version to fail")
	}
}

func TestImportSharedQuotesVersionHandling(t *testing.T) {
	t.Parallel()

	engine := newProfiledTestEngine(t, "")
	base := SharedQuoteEnvelope{
		SchemaVersion: ShareSchemaVersion,
		ExportedAt:    time.Now().UTC(),
		Quotes: []SharedQuoteEntry{
			{
				GlobalID:         "quote-1",
				AuthorUserID:     "author-1",
				AuthorName:       "Alice",
				SourceUserID:     "author-1",
				SourceName:       "Alice",
				SourceBackend:    "redmine",
				SourceNamespace:  "redmine:testdb",
				SourceEntityType: "issue_description",
				SourceEntityID:   "123",
				SourceLabel:      "Redmine issue #123",
				SourceURL:        "https://redmine.test/issues/123",
				Version:          1,
				Content:          "first version",
				Tags:             []string{"alpha"},
				CreatedAtUTC:     time.Unix(100, 0).UTC(),
				UpdatedAtUTC:     time.Unix(100, 0).UTC(),
			},
		},
	}

	payload, _ := json.Marshal(base)
	result, err := engine.ImportSharedQuotes(context.Background(), payload)
	if err != nil {
		t.Fatalf("initial import error = %v", err)
	}
	if result.Inserted != 1 {
		t.Fatalf("initial import inserted = %d, want 1", result.Inserted)
	}

	result, err = engine.ImportSharedQuotes(context.Background(), payload)
	if err != nil {
		t.Fatalf("duplicate import error = %v", err)
	}
	if result.Duplicates != 1 {
		t.Fatalf("duplicate import duplicates = %d, want 1", result.Duplicates)
	}

	newer := SharedQuoteEnvelope{
		SchemaVersion: ShareSchemaVersion,
		ExportedAt:    time.Now().UTC(),
		Quotes: []SharedQuoteEntry{
			{
				GlobalID:         "quote-1",
				AuthorUserID:     "author-1",
				AuthorName:       "Alice",
				SourceUserID:     "author-1",
				SourceName:       "Alice",
				SourceBackend:    "redmine",
				SourceNamespace:  "redmine:testdb",
				SourceEntityType: "issue_description",
				SourceEntityID:   "123",
				SourceLabel:      "Redmine issue #123",
				SourceURL:        "https://redmine.test/issues/123",
				Version:          2,
				Content:          "second version",
				Tags:             []string{"beta"},
				CreatedAtUTC:     time.Unix(100, 0).UTC(),
				UpdatedAtUTC:     time.Unix(200, 0).UTC(),
			},
		},
	}
	newerPayload, _ := json.Marshal(newer)
	result, err = engine.ImportSharedQuotes(context.Background(), newerPayload)
	if err != nil {
		t.Fatalf("newer import error = %v", err)
	}
	if result.Updated != 1 {
		t.Fatalf("newer import updated = %d, want 1", result.Updated)
	}

	stale := SharedQuoteEnvelope{
		SchemaVersion: ShareSchemaVersion,
		ExportedAt:    time.Now().UTC(),
		Quotes: []SharedQuoteEntry{
			{
				GlobalID:         "quote-1",
				AuthorUserID:     "author-1",
				AuthorName:       "Alice",
				SourceUserID:     "author-1",
				SourceName:       "Alice",
				SourceBackend:    "redmine",
				SourceNamespace:  "redmine:testdb",
				SourceEntityType: "issue_description",
				SourceEntityID:   "123",
				SourceLabel:      "Redmine issue #123",
				SourceURL:        "https://redmine.test/issues/123",
				Version:          1,
				Content:          "first version",
				Tags:             []string{"alpha"},
				CreatedAtUTC:     time.Unix(100, 0).UTC(),
				UpdatedAtUTC:     time.Unix(100, 0).UTC(),
			},
		},
	}
	stalePayload, _ := json.Marshal(stale)
	result, err = engine.ImportSharedQuotes(context.Background(), stalePayload)
	if err != nil {
		t.Fatalf("stale import error = %v", err)
	}
	if result.Stale != 1 {
		t.Fatalf("stale import stale = %d, want 1", result.Stale)
	}

	quotes, err := engine.ListQuotes(context.Background())
	if err != nil {
		t.Fatalf("ListQuotes() error = %v", err)
	}
	if len(quotes) != 1 {
		t.Fatalf("quote count = %d, want 1", len(quotes))
	}
	if quotes[0].Content != "second version" || quotes[0].Version != 2 {
		t.Fatalf("quote after update = %+v", quotes[0])
	}
	if quotes[0].SourceBackend != "redmine" || quotes[0].SourceEntityID != "123" {
		t.Fatalf("quote source provenance = %+v", quotes[0])
	}
}

func newTestEngine(t *testing.T, host string) *Engine {
	t.Helper()

	store, err := db.Open(filepath.Join(t.TempDir(), "engine.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	settings := DefaultSettings()
	settings.Provider = ProviderConfig{
		Host:  host,
		Port:  0,
		HTTPS: false,
		Model: "test-model",
	}
	return New(store, settings)
}

func newProfiledTestEngine(t *testing.T, host string) *Engine {
	t.Helper()
	engine := newTestEngine(t, host)
	profile := &UserProfile{
		UserID:      "local-user",
		DisplayName: "Local User",
		CreatedAt:   time.Unix(1, 0),
		UpdatedAt:   time.Unix(1, 0),
	}
	if err := engine.SaveUserProfile(context.Background(), profile); err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}
	engine.UpdateUserProfile(profile)
	return engine
}
