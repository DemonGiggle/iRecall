package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	irecallapp "github.com/gigol/irecall/app"
	"github.com/gigol/irecall/frontend"
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

	server, err := NewServer(app, frontend.Assets, 9527)
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
