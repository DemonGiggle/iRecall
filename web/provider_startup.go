package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gigol/irecall/app"
	"github.com/gigol/irecall/core"
)

type optionalBoolFlag struct {
	set   bool
	value bool
}

func (f *optionalBoolFlag) String() string {
	if !f.set {
		return ""
	}
	return strconv.FormatBool(f.value)
}

func (f *optionalBoolFlag) Set(value string) error {
	parsed, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return err
	}
	f.value = parsed
	f.set = true
	return nil
}

func (f *optionalBoolFlag) IsBoolFlag() bool {
	return true
}

type ProviderStartupOptions struct {
	Host       string
	Port       int
	HTTPS      optionalBoolFlag
	APIKeyPath string
	Model      string
}

func bindProviderStartupFlags(fs *flag.FlagSet) *ProviderStartupOptions {
	opts := &ProviderStartupOptions{}
	fs.StringVar(&opts.Host, "provider-host", "", "override provider host/server for --api-only mode")
	fs.IntVar(&opts.Port, "provider-port", 0, "override provider port for --api-only mode")
	fs.Var(&opts.HTTPS, "provider-https", "override provider HTTPS setting for --api-only mode")
	fs.StringVar(&opts.APIKeyPath, "provider-api-key-path", "", "read provider API key from this file for --api-only mode")
	fs.StringVar(&opts.Model, "provider-model", "", "override provider model for --api-only mode")
	return opts
}

func (o ProviderStartupOptions) HasOverrides() bool {
	return strings.TrimSpace(o.Host) != "" ||
		o.Port != 0 ||
		o.HTTPS.set ||
		strings.TrimSpace(o.APIKeyPath) != "" ||
		strings.TrimSpace(o.Model) != ""
}

func (o ProviderStartupOptions) Validate(serverOptions ServerOptions) error {
	if !serverOptions.APIOnly && o.HasOverrides() {
		return errors.New("provider startup flags require --api-only")
	}
	if o.Port != 0 && (o.Port < 1 || o.Port > 65535) {
		return fmt.Errorf("provider port must be a number between 1 and 65535")
	}
	return nil
}

func (o ProviderStartupOptions) Resolve(base core.ProviderConfig) (core.ProviderConfig, error) {
	resolved := base
	if host := strings.TrimSpace(o.Host); host != "" {
		resolved.Host = host
	}
	if o.Port != 0 {
		resolved.Port = o.Port
	}
	if o.HTTPS.set {
		resolved.HTTPS = o.HTTPS.value
	}
	if model := strings.TrimSpace(o.Model); model != "" {
		resolved.Model = model
	}
	if apiKeyPath := strings.TrimSpace(o.APIKeyPath); apiKeyPath != "" {
		apiKey, err := readSecretFile(apiKeyPath, "provider API key")
		if err != nil {
			return core.ProviderConfig{}, err
		}
		resolved.APIKey = apiKey
	}
	return resolved, nil
}

func applyAPIOnlyProviderStartupConfig(runtimeApp *app.App, serverOptions ServerOptions, opts ProviderStartupOptions) error {
	if err := opts.Validate(serverOptions); err != nil {
		return err
	}
	if !serverOptions.APIOnly || !opts.HasOverrides() {
		return nil
	}
	if runtimeApp == nil {
		return errors.New("app is not initialized")
	}

	base := core.DefaultSettings().Provider
	if settings := runtimeApp.GetSettings(); settings != nil {
		base = settings.Provider
	}
	provider, err := opts.Resolve(base)
	if err != nil {
		return err
	}
	return runtimeApp.ApplyRuntimeProvider(provider)
}

func readSecretFile(path string, label string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s file %q: %w", label, path, err)
	}
	value := strings.TrimSpace(string(data))
	if value == "" {
		return "", fmt.Errorf("%s file %q is empty", label, path)
	}
	return value, nil
}
