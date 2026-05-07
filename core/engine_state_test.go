package core

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gigol/irecall/core/db"
)

func TestProviderConnectivityAndModelDiscovery(t *testing.T) {
	t.Parallel()

	var authHeaders []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			t.Fatalf("path = %q, want /v1/models", r.URL.Path)
		}
		authHeaders = append(authHeaders, r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]string{
				{"id": "model-b"},
				{"id": "model-a"},
			},
		})
	}))
	defer srv.Close()

	engine := newTestEngine(t, "localhost")
	cfg := ProviderConfig{
		Host:   srv.Listener.Addr().String(),
		Model:  "runtime-model",
		APIKey: "runtime-secret",
	}

	models, err := engine.FetchModels(context.Background(), cfg)
	if err != nil {
		t.Fatalf("FetchModels() error = %v", err)
	}
	if want := []string{"model-a", "model-b"}; !reflect.DeepEqual(models, want) {
		t.Fatalf("FetchModels() = %#v, want %#v", models, want)
	}

	if err := engine.TestProvider(context.Background(), cfg); err != nil {
		t.Fatalf("TestProvider() error = %v", err)
	}

	if len(authHeaders) != 2 {
		t.Fatalf("auth header count = %d, want 2", len(authHeaders))
	}
	for _, got := range authHeaders {
		if got != "Bearer runtime-secret" {
			t.Fatalf("Authorization = %q, want Bearer runtime-secret", got)
		}
	}
}

func TestProviderConnectivityErrors(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}))
	defer srv.Close()

	engine := newTestEngine(t, "localhost")
	cfg := ProviderConfig{Host: srv.Listener.Addr().String(), Model: "runtime-model"}

	if _, err := engine.FetchModels(context.Background(), cfg); err == nil || !strings.Contains(err.Error(), "provider returned 502") {
		t.Fatalf("FetchModels() error = %v, want provider 502 failure", err)
	}
	if err := engine.TestProvider(context.Background(), cfg); err == nil || !strings.Contains(err.Error(), "provider returned 502") {
		t.Fatalf("TestProvider() error = %v, want provider 502 failure", err)
	}
}

func TestLoadAndSaveSettingsLifecycle(t *testing.T) {
	t.Parallel()

	engine := newTestEngine(t, "localhost")

	defaults := DefaultSettings()
	loaded, err := engine.LoadSettings(context.Background())
	if err != nil {
		t.Fatalf("LoadSettings() initial error = %v", err)
	}
	if !reflect.DeepEqual(loaded, defaults) {
		t.Fatalf("LoadSettings() initial = %+v, want defaults %+v", loaded, defaults)
	}

	invalid := *defaults
	invalid.Web.Port = 70000
	if err := engine.SaveSettings(context.Background(), &invalid); err == nil || !strings.Contains(err.Error(), "web port") {
		t.Fatalf("SaveSettings(invalid port) error = %v, want validation failure", err)
	}

	next := *defaults
	next.Theme = ""
	next.Web.Port = 9988
	next.Provider = ProviderConfig{
		Host:         "provider.example",
		Port:         8443,
		HTTPS:        true,
		APIKey:       "runtime-secret",
		Model:        "response-model",
		KeywordModel: "keyword-model",
	}
	if err := engine.SaveSettings(context.Background(), &next); err != nil {
		t.Fatalf("SaveSettings() error = %v", err)
	}
	if engine.cfg.Theme != defaults.Theme {
		t.Fatalf("engine.cfg.Theme = %q, want %q", engine.cfg.Theme, defaults.Theme)
	}
	if engine.cfg.Web.Port != 9988 {
		t.Fatalf("engine.cfg.Web.Port = %d, want 9988", engine.cfg.Web.Port)
	}
	if engine.cfg.Provider.Host != "provider.example" || engine.cfg.Provider.KeywordModel != "keyword-model" {
		t.Fatalf("engine.cfg.Provider = %+v, want saved provider values", engine.cfg.Provider)
	}

	reloaded, err := engine.LoadSettings(context.Background())
	if err != nil {
		t.Fatalf("LoadSettings() reloaded error = %v", err)
	}
	if reloaded.Theme != defaults.Theme {
		t.Fatalf("reloaded.Theme = %q, want default %q", reloaded.Theme, defaults.Theme)
	}
	if reloaded.Web.Port != 9988 {
		t.Fatalf("reloaded.Web.Port = %d, want 9988", reloaded.Web.Port)
	}
	if reloaded.Provider.Host != "provider.example" || reloaded.Provider.Model != "response-model" {
		t.Fatalf("reloaded.Provider = %+v, want persisted provider settings", reloaded.Provider)
	}

	if err := engine.store.SetSetting("settings", "{invalid-json"); err != nil {
		t.Fatalf("SetSetting(invalid settings) error = %v", err)
	}
	fallback, err := engine.LoadSettings(context.Background())
	if err != nil {
		t.Fatalf("LoadSettings() invalid-json error = %v", err)
	}
	if !reflect.DeepEqual(fallback, defaults) {
		t.Fatalf("LoadSettings() invalid-json = %+v, want defaults %+v", fallback, defaults)
	}
}

func TestQuoteCRUDAndRecallHistoryLifecycle(t *testing.T) {
	t.Parallel()

	engine := newProfiledTestEngine(t, "localhost")

	count, err := engine.CountQuotes(context.Background())
	if err != nil {
		t.Fatalf("CountQuotes() initial error = %v", err)
	}
	if count != 0 {
		t.Fatalf("CountQuotes() initial = %d, want 0", count)
	}

	first, err := engine.AddQuote(context.Background(), "first quote")
	if err != nil {
		t.Fatalf("AddQuote(first) error = %v", err)
	}
	second, err := engine.AddQuote(context.Background(), "second quote")
	if err != nil {
		t.Fatalf("AddQuote(second) error = %v", err)
	}

	updated, err := engine.UpdateQuote(context.Background(), first.ID, "updated first quote")
	if err != nil {
		t.Fatalf("UpdateQuote() error = %v", err)
	}
	if updated.Content != "updated first quote" {
		t.Fatalf("updated.Content = %q, want updated first quote", updated.Content)
	}

	count, err = engine.CountQuotes(context.Background())
	if err != nil {
		t.Fatalf("CountQuotes() after add error = %v", err)
	}
	if count != 2 {
		t.Fatalf("CountQuotes() after add = %d, want 2", count)
	}

	if _, err := engine.SaveRecallHistory(context.Background(), "", "response", []Quote{*updated}); err == nil || !strings.Contains(err.Error(), "history question is empty") {
		t.Fatalf("SaveRecallHistory(empty question) error = %v, want question validation", err)
	}
	if _, err := engine.SaveRecallHistory(context.Background(), "question", "", []Quote{*updated}); err == nil || !strings.Contains(err.Error(), "history response is empty") {
		t.Fatalf("SaveRecallHistory(empty response) error = %v, want response validation", err)
	}

	entry, err := engine.SaveRecallHistory(context.Background(), "How do I update quotes?", "Edit and save the quote.", []Quote{*updated, *second})
	if err != nil {
		t.Fatalf("SaveRecallHistory() error = %v", err)
	}
	if len(entry.Quotes) != 2 {
		t.Fatalf("history quote count = %d, want 2", len(entry.Quotes))
	}

	summaries, err := engine.ListRecallHistory(context.Background())
	if err != nil {
		t.Fatalf("ListRecallHistory() error = %v", err)
	}
	if len(summaries) != 1 || summaries[0].ID != entry.ID {
		t.Fatalf("history summaries = %+v, want one saved entry", summaries)
	}

	loadedEntry, err := engine.GetRecallHistory(context.Background(), entry.ID)
	if err != nil {
		t.Fatalf("GetRecallHistory() error = %v", err)
	}
	if loadedEntry.Question != "How do I update quotes?" || len(loadedEntry.Quotes) != 2 {
		t.Fatalf("loaded history = %+v, want saved question and quotes", loadedEntry)
	}

	if err := engine.DeleteRecallHistory(context.Background(), []int64{entry.ID}); err != nil {
		t.Fatalf("DeleteRecallHistory() error = %v", err)
	}
	summaries, err = engine.ListRecallHistory(context.Background())
	if err != nil {
		t.Fatalf("ListRecallHistory() after delete error = %v", err)
	}
	if len(summaries) != 0 {
		t.Fatalf("history summaries after delete = %+v, want empty", summaries)
	}

	if err := engine.DeleteQuote(context.Background(), first.ID); err != nil {
		t.Fatalf("DeleteQuote() error = %v", err)
	}
	if err := engine.DeleteQuotes(context.Background(), []int64{second.ID}); err != nil {
		t.Fatalf("DeleteQuotes() error = %v", err)
	}
	count, err = engine.CountQuotes(context.Background())
	if err != nil {
		t.Fatalf("CountQuotes() after delete error = %v", err)
	}
	if count != 0 {
		t.Fatalf("CountQuotes() after delete = %d, want 0", count)
	}
}

func TestLoadUserProfileAndBootstrapQuoteIdentity(t *testing.T) {
	t.Parallel()

	engine := newTestEngine(t, "localhost")
	if err := engine.BootstrapQuoteIdentity(context.Background()); err == nil || !strings.Contains(err.Error(), "user profile not loaded") {
		t.Fatalf("BootstrapQuoteIdentity() error = %v, want profile-not-loaded failure", err)
	}

	profile, err := engine.LoadUserProfile(context.Background())
	if err != nil {
		t.Fatalf("LoadUserProfile() error = %v", err)
	}
	if profile.UserID == "" {
		t.Fatal("LoadUserProfile() returned empty user ID")
	}
	if engine.profile == nil || engine.profile.UserID != profile.UserID {
		t.Fatalf("engine.profile = %+v, want loaded profile %+v", engine.profile, profile)
	}

	legacyID, err := engine.store.InsertQuote("legacy quote", db.QuoteIdentity{})
	if err != nil {
		t.Fatalf("InsertQuote(legacy) error = %v", err)
	}
	if err := engine.BootstrapQuoteIdentity(context.Background()); err != nil {
		t.Fatalf("BootstrapQuoteIdentity() error = %v", err)
	}

	quote, err := engine.loadQuote(legacyID)
	if err != nil {
		t.Fatalf("loadQuote() error = %v", err)
	}
	if quote.GlobalID == "" || quote.AuthorUserID != profile.UserID || quote.SourceUserID != profile.UserID {
		t.Fatalf("bootstrapped quote = %+v, want local user identity fields", quote)
	}
	if quote.SourceBackend != "local" || quote.SourceNamespace != "local:"+profile.UserID || quote.SourceEntityType != "quote" {
		t.Fatalf("bootstrapped quote provenance = %+v, want local provenance", quote)
	}
	if quote.SourceEntityID != quote.GlobalID || quote.SourceLabel != "Local quote" || quote.Version != 1 {
		t.Fatalf("bootstrapped quote identity = %+v, want normalized local identity", quote)
	}
}

func TestSharedQuoteHelpersAndSchemaV1Import(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		entry SharedQuoteEntry
		want  string
	}{
		{
			name: "missing global id",
			entry: SharedQuoteEntry{
				AuthorUserID: "author-1",
				SourceUserID: "source-1",
				Version:      1,
				Content:      "hello",
			},
			want: "missing global_id",
		},
		{
			name: "missing author user id",
			entry: SharedQuoteEntry{
				GlobalID:     "quote-1",
				SourceUserID: "source-1",
				Version:      1,
				Content:      "hello",
			},
			want: "missing author_user_id",
		},
		{
			name: "invalid version",
			entry: SharedQuoteEntry{
				GlobalID:     "quote-1",
				AuthorUserID: "author-1",
				SourceUserID: "source-1",
				Version:      0,
				Content:      "hello",
			},
			want: "invalid version",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateSharedQuoteEntry(tt.entry)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("validateSharedQuoteEntry() error = %v, want %q", err, tt.want)
			}
		})
	}

	normalized := normalizeImportedQuoteIdentity(ShareSchemaVersion, db.QuoteIdentity{
		GlobalID:         "quote-1",
		SourceUserID:     "source-1",
		SourceName:       "Source Name",
		SourceBackend:    "custom",
		SourceNamespace:  "custom:source-1",
		SourceEntityType: "ticket",
		SourceEntityID:   "42",
		SourceLabel:      "Ticket #42",
	})
	if normalized.SourceBackend != "custom" || normalized.SourceNamespace != "custom:source-1" || normalized.SourceLabel != "Ticket #42" {
		t.Fatalf("normalizeImportedQuoteIdentity() unexpectedly rewrote populated fields: %+v", normalized)
	}

	engine := newProfiledTestEngine(t, "")
	payload, err := json.Marshal(SharedQuoteEnvelope{
		SchemaVersion: 1,
		ExportedAt:    time.Now().UTC(),
		Quotes: []SharedQuoteEntry{
			{
				GlobalID:     "legacy-shared-1",
				AuthorUserID: "author-1",
				AuthorName:   "Alice",
				SourceUserID: "source-1",
				Version:      1,
				Content:      "legacy import",
				Tags:         []string{"legacy", "shared"},
				CreatedAtUTC: time.Unix(100, 0).UTC(),
				UpdatedAtUTC: time.Unix(200, 0).UTC(),
			},
		},
	})
	if err != nil {
		t.Fatalf("Marshal(schema v1 payload) error = %v", err)
	}

	result, err := engine.ImportSharedQuotes(context.Background(), payload)
	if err != nil {
		t.Fatalf("ImportSharedQuotes(schema v1) error = %v", err)
	}
	if result.Inserted != 1 {
		t.Fatalf("ImportSharedQuotes(schema v1) result = %+v, want inserted=1", result)
	}

	quotes, err := engine.ListQuotes(context.Background())
	if err != nil {
		t.Fatalf("ListQuotes() error = %v", err)
	}
	if len(quotes) != 1 {
		t.Fatalf("quote count = %d, want 1", len(quotes))
	}
	if quotes[0].SourceBackend != "shared_import" || quotes[0].SourceNamespace != "share:source-1" || quotes[0].SourceEntityType != "shared_quote" {
		t.Fatalf("schema v1 import provenance = %+v, want shared import backfill", quotes[0])
	}
	if quotes[0].SourceEntityID != "legacy-shared-1" || quotes[0].SourceLabel != "Shared import" {
		t.Fatalf("schema v1 import label/id = %+v, want legacy backfill", quotes[0])
	}
	if !reflect.DeepEqual(quotes[0].Tags, []string{"legacy", "shared"}) {
		t.Fatalf("schema v1 import tags = %#v, want %#v", quotes[0].Tags, []string{"legacy", "shared"})
	}
}
