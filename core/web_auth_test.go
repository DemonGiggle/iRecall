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
