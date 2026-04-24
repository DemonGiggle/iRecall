package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gigol/irecall/config"
	"github.com/gigol/irecall/core"
)

type App struct {
	ctx      context.Context
	engine   *core.Engine
	settings *core.Settings
	profile  *core.UserProfile
	paths    AppPaths
}

type AppPaths struct {
	RootDir   string `json:"rootDir"`
	DataDir   string `json:"dataDir"`
	ConfigDir string `json:"configDir"`
	StateDir  string `json:"stateDir"`
	DBPath    string `json:"dbPath"`
	LogPath   string `json:"logPath"`
}

type BootstrapState struct {
	ProductName string            `json:"productName"`
	Greeting    string            `json:"greeting"`
	Profile     *core.UserProfile `json:"profile"`
	Settings    *core.Settings    `json:"settings"`
	Paths       AppPaths          `json:"paths"`
	Pages       []string          `json:"pages"`
	Docs        map[string]string `json:"docs"`
}

type RecallResult struct {
	Question string       `json:"question"`
	Keywords []string     `json:"keywords"`
	Quotes   []core.Quote `json:"quotes"`
	Response string       `json:"response"`
}

type AuthStatus struct {
	Runtime            string `json:"runtime"`
	PasswordConfigured bool   `json:"passwordConfigured"`
	Authenticated      bool   `json:"authenticated"`
	CurrentPort        int    `json:"currentPort"`
}

type APITokenStatus struct {
	HasToken    bool   `json:"hasToken"`
	TokenPrefix string `json:"tokenPrefix"`
}

type APITokenCreateResult struct {
	Token       string `json:"token"`
	TokenPrefix string `json:"tokenPrefix"`
}

func NewApp(root string) (*App, error) {
	runtimeState, err := OpenRuntime(root)
	if err != nil {
		return nil, err
	}

	return &App{
		engine:   runtimeState.Engine,
		settings: runtimeState.Settings,
		profile:  runtimeState.Profile,
		paths:    runtimeState.Paths,
	}, nil
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Shutdown(context.Context) {
	if a.engine != nil {
		_ = a.engine.Close()
	}
}

func (a *App) BootstrapState() BootstrapState {
	greeting := ""
	if a.profile != nil && strings.TrimSpace(a.profile.DisplayName) != "" {
		greeting = "Hi! " + strings.TrimSpace(a.profile.DisplayName)
	}
	return BootstrapState{
		ProductName: "iRecall",
		Greeting:    greeting,
		Profile:     a.profile,
		Settings:    a.settings,
		Paths:       a.paths,
		Pages:       []string{"Recall", "History", "Quotes", "Settings"},
		Docs: map[string]string{
			"uiDesign":       "docs/UI_DESIGN.md",
			"desktopMapping": "docs/WAILS_DESKTOP.md",
		},
	}
}

func (a *App) ListQuotes() ([]core.Quote, error) {
	return a.engine.ListQuotes(a.context())
}

func (a *App) AddQuote(content string) (*core.Quote, error) {
	return a.engine.AddQuote(a.context(), content)
}

func (a *App) SaveRecallAsQuote(question, response string, keywords []string) (*core.Quote, error) {
	return a.engine.SaveRecallAsQuote(a.context(), question, response, keywords)
}

func (a *App) RefineQuoteDraft(content string) (string, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return "", errors.New("quote draft is empty")
	}
	return a.engine.RefineQuoteDraft(a.context(), content)
}

func (a *App) UpdateQuote(id int64, content string) (*core.Quote, error) {
	return a.engine.UpdateQuote(a.context(), id, content)
}

func (a *App) DeleteQuotes(ids []int64) error {
	return a.engine.DeleteQuotes(a.context(), ids)
}

func (a *App) ListRecallHistory() ([]core.RecallHistorySummary, error) {
	return a.engine.ListRecallHistory(a.context())
}

func (a *App) GetRecallHistory(id int64) (*core.RecallHistoryEntry, error) {
	return a.engine.GetRecallHistory(a.context(), id)
}

func (a *App) DeleteRecallHistory(ids []int64) error {
	return a.engine.DeleteRecallHistory(a.context(), ids)
}

func (a *App) ExportQuotesToFile(ids []int64, path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("export path is empty")
	}
	payload, err := a.engine.ExportQuotes(a.context(), ids)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create export directory: %w", err)
	}
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		return fmt.Errorf("write export payload: %w", err)
	}
	return nil
}

func (a *App) PreviewQuoteExport(ids []int64) (string, error) {
	payload, err := a.engine.ExportQuotes(a.context(), ids)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func (a *App) ImportQuotesFromFile(path string) (core.ImportResult, error) {
	if strings.TrimSpace(path) == "" {
		return core.ImportResult{}, errors.New("import path is empty")
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		return core.ImportResult{}, fmt.Errorf("read import payload: %w", err)
	}
	return a.engine.ImportSharedQuotes(a.context(), payload)
}

func (a *App) ImportQuotesPayload(payload string) (core.ImportResult, error) {
	payload = strings.TrimSpace(payload)
	if payload == "" {
		return core.ImportResult{}, errors.New("import payload is empty")
	}
	return a.engine.ImportSharedQuotes(a.context(), []byte(payload))
}

func (a *App) GetSettings() *core.Settings {
	return a.settings
}

func (a *App) SaveSettings(settings core.Settings) (*core.Settings, error) {
	nextRuntime, err := SwitchRuntime(&RuntimeState{
		Engine:   a.engine,
		Settings: a.settings,
		Profile:  a.profile,
		Paths:    a.paths,
	}, &settings)
	if err != nil {
		if nextRuntime != nil {
			a.engine = nextRuntime.Engine
			a.settings = nextRuntime.Settings
			a.profile = nextRuntime.Profile
			a.paths = nextRuntime.Paths
		}
		return nil, err
	}
	a.engine = nextRuntime.Engine
	a.settings = nextRuntime.Settings
	a.profile = nextRuntime.Profile
	a.paths = nextRuntime.Paths
	return a.settings, nil
}

func (a *App) AuthStatus() (AuthStatus, error) {
	hasPassword, err := a.engine.HasWebPassword(a.context())
	if err != nil {
		return AuthStatus{}, err
	}
	return AuthStatus{
		Runtime:            "desktop",
		PasswordConfigured: hasPassword,
		Authenticated:      true,
		CurrentPort:        0,
	}, nil
}

func (a *App) SetupPassword(password, confirm string) error {
	return a.engine.SetupWebPassword(a.context(), password, confirm)
}

func (a *App) ResetPassword() error {
	return a.engine.ResetWebPassword(a.context())
}

func (a *App) Login(password string) error {
	ok, err := a.engine.VerifyWebPassword(a.context(), password)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("invalid password")
	}
	return nil
}

func (a *App) ChangePassword(current, next, confirm string) error {
	return a.engine.ChangeWebPassword(a.context(), current, next, confirm)
}

func (a *App) GetAPITokenStatus() (APITokenStatus, error) {
	status, err := a.engine.GetWebAPITokenStatus(a.context())
	if err != nil {
		return APITokenStatus{}, err
	}
	return APITokenStatus{
		HasToken:    status.HasToken,
		TokenPrefix: status.TokenPrefix,
	}, nil
}

func (a *App) CreateAPIToken() (APITokenCreateResult, error) {
	token, status, err := a.engine.GenerateWebAPIToken(a.context())
	if err != nil {
		return APITokenCreateResult{}, err
	}
	return APITokenCreateResult{
		Token:       token,
		TokenPrefix: status.TokenPrefix,
	}, nil
}

func (a *App) VerifyAPIToken(token string) (bool, error) {
	return a.engine.VerifyWebAPIToken(a.context(), token)
}

func (a *App) FetchModels(settings core.ProviderConfig) ([]string, error) {
	return a.engine.FetchModels(a.context(), settings)
}

func (a *App) GetUserProfile() *core.UserProfile {
	return a.profile
}

func (a *App) SaveUserProfile(name string) (*core.UserProfile, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("display name is empty")
	}
	profile := core.UserProfile{}
	if a.profile != nil {
		profile = *a.profile
	}
	profile.DisplayName = name
	if err := a.engine.SaveUserProfile(a.context(), &profile); err != nil {
		return nil, err
	}
	a.engine.UpdateUserProfile(&profile)
	a.profile = &profile
	return a.profile, nil
}

func (a *App) RunRecall(question string) (*RecallResult, error) {
	question = strings.TrimSpace(question)
	if question == "" {
		return nil, errors.New("question is empty")
	}

	keywords, err := a.engine.ExtractKeywords(a.context(), question)
	if err != nil {
		return nil, err
	}
	quotes, err := a.engine.SearchQuotes(a.context(), keywords)
	if err != nil {
		return nil, err
	}

	tokenCh := make(chan string, 64)
	errCh := make(chan error, 1)
	go func() {
		errCh <- a.engine.GenerateResponse(a.context(), question, quotes, tokenCh)
	}()

	var sb strings.Builder
	for token := range tokenCh {
		sb.WriteString(token)
	}
	if err := <-errCh; err != nil {
		return nil, err
	}
	if _, err := a.engine.SaveRecallHistory(a.context(), question, sb.String(), quotes); err != nil {
		return nil, err
	}

	return &RecallResult{
		Question: question,
		Keywords: keywords,
		Quotes:   quotes,
		Response: sb.String(),
	}, nil
}

func (a *App) context() context.Context {
	if a.ctx != nil {
		return a.ctx
	}
	return context.Background()
}

func resolvePaths(root string) (AppPaths, error) {
	var err error
	root = strings.TrimSpace(root)
	if root == "" {
		root = strings.TrimSpace(config.RootPath())
	}
	if root == "" {
		return AppPaths{
			RootDir:   "",
			DataDir:   config.DataDir(),
			ConfigDir: config.ConfigDir(),
			StateDir:  config.StateDir(),
			DBPath:    config.DBPath(),
			LogPath:   config.LogPath(),
		}, nil
	}
	root, err = normalizeRootDir(root)
	if err != nil {
		return AppPaths{}, err
	}
	return AppPaths{
		RootDir:   root,
		DataDir:   filepath.Join(root, "data"),
		ConfigDir: filepath.Join(root, "config"),
		StateDir:  filepath.Join(root, "state"),
		DBPath:    filepath.Join(root, "data", "irecall.db"),
		LogPath:   filepath.Join(root, "state", "irecall.log"),
	}, nil
}

func ensurePaths(paths AppPaths) error {
	for _, dir := range []string{paths.RootDir, paths.DataDir, paths.ConfigDir, paths.StateDir} {
		if strings.TrimSpace(dir) == "" {
			continue
		}
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return fmt.Errorf("create desktop app directory %q: %w", dir, err)
		}
	}
	return nil
}
