package main

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gigol/irecall/app"
	"github.com/gigol/irecall/core"
)

const sessionCookieName = "irecall_session"

type Server struct {
	app         *app.App
	currentPort int
	assets      fs.FS

	mu       sync.Mutex
	sessions map[string]time.Time
}

func NewServer(app *app.App, assets embed.FS, currentPort int) (*Server, error) {
	sub, err := fs.Sub(assets, "dist")
	if err != nil {
		return nil, fmt.Errorf("open frontend assets: %w", err)
	}
	return &Server{
		app:         app,
		currentPort: currentPort,
		assets:      sub,
		sessions:    make(map[string]time.Time),
	}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/status", s.handleAuthStatus)
	mux.HandleFunc("/api/auth/login", s.handleAuthLogin)
	mux.HandleFunc("/api/auth/logout", s.handleAuthLogout)
	mux.Handle("/api/auth/change-password", s.requireAuth(http.HandlerFunc(s.handleChangePassword)))

	mux.Handle("/api/app/bootstrap-state", s.requireAuth(http.HandlerFunc(s.handleBootstrapState)))
	mux.Handle("/api/app/list-quotes", s.requireAuth(http.HandlerFunc(s.handleListQuotes)))
	mux.Handle("/api/app/add-quote", s.requireAuth(http.HandlerFunc(s.handleAddQuote)))
	mux.Handle("/api/app/save-recall-as-quote", s.requireAuth(http.HandlerFunc(s.handleSaveRecallAsQuote)))
	mux.Handle("/api/app/refine-quote-draft", s.requireAuth(http.HandlerFunc(s.handleRefineQuoteDraft)))
	mux.Handle("/api/app/update-quote", s.requireAuth(http.HandlerFunc(s.handleUpdateQuote)))
	mux.Handle("/api/app/delete-quotes", s.requireAuth(http.HandlerFunc(s.handleDeleteQuotes)))
	mux.Handle("/api/app/preview-quote-export", s.requireAuth(http.HandlerFunc(s.handlePreviewQuoteExport)))
	mux.Handle("/api/app/import-quotes-payload", s.requireAuth(http.HandlerFunc(s.handleImportQuotesPayload)))
	mux.Handle("/api/app/save-user-profile", s.requireAuth(http.HandlerFunc(s.handleSaveUserProfile)))
	mux.Handle("/api/app/save-settings", s.requireAuth(http.HandlerFunc(s.handleSaveSettings)))
	mux.Handle("/api/app/fetch-models", s.requireAuth(http.HandlerFunc(s.handleFetchModels)))
	mux.Handle("/api/app/run-recall", s.requireAuth(http.HandlerFunc(s.handleRunRecall)))
	mux.Handle("/api/app/list-recall-history", s.requireAuth(http.HandlerFunc(s.handleListRecallHistory)))
	mux.Handle("/api/app/get-recall-history", s.requireAuth(http.HandlerFunc(s.handleGetRecallHistory)))
	mux.Handle("/api/app/delete-recall-history", s.requireAuth(http.HandlerFunc(s.handleDeleteRecallHistory)))

	mux.HandleFunc("/bridge.js", s.handleBridge)
	mux.HandleFunc("/", s.handleFrontend)
	return mux
}

func (s *Server) handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}
	hasPassword, err := s.app.AuthStatus()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"runtime":            "web",
		"passwordConfigured": hasPassword.PasswordConfigured,
		"authenticated":      s.isAuthenticated(r),
		"currentPort":        s.currentPort,
	})
}

func (s *Server) handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.app.Login(req.Password); err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}
	if err := s.startSession(w); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleAuthLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}
	s.endSession(w, r)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}
	var req struct {
		Current string `json:"current"`
		Next    string `json:"next"`
		Confirm string `json:"confirm"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.app.ChangePassword(req.Current, req.Next, req.Confirm); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleBootstrapState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}
	writeJSON(w, http.StatusOK, s.app.BootstrapState())
}

func (s *Server) handleListQuotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}
	value, err := s.app.ListQuotes()
	writeAppJSON(w, value, err)
}

func (s *Server) handleAddQuote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content string `json:"content"`
	}
	if !requirePostJSON(w, r, &req) {
		return
	}
	value, err := s.app.AddQuote(req.Content)
	writeAppJSON(w, value, err)
}

func (s *Server) handleSaveRecallAsQuote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Question string   `json:"question"`
		Response string   `json:"response"`
		Keywords []string `json:"keywords"`
	}
	if !requirePostJSON(w, r, &req) {
		return
	}
	value, err := s.app.SaveRecallAsQuote(req.Question, req.Response, req.Keywords)
	writeAppJSON(w, value, err)
}

func (s *Server) handleRefineQuoteDraft(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content string `json:"content"`
	}
	if !requirePostJSON(w, r, &req) {
		return
	}
	value, err := s.app.RefineQuoteDraft(req.Content)
	writeAppJSON(w, value, err)
}

func (s *Server) handleUpdateQuote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID      int64  `json:"id"`
		Content string `json:"content"`
	}
	if !requirePostJSON(w, r, &req) {
		return
	}
	value, err := s.app.UpdateQuote(req.ID, req.Content)
	writeAppJSON(w, value, err)
}

func (s *Server) handleDeleteQuotes(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDs []int64 `json:"ids"`
	}
	if !requirePostJSON(w, r, &req) {
		return
	}
	if err := s.app.DeleteQuotes(req.IDs); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handlePreviewQuoteExport(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDs []int64 `json:"ids"`
	}
	if !requirePostJSON(w, r, &req) {
		return
	}
	value, err := s.app.PreviewQuoteExport(req.IDs)
	writeAppJSON(w, value, err)
}

func (s *Server) handleImportQuotesPayload(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Payload string `json:"payload"`
	}
	if !requirePostJSON(w, r, &req) {
		return
	}
	value, err := s.app.ImportQuotesPayload(req.Payload)
	writeAppJSON(w, value, err)
}

func (s *Server) handleSaveUserProfile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if !requirePostJSON(w, r, &req) {
		return
	}
	value, err := s.app.SaveUserProfile(req.Name)
	writeAppJSON(w, value, err)
}

func (s *Server) handleSaveSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var settings struct {
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
		Theme string `json:"Theme"`
		Web   struct {
			Port int `json:"Port"`
		} `json:"Web"`
	}
	if err := json.Unmarshal(body, &settings); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := s.app.SaveSettings(backendToCoreSettings(settings))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleFetchModels(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Host   string `json:"Host"`
		Port   int    `json:"Port"`
		HTTPS  bool   `json:"HTTPS"`
		APIKey string `json:"APIKey"`
		Model  string `json:"Model"`
	}
	if !requirePostJSON(w, r, &req) {
		return
	}
	value, err := s.app.FetchModels(backendToCoreProvider(req))
	writeAppJSON(w, value, err)
}

func (s *Server) handleRunRecall(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Question string `json:"question"`
	}
	if !requirePostJSON(w, r, &req) {
		return
	}
	value, err := s.app.RunRecall(req.Question)
	writeAppJSON(w, value, err)
}

func (s *Server) handleListRecallHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}
	value, err := s.app.ListRecallHistory()
	writeAppJSON(w, value, err)
}

func (s *Server) handleGetRecallHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, errors.New("history id is required"))
		return
	}
	value, err := s.app.GetRecallHistory(id)
	writeAppJSON(w, value, err)
}

func (s *Server) handleDeleteRecallHistory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDs []int64 `json:"ids"`
	}
	if !requirePostJSON(w, r, &req) {
		return
	}
	if err := s.app.DeleteRecallHistory(req.IDs); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleBridge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	_, _ = io.WriteString(w, webBridgeJS)
}

func (s *Server) handleFrontend(w http.ResponseWriter, r *http.Request) {
	clean := path.Clean(strings.TrimPrefix(r.URL.Path, "/"))
	if clean == "." || clean == "" {
		s.serveIndex(w)
		return
	}
	file, err := s.assets.Open(clean)
	if err != nil {
		s.serveIndex(w)
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || info.IsDir() {
		s.serveIndex(w)
		return
	}
	http.FileServer(http.FS(s.assets)).ServeHTTP(w, r)
}

func (s *Server) serveIndex(w http.ResponseWriter) {
	index, err := fs.ReadFile(s.assets, "index.html")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	html := strings.Replace(string(index), "<head>", "<head>\n    <script src=\"/bridge.js\"></script>", 1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = io.WriteString(w, html)
}

func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.isAuthenticated(r) {
			writeError(w, http.StatusUnauthorized, errors.New("authentication required"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	expiresAt, ok := s.sessions[cookie.Value]
	if !ok {
		return false
	}
	if time.Now().After(expiresAt) {
		delete(s.sessions, cookie.Value)
		return false
	}
	s.sessions[cookie.Value] = time.Now().Add(24 * time.Hour)
	return true
}

func (s *Server) startSession(w http.ResponseWriter) error {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("generate session token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)
	s.mu.Lock()
	s.sessions[token] = time.Now().Add(24 * time.Hour)
	s.mu.Unlock()
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int((24 * time.Hour).Seconds()),
	})
	return nil
}

func (s *Server) endSession(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		s.mu.Lock()
		delete(s.sessions, cookie.Value)
		s.mu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}

func requirePostJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return false
	}
	if err := decodeJSON(r, target); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return false
	}
	return true
}

func decodeJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		return err
	}
	return nil
}

func writeAppJSON[T any](w http.ResponseWriter, value T, err error) {
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, value)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]any{
		"error": err.Error(),
	})
}

func writeMethodNotAllowed(w http.ResponseWriter) {
	writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
}

func backendToCoreProvider(v struct {
	Host   string `json:"Host"`
	Port   int    `json:"Port"`
	HTTPS  bool   `json:"HTTPS"`
	APIKey string `json:"APIKey"`
	Model  string `json:"Model"`
}) core.ProviderConfig {
	return core.ProviderConfig{
		Host:   v.Host,
		Port:   v.Port,
		HTTPS:  v.HTTPS,
		APIKey: v.APIKey,
		Model:  v.Model,
	}
}

func backendToCoreSettings(v struct {
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
	Theme string `json:"Theme"`
	Web   struct {
		Port int `json:"Port"`
	} `json:"Web"`
}) core.Settings {
	return core.Settings{
		Provider: core.ProviderConfig{
			Host:   v.Provider.Host,
			Port:   v.Provider.Port,
			HTTPS:  v.Provider.HTTPS,
			APIKey: v.Provider.APIKey,
			Model:  v.Provider.Model,
		},
		Search: core.SearchConfig{
			MaxResults:   v.Search.MaxResults,
			MinRelevance: v.Search.MinRelevance,
		},
		Theme: v.Theme,
		Web: core.WebConfig{
			Port: v.Web.Port,
		},
	}
}
