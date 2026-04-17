package core

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const webPasswordHashSettingKey = "web.password_hash"

const minWebPasswordLength = 12

var commonWebPasswords = []string{
	"12345678",
	"123456789",
	"1234567890",
	"password",
	"password1",
	"qwerty123",
	"letmein123",
	"admin1234",
	"changeme123",
}

func (e *Engine) HasWebPassword(ctx context.Context) (bool, error) {
	_ = ctx
	value, err := e.store.GetSetting(webPasswordHashSettingKey)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(value) != "", nil
}

func (e *Engine) SetupWebPassword(ctx context.Context, password, confirm string) error {
	_ = ctx
	if err := validatePasswordChange("", password, confirm, false); err != nil {
		return err
	}
	hasPassword, err := e.HasWebPassword(ctx)
	if err != nil {
		return err
	}
	if hasPassword {
		return errors.New("web password is already configured")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash web password: %w", err)
	}
	return e.store.SetSetting(webPasswordHashSettingKey, string(hash))
}

func (e *Engine) VerifyWebPassword(ctx context.Context, password string) (bool, error) {
	_ = ctx
	hash, err := e.store.GetSetting(webPasswordHashSettingKey)
	if err != nil {
		return false, err
	}
	hash = strings.TrimSpace(hash)
	if hash == "" {
		return false, errors.New("web password is not configured")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, fmt.Errorf("verify web password: %w", err)
	}
	return true, nil
}

func (e *Engine) ChangeWebPassword(ctx context.Context, current, next, confirm string) error {
	if err := validatePasswordChange(current, next, confirm, true); err != nil {
		return err
	}
	ok, err := e.VerifyWebPassword(ctx, current)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("current password is incorrect")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(next), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash web password: %w", err)
	}
	return e.store.SetSetting(webPasswordHashSettingKey, string(hash))
}

func validatePasswordChange(current, next, confirm string, requireCurrent bool) error {
	if requireCurrent && strings.TrimSpace(current) == "" {
		return errors.New("current password is required")
	}
	if strings.TrimSpace(next) == "" {
		return errors.New("password is required")
	}
	if next != confirm {
		return errors.New("passwords do not match")
	}
	if len(next) < minWebPasswordLength {
		return fmt.Errorf("password must be at least %d characters", minWebPasswordLength)
	}
	if isWeakWebPassword(next) {
		return errors.New("password is too weak; use a longer password with a mix of upper, lower, digit, or symbol characters")
	}
	return nil
}

func isWeakWebPassword(password string) bool {
	normalized := strings.ToLower(strings.TrimSpace(password))
	if slices.Contains(commonWebPasswords, normalized) {
		return true
	}
	if len(normalized) > 0 && normalized == strings.Repeat(string(normalized[0]), len(normalized)) {
		return true
	}

	var hasLower, hasUpper, hasDigit, hasSymbol bool
	for _, r := range password {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		default:
			hasSymbol = true
		}
	}

	var classes int
	for _, present := range []bool{hasLower, hasUpper, hasDigit, hasSymbol} {
		if present {
			classes++
		}
	}
	return classes < 3
}
