package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/gigol/irecall/config"
	"github.com/gigol/irecall/core"
)

func TestDesktopBackendQuoteShareRoundTrip(t *testing.T) {
	t.Parallel()

	alice, err := NewAppWithOptions(filepath.Join(t.TempDir(), "alice"), AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions(alice) error = %v", err)
	}
	t.Cleanup(func() { alice.Shutdown(context.Background()) })

	if _, err := alice.SaveUserProfile("Alice"); err != nil {
		t.Fatalf("SaveUserProfile(alice) error = %v", err)
	}
	q, err := alice.AddQuote("desktop share roundtrip")
	if err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}

	exportPath := filepath.Join(t.TempDir(), "quotes", "share.json")
	if err := alice.ExportQuotesToFile([]int64{q.ID}, exportPath); err != nil {
		t.Fatalf("ExportQuotesToFile() error = %v", err)
	}
	if _, err := os.Stat(exportPath); err != nil {
		t.Fatalf("Stat(exportPath) error = %v", err)
	}

	bob, err := NewAppWithOptions(filepath.Join(t.TempDir(), "bob"), AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions(bob) error = %v", err)
	}
	t.Cleanup(func() { bob.Shutdown(context.Background()) })

	if _, err := bob.SaveUserProfile("Bob"); err != nil {
		t.Fatalf("SaveUserProfile(bob) error = %v", err)
	}
	result, err := bob.ImportQuotesFromFile(exportPath)
	if err != nil {
		t.Fatalf("ImportQuotesFromFile() error = %v", err)
	}
	if result.Inserted != 1 {
		t.Fatalf("import result = %+v, want inserted=1", result)
	}

	quotes, err := bob.ListQuotes()
	if err != nil {
		t.Fatalf("ListQuotes() error = %v", err)
	}
	if len(quotes) != 1 {
		t.Fatalf("quote count = %d, want 1", len(quotes))
	}
	if quotes[0].Content != "desktop share roundtrip" {
		t.Fatalf("imported content = %q, want desktop share roundtrip", quotes[0].Content)
	}
	if quotes[0].SourceName != "Alice" {
		t.Fatalf("source name = %q, want Alice", quotes[0].SourceName)
	}
}

func TestDesktopBackendBootstrapState(t *testing.T) {
	t.Parallel()

	app, err := NewAppWithOptions(filepath.Join(t.TempDir(), "desktop"), AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	profile, err := app.SaveUserProfile("Alice")
	if err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}

	state := app.BootstrapState()
	if state.ProductName != "iRecall" {
		t.Fatalf("product name = %q, want iRecall", state.ProductName)
	}
	if state.Greeting != "Hi! Alice" {
		t.Fatalf("greeting = %q, want Hi! Alice", state.Greeting)
	}
	if state.Profile == nil || state.Profile.DisplayName != profile.DisplayName {
		t.Fatalf("profile = %+v, want display name %q", state.Profile, profile.DisplayName)
	}
	if len(state.Pages) != 4 {
		t.Fatalf("page count = %d, want 4", len(state.Pages))
	}
	if state.Pages[1] != "History" {
		t.Fatalf("pages = %v, want History tab in bootstrap state", state.Pages)
	}
	if state.Docs["uiDesign"] != "docs/UI_DESIGN.md" {
		t.Fatalf("ui design doc = %q, want docs/UI_DESIGN.md", state.Docs["uiDesign"])
	}
}

func TestDesktopBackendCountQuotes(t *testing.T) {
	t.Parallel()

	app, err := NewAppWithOptions(filepath.Join(t.TempDir(), "desktop-count"), AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	count, err := app.CountQuotes()
	if err != nil {
		t.Fatalf("CountQuotes() before insert error = %v", err)
	}
	if count != 0 {
		t.Fatalf("CountQuotes() before insert = %d, want 0", count)
	}

	if _, err := app.AddQuote("count me"); err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}

	count, err = app.CountQuotes()
	if err != nil {
		t.Fatalf("CountQuotes() after insert error = %v", err)
	}
	if count != 1 {
		t.Fatalf("CountQuotes() after insert = %d, want 1", count)
	}
}

func TestDesktopBackendRecallHistoryLifecycle(t *testing.T) {
	t.Parallel()

	app, err := NewAppWithOptions(filepath.Join(t.TempDir(), "desktop-history"), AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	if _, err := app.SaveUserProfile("Alice"); err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}

	quote, err := app.AddQuote("History-enabled desktop quote")
	if err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}

	if _, err := app.engine.SaveRecallHistory(context.Background(),
		"How do I check history?",
		"Open the History tab and inspect the saved session.",
		[]core.Quote{*quote},
	); err != nil {
		t.Fatalf("SaveRecallHistory() error = %v", err)
	}
	if _, err := app.engine.SaveRecallHistory(context.Background(),
		"How do I move between pages?",
		"Use previous and next controls.",
		[]core.Quote{*quote},
	); err != nil {
		t.Fatalf("SaveRecallHistory(second) error = %v", err)
	}

	count, err := app.CountRecallHistory()
	if err != nil {
		t.Fatalf("CountRecallHistory() error = %v", err)
	}
	if count != 2 {
		t.Fatalf("CountRecallHistory() = %d, want 2", count)
	}

	historyPage, err := app.ListRecallHistoryPage(1, 0)
	if err != nil {
		t.Fatalf("ListRecallHistoryPage() error = %v", err)
	}
	if len(historyPage) != 1 || historyPage[0].Question != "How do I move between pages?" {
		t.Fatalf("history page = %+v, want newest history entry", historyPage)
	}

	history, err := app.ListRecallHistory()
	if err != nil {
		t.Fatalf("ListRecallHistory() error = %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("history count = %d, want 2", len(history))
	}

	entry, err := app.GetRecallHistory(history[1].ID)
	if err != nil {
		t.Fatalf("GetRecallHistory() error = %v", err)
	}
	if entry.Question != "How do I check history?" {
		t.Fatalf("history question = %q, want exact saved question", entry.Question)
	}
	if len(entry.Quotes) != 1 || entry.Quotes[0].ID != quote.ID {
		t.Fatalf("history quotes = %+v, want original quote", entry.Quotes)
	}

	if err := app.DeleteRecallHistory([]int64{history[0].ID, history[1].ID}); err != nil {
		t.Fatalf("DeleteRecallHistory() error = %v", err)
	}

	history, err = app.ListRecallHistory()
	if err != nil {
		t.Fatalf("ListRecallHistory() after delete error = %v", err)
	}
	if len(history) != 0 {
		t.Fatalf("history count after delete = %d, want 0", len(history))
	}
}

func TestDesktopBackendSaveRecallAsQuote(t *testing.T) {
	t.Parallel()

	app, err := NewAppWithOptions(filepath.Join(t.TempDir(), "desktop-recall-quote"), AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	if _, err := app.SaveUserProfile("Alice"); err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}

	quote, err := app.SaveRecallAsQuote(
		"How do I export quotes?",
		"Use the share flow and select the quotes you want to export.",
		[]string{"export", "sharing"},
	)
	if err != nil {
		t.Fatalf("SaveRecallAsQuote() error = %v", err)
	}
	if quote.ID == 0 {
		t.Fatalf("quote id = %d, want persisted quote", quote.ID)
	}
	if quote.Content == "" || quote.Content[:9] != "Question:" {
		t.Fatalf("quote content = %q, want formatted recall quote", quote.Content)
	}
	if len(quote.Tags) == 0 {
		t.Fatalf("quote tags = %#v, want saved tags", quote.Tags)
	}
}

func TestDesktopBackendImportQuotesPayload(t *testing.T) {
	t.Parallel()

	app, err := NewAppWithOptions(filepath.Join(t.TempDir(), "desktop-import-payload"), AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	if _, err := app.SaveUserProfile("Alice"); err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}

	quote, err := app.AddQuote("payload import roundtrip")
	if err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}

	payload, err := app.PreviewQuoteExport([]int64{quote.ID})
	if err != nil {
		t.Fatalf("PreviewQuoteExport() error = %v", err)
	}

	target, err := NewAppWithOptions(filepath.Join(t.TempDir(), "desktop-import-target"), AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions(target) error = %v", err)
	}
	t.Cleanup(func() { target.Shutdown(context.Background()) })

	if _, err := target.SaveUserProfile("Bob"); err != nil {
		t.Fatalf("SaveUserProfile(target) error = %v", err)
	}

	result, err := target.ImportQuotesPayload(payload)
	if err != nil {
		t.Fatalf("ImportQuotesPayload() error = %v", err)
	}
	if result.Inserted != 1 {
		t.Fatalf("ImportQuotesPayload() result = %+v, want inserted=1", result)
	}
}

func TestNewAppUsesDefaultStorageWhenRootIsEmpty(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", filepath.Join(t.TempDir(), "xdg-data"))
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "xdg-config"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(t.TempDir(), "xdg-state"))

	app, err := NewAppWithOptions("", AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions(\"\") error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	if app.paths.RootDir != "" {
		t.Fatalf("app.paths.RootDir = %q, want empty for default storage", app.paths.RootDir)
	}
	for _, path := range []string{app.paths.DataDir, app.paths.ConfigDir, app.paths.StateDir} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("Stat(%q) error = %v", path, err)
		}
		if !info.IsDir() {
			t.Fatalf("%q is not a directory", path)
		}
	}
}

func TestDesktopBackendSaveSettingsSwitchesStorageRoot(t *testing.T) {
	xdgConfig := filepath.Join(t.TempDir(), "xdg-config")
	t.Setenv("XDG_CONFIG_HOME", xdgConfig)

	sourceRoot := filepath.Join(t.TempDir(), "desktop-source")
	targetRoot := filepath.Join(t.TempDir(), "desktop-target")
	app, err := NewApp(sourceRoot)
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	if _, err := app.SaveUserProfile("Alice"); err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}
	if _, err := app.AddQuote("switch storage root"); err != nil {
		t.Fatalf("AddQuote() error = %v", err)
	}

	settings := *app.GetSettings()
	settings.RootDir = targetRoot
	settings.Theme = "forest"

	saved, err := app.SaveSettings(settings)
	if err != nil {
		t.Fatalf("SaveSettings() error = %v", err)
	}

	absTarget, err := filepath.Abs(targetRoot)
	if err != nil {
		t.Fatalf("Abs(targetRoot) error = %v", err)
	}
	if saved.RootDir != absTarget {
		t.Fatalf("saved root = %q, want %q", saved.RootDir, absTarget)
	}
	if app.paths.RootDir != absTarget {
		t.Fatalf("app.paths.RootDir = %q, want %q", app.paths.RootDir, absTarget)
	}
	if state := app.BootstrapState(); state.Paths.RootDir != absTarget {
		t.Fatalf("bootstrap root = %q, want %q", state.Paths.RootDir, absTarget)
	}

	quotes, err := app.ListQuotes()
	if err != nil {
		t.Fatalf("ListQuotes() error = %v", err)
	}
	if len(quotes) != 1 || quotes[0].Content != "switch storage root" {
		t.Fatalf("quotes after switch = %+v, want migrated quote", quotes)
	}

	if _, err := os.Stat(filepath.Join(absTarget, "data", "irecall.db")); err != nil {
		t.Fatalf("Stat(target db) error = %v", err)
	}
	preferredRoot, err := config.LoadPreferredRootPath()
	if err != nil {
		t.Fatalf("LoadPreferredRootPath() error = %v", err)
	}
	if preferredRoot != absTarget {
		t.Fatalf("preferred root = %q, want %q", preferredRoot, absTarget)
	}
	if _, err := os.Stat(filepath.Join(xdgConfig, "irecall", config.PreferredRootFileName)); err != nil {
		t.Fatalf("preferred root marker missing from isolated config dir: %v", err)
	}
}

func TestApplyRuntimeProviderInitializesMissingSettings(t *testing.T) {
	t.Parallel()

	root := filepath.Join(t.TempDir(), "runtime-provider")
	app, err := NewAppWithOptions(root, AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	app.settings = nil
	provider := core.ProviderConfig{
		Host:         "provider.example/api",
		Port:         443,
		HTTPS:        true,
		APIKey:       "runtime-secret",
		Model:        "runtime-model",
		KeywordModel: "runtime-keyword-model",
	}

	if err := app.ApplyRuntimeProvider(provider); err != nil {
		t.Fatalf("ApplyRuntimeProvider() error = %v", err)
	}
	if app.settings == nil {
		t.Fatal("settings = nil, want initialized settings")
	}
	if app.settings.RootDir != app.paths.RootDir {
		t.Fatalf("settings.RootDir = %q, want %q", app.settings.RootDir, app.paths.RootDir)
	}
	if app.settings.Provider != provider {
		t.Fatalf("settings.Provider = %+v, want %+v", app.settings.Provider, provider)
	}
}

func TestDesktopBackendUpdateAndDeleteQuotes(t *testing.T) {
	t.Parallel()

	app, err := NewAppWithOptions(filepath.Join(t.TempDir(), "desktop-update-delete"), AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	first, err := app.AddQuote("first quote")
	if err != nil {
		t.Fatalf("AddQuote(first) error = %v", err)
	}
	second, err := app.AddQuote("second quote")
	if err != nil {
		t.Fatalf("AddQuote(second) error = %v", err)
	}

	updated, err := app.UpdateQuote(first.ID, "updated first quote")
	if err != nil {
		t.Fatalf("UpdateQuote() error = %v", err)
	}
	if updated.Content != "updated first quote" {
		t.Fatalf("updated content = %q, want updated first quote", updated.Content)
	}

	if err := app.DeleteQuotes([]int64{second.ID}); err != nil {
		t.Fatalf("DeleteQuotes() error = %v", err)
	}

	quotes, err := app.ListQuotes()
	if err != nil {
		t.Fatalf("ListQuotes() error = %v", err)
	}
	if len(quotes) != 1 {
		t.Fatalf("quote count = %d, want 1", len(quotes))
	}
	if quotes[0].ID != first.ID || quotes[0].Content != "updated first quote" {
		t.Fatalf("remaining quote = %+v, want updated first quote only", quotes[0])
	}
}

func TestDesktopBackendPasswordAndAPITokenHelpers(t *testing.T) {
	t.Parallel()

	app, err := NewAppWithOptions(filepath.Join(t.TempDir(), "desktop-auth-helpers"), AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	originalPassword := "Secret-pass-123!"
	nextPassword := "EvenBetter-456!"

	if err := app.SetupPassword(originalPassword, originalPassword); err != nil {
		t.Fatalf("SetupPassword() error = %v", err)
	}
	if err := app.ChangePassword(originalPassword, nextPassword, nextPassword); err != nil {
		t.Fatalf("ChangePassword() error = %v", err)
	}
	if err := app.Login(originalPassword); err == nil || err.Error() != "invalid password" {
		t.Fatalf("Login(old password) error = %v, want invalid password", err)
	}
	if err := app.Login(nextPassword); err != nil {
		t.Fatalf("Login(new password) error = %v", err)
	}

	tokenResult, err := app.CreateAPITokenWithPassword(nextPassword)
	if err != nil {
		t.Fatalf("CreateAPITokenWithPassword() error = %v", err)
	}
	if tokenResult.Token == "" || tokenResult.TokenPrefix == "" {
		t.Fatalf("token result = %+v, want token and prefix", tokenResult)
	}

	ok, err := app.VerifyAPIToken(tokenResult.Token)
	if err != nil {
		t.Fatalf("VerifyAPIToken(created) error = %v", err)
	}
	if !ok {
		t.Fatal("VerifyAPIToken(created) = false, want true")
	}

	if err := app.RevokeAPITokenWithPassword(nextPassword); err != nil {
		t.Fatalf("RevokeAPITokenWithPassword() error = %v", err)
	}
	ok, err = app.VerifyAPIToken(tokenResult.Token)
	if err != nil {
		t.Fatalf("VerifyAPIToken(revoked) error = %v", err)
	}
	if ok {
		t.Fatal("VerifyAPIToken(revoked) = true, want false")
	}

	if err := app.ResetPassword(); err != nil {
		t.Fatalf("ResetPassword() error = %v", err)
	}
	if err := app.Login(nextPassword); err == nil || err.Error() != "web password is not configured" {
		t.Fatalf("Login(after reset) error = %v, want not configured", err)
	}
}

func TestDesktopBackendGetUserProfileReflectsSavedProfile(t *testing.T) {
	t.Parallel()

	app, err := NewAppWithOptions(filepath.Join(t.TempDir(), "desktop-profile"), AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	initial := app.GetUserProfile()
	if initial == nil {
		t.Fatal("GetUserProfile() = nil, want initialized profile")
	}
	if initial.DisplayName != "" {
		t.Fatalf("initial display name = %q, want empty", initial.DisplayName)
	}

	saved, err := app.SaveUserProfile("Alice")
	if err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}
	got := app.GetUserProfile()
	if got == nil || got.DisplayName != "Alice" || got.UserID != saved.UserID {
		t.Fatalf("GetUserProfile() = %+v, want saved profile %+v", got, saved)
	}
}

func TestDesktopBackendRunRecallWithMockLLM(t *testing.T) {
	t.Parallel()

	app, err := NewAppWithOptions(filepath.Join(t.TempDir(), "desktop-run-recall"), AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

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

	result, err := app.RunRecall("alpha beta")
	if err != nil {
		t.Fatalf("RunRecall() error = %v", err)
	}
	if result.Question != "alpha beta" {
		t.Fatalf("question = %q, want alpha beta", result.Question)
	}
	if len(result.Keywords) != 2 || result.Keywords[0] != "alpha" || result.Keywords[1] != "beta" {
		t.Fatalf("keywords = %#v, want alpha/beta", result.Keywords)
	}
	if len(result.Quotes) != 1 || result.Quotes[0].Content != "alpha beta note" {
		t.Fatalf("quotes = %+v, want matching note", result.Quotes)
	}
	if result.Response != "alpha beta note" {
		t.Fatalf("response = %q, want joined mock quote content", result.Response)
	}
}

func TestDesktopBackendCanSkipPreferredRootPersistence(t *testing.T) {
	xdgConfig := filepath.Join(t.TempDir(), "xdg-config")
	t.Setenv("XDG_CONFIG_HOME", xdgConfig)

	sourceRoot := filepath.Join(t.TempDir(), "source")
	targetRoot := filepath.Join(t.TempDir(), "target")
	app, err := NewAppWithOptions(sourceRoot, AppOptions{})
	if err != nil {
		t.Fatalf("NewAppWithOptions() error = %v", err)
	}
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	settings := *app.GetSettings()
	settings.RootDir = targetRoot
	if _, err := app.SaveSettings(settings); err != nil {
		t.Fatalf("SaveSettings() error = %v", err)
	}

	preferredRoot, err := config.LoadPreferredRootPath()
	if err != nil {
		t.Fatalf("LoadPreferredRootPath() error = %v", err)
	}
	if preferredRoot != "" {
		t.Fatalf("preferred root = %q, want empty when app persistence disabled", preferredRoot)
	}
	if _, err := os.Stat(filepath.Join(xdgConfig, "irecall", config.PreferredRootFileName)); !os.IsNotExist(err) {
		t.Fatalf("preferred root marker exists or stat failed: %v", err)
	}
}
