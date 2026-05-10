package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"testing"

	irecallapp "github.com/gigol/irecall/app"
	"github.com/gigol/irecall/config"
	"github.com/gigol/irecall/core"
	frontendassets "github.com/gigol/irecall/frontend"
)

func TestBearerTokenAuthenticatesAppRoutes(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	tokenResult, err := app.CreateAPIToken()
	if err != nil {
		t.Fatalf("CreateAPIToken() error = %v", err)
	}
	server := newTestServer(t, app)

	req := httptest.NewRequest(http.MethodGet, "/api/app/list-quotes", nil)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/app/list-quotes without auth = %d, want %d", res.Code, http.StatusUnauthorized)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/app/list-quotes", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResult.Token)
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("GET /api/app/list-quotes with bearer token = %d, want %d", res.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/app/count-quotes", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResult.Token)
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("GET /api/app/count-quotes with bearer token = %d, want %d", res.Code, http.StatusOK)
	}
	var count struct {
		Count int64 `json:"count"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &count); err != nil {
		t.Fatalf("decode count response: %v", err)
	}
	if count.Count != 0 {
		t.Fatalf("count response = %d, want 0", count.Count)
	}
}

func TestRecallHistoryPaginationRoutesViaBearerAuth(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	if _, err := app.SaveUserProfile("Alice"); err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}
	if _, err := app.AddQuote("first second question history pagination reference"); err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}
	settings := *app.GetSettings()
	settings.Debug.MockLLM = true
	if _, err := app.SaveSettings(settings); err != nil {
		t.Fatalf("SaveSettings(MockLLM) error = %v", err)
	}
	if _, err := app.RunRecall("first question"); err != nil {
		t.Fatalf("RunRecall(first) error = %v", err)
	}
	if _, err := app.RunRecall("second question"); err != nil {
		t.Fatalf("RunRecall(second) error = %v", err)
	}
	tokenResult, err := app.CreateAPIToken()
	if err != nil {
		t.Fatalf("CreateAPIToken() error = %v", err)
	}
	server := newTestServer(t, app)

	countReq := httptest.NewRequest(http.MethodGet, "/api/app/count-recall-history", nil)
	countReq.Header.Set("Authorization", "Bearer "+tokenResult.Token)
	countRes := httptest.NewRecorder()
	server.ServeHTTP(countRes, countReq)
	if countRes.Code != http.StatusOK {
		t.Fatalf("GET /api/app/count-recall-history = %d, body = %s", countRes.Code, countRes.Body.String())
	}
	var count struct {
		Count int64 `json:"count"`
	}
	if err := json.Unmarshal(countRes.Body.Bytes(), &count); err != nil {
		t.Fatalf("decode history count response: %v", err)
	}
	if count.Count != 2 {
		t.Fatalf("history count = %d, want 2", count.Count)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/app/list-recall-history?limit=1&offset=1", nil)
	listReq.Header.Set("Authorization", "Bearer "+tokenResult.Token)
	listRes := httptest.NewRecorder()
	server.ServeHTTP(listRes, listReq)
	if listRes.Code != http.StatusOK {
		t.Fatalf("GET /api/app/list-recall-history = %d, body = %s", listRes.Code, listRes.Body.String())
	}
	var entries []core.RecallHistorySummary
	if err := json.Unmarshal(listRes.Body.Bytes(), &entries); err != nil {
		t.Fatalf("decode paged history response: %v", err)
	}
	if len(entries) != 1 || entries[0].Question != "first question" {
		t.Fatalf("paged history entries = %+v, want older entry only", entries)
	}
}

func TestAPIOnlyModeStartsWithoutWebPasswordAndRequiresBearerToken(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	tokenResult, err := app.CreateAPIToken()
	if err != nil {
		t.Fatalf("CreateAPIToken() error = %v", err)
	}
	server := newTestServerWithOptions(t, app, ServerOptions{APIOnly: true})

	statusReq := httptest.NewRequest(http.MethodGet, "/api/auth/status", nil)
	statusRes := httptest.NewRecorder()
	server.ServeHTTP(statusRes, statusReq)
	if statusRes.Code != http.StatusOK {
		t.Fatalf("GET /api/auth/status = %d, want %d", statusRes.Code, http.StatusOK)
	}
	var status struct {
		PasswordConfigured bool `json:"passwordConfigured"`
		Authenticated      bool `json:"authenticated"`
	}
	if err := json.Unmarshal(statusRes.Body.Bytes(), &status); err != nil {
		t.Fatalf("decode auth status: %v", err)
	}
	if status.PasswordConfigured {
		t.Fatalf("passwordConfigured = true, want false")
	}
	if status.Authenticated {
		t.Fatalf("authenticated = true, want false")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/app/list-quotes", nil)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/app/list-quotes without auth = %d, want %d", res.Code, http.StatusUnauthorized)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/app/list-quotes", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResult.Token)
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("GET /api/app/list-quotes with bearer token = %d, want %d", res.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/app/count-quotes", nil)
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/app/count-quotes without auth = %d, want %d", res.Code, http.StatusUnauthorized)
	}
}

func TestAPIOnlyModeDisablesBrowserSessionRoutes(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	server := newTestServerWithOptions(t, app, ServerOptions{APIOnly: true})

	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", jsonBody(t, map[string]string{
		"password": "Secret-pass-123!",
	}))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	server.ServeHTTP(loginRes, loginReq)
	if loginRes.Code != http.StatusForbidden {
		t.Fatalf("POST /api/auth/login in api-only mode = %d, want %d", loginRes.Code, http.StatusForbidden)
	}

	tokenReq := httptest.NewRequest(http.MethodPost, "/api/app/create-api-token", nil)
	tokenRes := httptest.NewRecorder()
	server.ServeHTTP(tokenRes, tokenReq)
	if tokenRes.Code != http.StatusForbidden {
		t.Fatalf("POST /api/app/create-api-token in api-only mode = %d, want %d", tokenRes.Code, http.StatusForbidden)
	}
}

func TestCreateAPITokenRequiresSession(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	if err := app.SetupPassword("Secret-pass-123!", "Secret-pass-123!"); err != nil {
		t.Fatalf("SetupPassword() error = %v", err)
	}
	server := newTestServer(t, app)

	req := httptest.NewRequest(http.MethodPost, "/api/app/create-api-token", nil)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("POST /api/app/create-api-token without session = %d, want %d", res.Code, http.StatusUnauthorized)
	}

	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", jsonBody(t, map[string]string{
		"password": "Secret-pass-123!",
	}))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	server.ServeHTTP(loginRes, loginReq)
	if loginRes.Code != http.StatusOK {
		t.Fatalf("POST /api/auth/login = %d, want %d", loginRes.Code, http.StatusOK)
	}
	cookies := loginRes.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatalf("login returned no cookies")
	}

	statusReq := httptest.NewRequest(http.MethodGet, "/api/app/get-api-token-status", nil)
	statusReq.AddCookie(cookies[0])
	statusRes := httptest.NewRecorder()
	server.ServeHTTP(statusRes, statusReq)
	if statusRes.Code != http.StatusOK {
		t.Fatalf("GET /api/app/get-api-token-status before create = %d, want %d", statusRes.Code, http.StatusOK)
	}
	var status struct {
		HasToken bool `json:"hasToken"`
	}
	if err := json.Unmarshal(statusRes.Body.Bytes(), &status); err != nil {
		t.Fatalf("decode status before create: %v", err)
	}
	if status.HasToken {
		t.Fatalf("status before create HasToken = true, want false")
	}

	createReq := httptest.NewRequest(http.MethodPost, "/api/app/create-api-token", nil)
	createReq.AddCookie(cookies[0])
	createRes := httptest.NewRecorder()
	server.ServeHTTP(createRes, createReq)
	if createRes.Code != http.StatusOK {
		t.Fatalf("POST /api/app/create-api-token with session = %d, want %d", createRes.Code, http.StatusOK)
	}
	var created struct {
		Token       string `json:"token"`
		TokenPrefix string `json:"tokenPrefix"`
	}
	if err := json.Unmarshal(createRes.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created token: %v", err)
	}
	if created.Token == "" || created.TokenPrefix == "" {
		t.Fatalf("created token response = %#v, want token and prefix", created)
	}

	bearerReq := httptest.NewRequest(http.MethodGet, "/api/app/list-quotes", nil)
	bearerReq.Header.Set("Authorization", "Bearer "+created.Token)
	bearerRes := httptest.NewRecorder()
	server.ServeHTTP(bearerRes, bearerReq)
	if bearerRes.Code != http.StatusOK {
		t.Fatalf("GET /api/app/list-quotes with created bearer token = %d, want %d", bearerRes.Code, http.StatusOK)
	}

	renewWithBearerReq := httptest.NewRequest(http.MethodPost, "/api/app/create-api-token", nil)
	renewWithBearerReq.Header.Set("Authorization", "Bearer "+created.Token)
	renewWithBearerRes := httptest.NewRecorder()
	server.ServeHTTP(renewWithBearerRes, renewWithBearerReq)
	if renewWithBearerRes.Code != http.StatusUnauthorized {
		t.Fatalf("POST /api/app/create-api-token with bearer only = %d, want %d", renewWithBearerRes.Code, http.StatusUnauthorized)
	}
}

func TestHandleSaveSettingsPreservesExistingRootWhenOmitted(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "xdg-config"))

	root := filepath.Join(t.TempDir(), "web-root")
	runtimeApp, err := irecallapp.NewApp(root)
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	t.Cleanup(func() { runtimeApp.Shutdown(context.Background()) })

	current := runtimeApp.GetSettings()
	if current == nil {
		t.Fatal("GetSettings() returned nil")
	}
	if current.RootDir == "" {
		t.Fatal("GetSettings().RootDir = empty, want persisted root")
	}

	reqBody, err := json.Marshal(struct {
		Provider struct {
			Host         string `json:"Host"`
			Port         int    `json:"Port"`
			HTTPS        bool   `json:"HTTPS"`
			APIKey       string `json:"APIKey"`
			Model        string `json:"Model"`
			KeywordModel string `json:"KeywordModel"`
		} `json:"Provider"`
		Search struct {
			MaxResults   int     `json:"MaxResults"`
			MinRelevance float64 `json:"MinRelevance"`
		} `json:"Search"`
		Debug struct {
			MockLLM bool `json:"MockLLM"`
		} `json:"Debug"`
		Theme string `json:"Theme"`
		Web   struct {
			Port int `json:"Port"`
		} `json:"Web"`
	}{
		Provider: struct {
			Host         string `json:"Host"`
			Port         int    `json:"Port"`
			HTTPS        bool   `json:"HTTPS"`
			APIKey       string `json:"APIKey"`
			Model        string `json:"Model"`
			KeywordModel string `json:"KeywordModel"`
		}{
			Host:         current.Provider.Host,
			Port:         current.Provider.Port,
			HTTPS:        current.Provider.HTTPS,
			APIKey:       current.Provider.APIKey,
			Model:        current.Provider.Model,
			KeywordModel: current.Provider.KeywordModel,
		},
		Search: struct {
			MaxResults   int     `json:"MaxResults"`
			MinRelevance float64 `json:"MinRelevance"`
		}{
			MaxResults:   current.Search.MaxResults,
			MinRelevance: current.Search.MinRelevance,
		},
		Debug: struct {
			MockLLM bool `json:"MockLLM"`
		}{
			MockLLM: current.Debug.MockLLM,
		},
		Theme: "forest",
		Web: struct {
			Port int `json:"Port"`
		}{
			Port: current.Web.Port,
		},
	})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	server, err := NewServer(runtimeApp, frontendassets.Assets, current.Web.Port, ServerOptions{})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/app/save-settings", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.handleSaveSettings(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("handleSaveSettings() status = %d, body = %s", rec.Code, rec.Body.String())
	}

	saved := runtimeApp.GetSettings()
	if saved == nil {
		t.Fatal("GetSettings() after save returned nil")
	}
	if saved.Theme != "forest" {
		t.Fatalf("saved theme = %q, want %q", saved.Theme, "forest")
	}
	if saved.RootDir != current.RootDir {
		t.Fatalf("saved root = %q, want %q", saved.RootDir, current.RootDir)
	}

	preferredRoot, err := config.LoadPreferredRootPath()
	if err != nil {
		t.Fatalf("LoadPreferredRootPath() error = %v", err)
	}
	if preferredRoot != current.RootDir {
		t.Fatalf("preferred root = %q, want %q", preferredRoot, current.RootDir)
	}
}

func TestHandleRegenerateQuoteKeywordsSupportsGlobalID(t *testing.T) {
	t.Parallel()

	callCount := 0
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path = %q, want /v1/chat/completions", r.URL.Path)
		}
		callCount++
		tags := []string{"legacy"}
		if callCount > 1 {
			tags = []string{"sqlite", "wal", "concurrency"}
		}
		content, err := json.Marshal(tags)
		if err != nil {
			t.Fatalf("json.Marshal(tags) error = %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{
					"message": map[string]string{
						"content": string(content),
					},
				},
			},
		})
	}))
	defer provider.Close()

	app := newTestApp(t)
	if _, err := app.SaveUserProfile("Test User"); err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}
	if err := app.ApplyRuntimeProvider(core.ProviderConfig{
		Host:  provider.Listener.Addr().String(),
		Port:  0,
		HTTPS: false,
		Model: "test-model",
	}); err != nil {
		t.Fatalf("ApplyRuntimeProvider() error = %v", err)
	}

	quote, err := app.AddQuote("SQLite WAL helps readers and writers overlap safely.")
	if err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}

	tokenResult, err := app.CreateAPIToken()
	if err != nil {
		t.Fatalf("CreateAPIToken() error = %v", err)
	}
	server := newTestServer(t, app)
	wantNewKeywords := []string{"sqlite", "wal", "concurrency"}

	req := httptest.NewRequest(http.MethodPost, "/api/app/regenerate-quote-keywords", jsonBody(t, map[string]any{
		"globalId": quote.GlobalID,
	}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenResult.Token)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("POST /api/app/regenerate-quote-keywords = %d, body = %s", res.Code, res.Body.String())
	}

	var payload struct {
		QuoteID     int64      `json:"quoteId"`
		GlobalID    string     `json:"globalId"`
		OldKeywords []string   `json:"oldKeywords"`
		NewKeywords []string   `json:"newKeywords"`
		Changed     bool       `json:"changed"`
		Status      string     `json:"status"`
		Quote       core.Quote `json:"quote"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.QuoteID != quote.ID {
		t.Fatalf("quoteId = %d, want %d", payload.QuoteID, quote.ID)
	}
	if payload.GlobalID != quote.GlobalID {
		t.Fatalf("globalId = %q, want %q", payload.GlobalID, quote.GlobalID)
	}
	if !slices.Equal(payload.OldKeywords, []string{"legacy"}) {
		t.Fatalf("oldKeywords = %#v, want %#v", payload.OldKeywords, []string{"legacy"})
	}
	if !slices.Equal(payload.NewKeywords, wantNewKeywords) {
		t.Fatalf("newKeywords = %#v, want %#v", payload.NewKeywords, wantNewKeywords)
	}
	if !payload.Changed {
		t.Fatal("changed = false, want true")
	}
	if payload.Status != "updated" {
		t.Fatalf("status = %q, want updated", payload.Status)
	}
	gotTags := slices.Clone(payload.Quote.Tags)
	slices.Sort(gotTags)
	wantTags := slices.Clone(wantNewKeywords)
	slices.Sort(wantTags)
	if !slices.Equal(gotTags, wantTags) {
		t.Fatalf("quote tags = %#v, want %#v", payload.Quote.Tags, wantNewKeywords)
	}
}

func TestQuoteRoutesUpdateAndDeleteQuotes(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	first, err := app.AddQuote("first quote")
	if err != nil {
		t.Fatalf("AddQuote(first) error = %v", err)
	}
	second, err := app.AddQuote("second quote")
	if err != nil {
		t.Fatalf("AddQuote(second) error = %v", err)
	}
	tokenResult, err := app.CreateAPIToken()
	if err != nil {
		t.Fatalf("CreateAPIToken() error = %v", err)
	}
	server := newTestServer(t, app)

	updateReq := httptest.NewRequest(http.MethodPost, "/api/app/update-quote", jsonBody(t, map[string]any{
		"id":      first.ID,
		"content": "updated first quote",
	}))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("Authorization", "Bearer "+tokenResult.Token)
	updateRes := httptest.NewRecorder()
	server.ServeHTTP(updateRes, updateReq)
	if updateRes.Code != http.StatusOK {
		t.Fatalf("POST /api/app/update-quote = %d, body = %s", updateRes.Code, updateRes.Body.String())
	}

	var updated core.Quote
	if err := json.Unmarshal(updateRes.Body.Bytes(), &updated); err != nil {
		t.Fatalf("decode update response: %v", err)
	}
	if updated.ID != first.ID || updated.Content != "updated first quote" {
		t.Fatalf("updated quote = %+v, want updated first quote", updated)
	}

	deleteReq := httptest.NewRequest(http.MethodPost, "/api/app/delete-quotes", jsonBody(t, map[string]any{
		"ids": []int64{second.ID},
	}))
	deleteReq.Header.Set("Content-Type", "application/json")
	deleteReq.Header.Set("Authorization", "Bearer "+tokenResult.Token)
	deleteRes := httptest.NewRecorder()
	server.ServeHTTP(deleteRes, deleteReq)
	if deleteRes.Code != http.StatusOK {
		t.Fatalf("POST /api/app/delete-quotes = %d, body = %s", deleteRes.Code, deleteRes.Body.String())
	}

	var deleted struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal(deleteRes.Body.Bytes(), &deleted); err != nil {
		t.Fatalf("decode delete response: %v", err)
	}
	if !deleted.OK {
		t.Fatal("delete response ok = false, want true")
	}

	quotes, err := app.ListQuotes()
	if err != nil {
		t.Fatalf("ListQuotes() error = %v", err)
	}
	if len(quotes) != 1 || quotes[0].ID != first.ID || quotes[0].Content != "updated first quote" {
		t.Fatalf("quotes after update/delete = %+v, want only updated first quote", quotes)
	}
}

func TestHandleSaveUserProfileViaBearerAuth(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	tokenResult, err := app.CreateAPIToken()
	if err != nil {
		t.Fatalf("CreateAPIToken() error = %v", err)
	}
	server := newTestServer(t, app)

	req := httptest.NewRequest(http.MethodPost, "/api/app/save-user-profile", jsonBody(t, map[string]string{
		"name": "Alice",
	}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenResult.Token)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("POST /api/app/save-user-profile = %d, body = %s", res.Code, res.Body.String())
	}

	var profile core.UserProfile
	if err := json.Unmarshal(res.Body.Bytes(), &profile); err != nil {
		t.Fatalf("decode profile response: %v", err)
	}
	if profile.DisplayName != "Alice" || profile.UserID == "" {
		t.Fatalf("profile = %+v, want saved display name and user ID", profile)
	}
	if got := app.GetUserProfile(); got == nil || got.DisplayName != "Alice" {
		t.Fatalf("app.GetUserProfile() = %+v, want updated profile", got)
	}
}

func TestHandleChangePasswordViaSession(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	if err := app.SetupPassword("Secret-pass-123!", "Secret-pass-123!"); err != nil {
		t.Fatalf("SetupPassword() error = %v", err)
	}
	server := newTestServer(t, app)

	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", jsonBody(t, map[string]string{
		"password": "Secret-pass-123!",
	}))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	server.ServeHTTP(loginRes, loginReq)
	if loginRes.Code != http.StatusOK {
		t.Fatalf("POST /api/auth/login = %d, want %d", loginRes.Code, http.StatusOK)
	}
	cookies := loginRes.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("login returned no cookies")
	}

	changeReq := httptest.NewRequest(http.MethodPost, "/api/auth/change-password", jsonBody(t, map[string]string{
		"current": "Secret-pass-123!",
		"next":    "EvenBetter-456!",
		"confirm": "EvenBetter-456!",
	}))
	changeReq.Header.Set("Content-Type", "application/json")
	changeReq.AddCookie(cookies[0])
	changeRes := httptest.NewRecorder()
	server.ServeHTTP(changeRes, changeReq)
	if changeRes.Code != http.StatusOK {
		t.Fatalf("POST /api/auth/change-password = %d, body = %s", changeRes.Code, changeRes.Body.String())
	}

	if err := app.Login("Secret-pass-123!"); err == nil {
		t.Fatal("Login(old password) error = nil, want failure")
	}
	if err := app.Login("EvenBetter-456!"); err != nil {
		t.Fatalf("Login(new password) error = %v", err)
	}
}

func TestHandleRunRecallViaBearerAuth(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	if _, err := app.SaveUserProfile("Alice"); err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}
	if _, err := app.AddQuote("alpha beta note"); err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}
	settings := *app.GetSettings()
	settings.Debug.MockLLM = true
	if _, err := app.SaveSettings(settings); err != nil {
		t.Fatalf("SaveSettings(MockLLM) error = %v", err)
	}
	tokenResult, err := app.CreateAPIToken()
	if err != nil {
		t.Fatalf("CreateAPIToken() error = %v", err)
	}
	server := newTestServer(t, app)

	req := httptest.NewRequest(http.MethodPost, "/api/app/run-recall", jsonBody(t, map[string]string{
		"question": "alpha beta",
	}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenResult.Token)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("POST /api/app/run-recall = %d, body = %s", res.Code, res.Body.String())
	}

	var result struct {
		Question string       `json:"question"`
		Keywords []string     `json:"keywords"`
		Quotes   []core.Quote `json:"quotes"`
		Response string       `json:"response"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode recall response: %v", err)
	}
	if result.Question != "alpha beta" {
		t.Fatalf("question = %q, want alpha beta", result.Question)
	}
	if !slices.Equal(result.Keywords, []string{"alpha", "beta"}) {
		t.Fatalf("keywords = %#v, want %#v", result.Keywords, []string{"alpha", "beta"})
	}
	if len(result.Quotes) != 1 || result.Quotes[0].Content != "alpha beta note" {
		t.Fatalf("quotes = %+v, want matching note", result.Quotes)
	}
	if result.Response != "alpha beta note" {
		t.Fatalf("response = %q, want mock joined quote content", result.Response)
	}
}

func newTestApp(t *testing.T) *irecallapp.App {
	t.Helper()

	app, err := irecallapp.NewAppWithOptions(t.TempDir(), irecallapp.AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions() error = %v", err)
	}
	t.Cleanup(func() {
		app.Shutdown(context.Background())
	})
	return app
}

func newTestServer(t *testing.T, app *irecallapp.App) http.Handler {
	return newTestServerWithOptions(t, app, ServerOptions{})
}

func newTestServerWithOptions(t *testing.T, app *irecallapp.App, options ServerOptions) http.Handler {
	t.Helper()

	server, err := NewServer(app, frontendassets.Assets, 9527, options)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}
	return server.Handler()
}

func TestUnsafeNoPasswordCheckBypassesAuth(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)

	server, err := NewServer(app, frontendassets.Assets, 9527, ServerOptions{UnsafeNoPasswordCheck: true})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}
	handler := server.Handler()

	statusReq := httptest.NewRequest(http.MethodGet, "/api/auth/status", nil)
	statusRes := httptest.NewRecorder()
	handler.ServeHTTP(statusRes, statusReq)
	if statusRes.Code != http.StatusOK {
		t.Fatalf("GET /api/auth/status = %d, want %d", statusRes.Code, http.StatusOK)
	}

	var status struct {
		PasswordConfigured bool `json:"passwordConfigured"`
		Authenticated      bool `json:"authenticated"`
	}
	if err := json.Unmarshal(statusRes.Body.Bytes(), &status); err != nil {
		t.Fatalf("decode auth status: %v", err)
	}
	if !status.PasswordConfigured || !status.Authenticated {
		t.Fatalf("auth status = %#v, want passwordConfigured and authenticated true", status)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/app/list-quotes", nil)
	listRes := httptest.NewRecorder()
	handler.ServeHTTP(listRes, listReq)
	if listRes.Code != http.StatusOK {
		t.Fatalf("GET /api/app/list-quotes = %d, want %d", listRes.Code, http.StatusOK)
	}
}

func jsonBody(t *testing.T, value any) *bytes.Reader {
	t.Helper()

	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	return bytes.NewReader(payload)
}
