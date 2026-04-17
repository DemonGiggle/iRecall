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

	if err := engine.SetupWebPassword(ctx, "secret-pass", "secret-pass"); err != nil {
		t.Fatalf("SetupWebPassword() error = %v", err)
	}

	ok, err := engine.VerifyWebPassword(ctx, "secret-pass")
	if err != nil {
		t.Fatalf("VerifyWebPassword() error = %v", err)
	}
	if !ok {
		t.Fatalf("VerifyWebPassword() = false, want true")
	}

	if err := engine.ChangeWebPassword(ctx, "secret-pass", "changed-pass", "changed-pass"); err != nil {
		t.Fatalf("ChangeWebPassword() error = %v", err)
	}

	ok, err = engine.VerifyWebPassword(ctx, "changed-pass")
	if err != nil {
		t.Fatalf("VerifyWebPassword(changed) error = %v", err)
	}
	if !ok {
		t.Fatalf("VerifyWebPassword(changed) = false, want true")
	}

	ok, err = engine.VerifyWebPassword(ctx, "secret-pass")
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
	if err := engine.SetupWebPassword(context.Background(), "one", "two"); err == nil {
		t.Fatalf("SetupWebPassword() error = nil, want mismatch error")
	}
}
