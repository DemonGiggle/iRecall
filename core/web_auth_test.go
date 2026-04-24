package core

import (
	"context"
	"testing"

	"github.com/gigol/irecall/core/db"
)

func TestWebPasswordLifecycle(t *testing.T) {
	t.Parallel()

	store, err := db.Open(t.TempDir() + "/irecall.db")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	engine := New(store, DefaultSettings())
	ctx := context.Background()

	hasPassword, err := engine.HasWebPassword(ctx)
	if err != nil {
		t.Fatalf("HasWebPassword() error = %v", err)
	}
	if hasPassword {
		t.Fatalf("HasWebPassword() = true, want false")
	}

	if err := engine.SetupWebPassword(ctx, "Secret-pass-123!", "Secret-pass-123!"); err != nil {
		t.Fatalf("SetupWebPassword() error = %v", err)
	}

	ok, err := engine.VerifyWebPassword(ctx, "Secret-pass-123!")
	if err != nil {
		t.Fatalf("VerifyWebPassword() error = %v", err)
	}
	if !ok {
		t.Fatalf("VerifyWebPassword() = false, want true")
	}

	if err := engine.ChangeWebPassword(ctx, "Secret-pass-123!", "Changed-pass-456!", "Changed-pass-456!"); err != nil {
		t.Fatalf("ChangeWebPassword() error = %v", err)
	}

	ok, err = engine.VerifyWebPassword(ctx, "Changed-pass-456!")
	if err != nil {
		t.Fatalf("VerifyWebPassword(changed) error = %v", err)
	}
	if !ok {
		t.Fatalf("VerifyWebPassword(changed) = false, want true")
	}

	ok, err = engine.VerifyWebPassword(ctx, "Secret-pass-123!")
	if err != nil {
		t.Fatalf("VerifyWebPassword(old) error = %v", err)
	}
	if ok {
		t.Fatalf("VerifyWebPassword(old) = true, want false")
	}
}

func TestResetWebPasswordClearsConfiguredPassword(t *testing.T) {
	t.Parallel()

	store, err := db.Open(t.TempDir() + "/irecall.db")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	engine := New(store, DefaultSettings())
	ctx := context.Background()

	if err := engine.SetupWebPassword(ctx, "Secret-pass-123!", "Secret-pass-123!"); err != nil {
		t.Fatalf("SetupWebPassword() error = %v", err)
	}

	if err := engine.ResetWebPassword(ctx); err != nil {
		t.Fatalf("ResetWebPassword() error = %v", err)
	}

	hasPassword, err := engine.HasWebPassword(ctx)
	if err != nil {
		t.Fatalf("HasWebPassword() error = %v", err)
	}
	if hasPassword {
		t.Fatalf("HasWebPassword() after reset = true, want false")
	}

	ok, err := engine.VerifyWebPassword(ctx, "Secret-pass-123!")
	if err == nil {
		t.Fatalf("VerifyWebPassword() error = nil, want not configured error")
	}
	if ok {
		t.Fatalf("VerifyWebPassword() after reset = true, want false")
	}
}

func TestSetupWebPasswordRejectsMismatch(t *testing.T) {
	t.Parallel()

	store, err := db.Open(t.TempDir() + "/irecall.db")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	engine := New(store, DefaultSettings())
	if err := engine.SetupWebPassword(context.Background(), "Strong-pass-123!", "Wrong-pass-123!"); err == nil {
		t.Fatalf("SetupWebPassword() error = nil, want mismatch error")
	}
}

func TestSetupWebPasswordRejectsWeakPassword(t *testing.T) {
	t.Parallel()

	store, err := db.Open(t.TempDir() + "/irecall.db")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	engine := New(store, DefaultSettings())
	err = engine.SetupWebPassword(context.Background(), "password123", "password123")
	if err == nil {
		t.Fatalf("SetupWebPassword() error = nil, want weak password error")
	}
}

func TestWebAPITokenLifecycle(t *testing.T) {
	t.Parallel()

	store, err := db.Open(t.TempDir() + "/irecall.db")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	engine := New(store, DefaultSettings())
	ctx := context.Background()

	status, err := engine.GetWebAPITokenStatus(ctx)
	if err != nil {
		t.Fatalf("GetWebAPITokenStatus() error = %v", err)
	}
	if status.HasToken {
		t.Fatalf("GetWebAPITokenStatus().HasToken = true, want false")
	}

	token, status, err := engine.GenerateWebAPIToken(ctx)
	if err != nil {
		t.Fatalf("GenerateWebAPIToken() error = %v", err)
	}
	if token == "" {
		t.Fatalf("GenerateWebAPIToken() token = empty")
	}
	if !status.HasToken {
		t.Fatalf("GenerateWebAPIToken().HasToken = false, want true")
	}
	if status.TokenPrefix == "" {
		t.Fatalf("GenerateWebAPIToken().TokenPrefix = empty")
	}
	if status.TokenPrefix == token {
		t.Fatalf("GenerateWebAPIToken().TokenPrefix = full token, want short prefix")
	}

	storedHash, err := engine.store.GetSetting(webAPITokenHashSettingKey)
	if err != nil {
		t.Fatalf("GetSetting(hash) error = %v", err)
	}
	if storedHash == "" {
		t.Fatalf("stored hash = empty")
	}
	if storedHash == token {
		t.Fatalf("stored hash = plaintext token, want hashed token")
	}

	ok, err := engine.VerifyWebAPIToken(ctx, token)
	if err != nil {
		t.Fatalf("VerifyWebAPIToken() error = %v", err)
	}
	if !ok {
		t.Fatalf("VerifyWebAPIToken() = false, want true")
	}

	nextToken, _, err := engine.GenerateWebAPIToken(ctx)
	if err != nil {
		t.Fatalf("GenerateWebAPIToken(renew) error = %v", err)
	}
	if nextToken == token {
		t.Fatalf("renewed token = original token, want replacement")
	}

	ok, err = engine.VerifyWebAPIToken(ctx, token)
	if err != nil {
		t.Fatalf("VerifyWebAPIToken(old) error = %v", err)
	}
	if ok {
		t.Fatalf("VerifyWebAPIToken(old) = true, want false")
	}

	ok, err = engine.VerifyWebAPIToken(ctx, nextToken)
	if err != nil {
		t.Fatalf("VerifyWebAPIToken(new) error = %v", err)
	}
	if !ok {
		t.Fatalf("VerifyWebAPIToken(new) = false, want true")
	}
}
