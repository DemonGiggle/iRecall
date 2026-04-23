package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	irecallapp "github.com/gigol/irecall/app"
	"github.com/gigol/irecall/config"
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
			Host   string `json:"Host"`
			Port   int    `json:"Port"`
			HTTPS  bool   `json:"HTTPS"`
			APIKey string `json:"APIKey"`
			Model  string `json:"Model"`
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
			Host   string `json:"Host"`
			Port   int    `json:"Port"`
			HTTPS  bool   `json:"HTTPS"`
			APIKey string `json:"APIKey"`
			Model  string `json:"Model"`
		}{
			Host:   current.Provider.Host,
			Port:   current.Provider.Port,
			HTTPS:  current.Provider.HTTPS,
			APIKey: current.Provider.APIKey,
			Model:  current.Provider.Model,
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

	server, err := NewServer(runtimeApp, frontendassets.Assets, current.Web.Port)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/app/save-settings", bytes.NewReader(reqBody))
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

func newTestApp(t *testing.T) *irecallapp.App {
	t.Helper()

	app, err := irecallapp.NewApp(t.TempDir())
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	t.Cleanup(func() {
		app.Shutdown(context.Background())
	})
	return app
}

func newTestServer(t *testing.T, app *irecallapp.App) http.Handler {
	t.Helper()

	server, err := NewServer(app, frontendassets.Assets, 9527)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}
	return server.Handler()
}

func jsonBody(t *testing.T, value any) *bytes.Reader {
	t.Helper()

	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	return bytes.NewReader(payload)
}
