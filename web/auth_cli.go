//go:build !wails

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	irecallapp "github.com/gigol/irecall/app"
	"github.com/gigol/irecall/config"
)

type authCLIOptions struct {
	DataPath       string
	PasswordStdin  bool
	WriteTokenFile string
}

func maybeHandleAuthCommand(args []string) bool {
	if len(args) == 0 || args[0] != "auth" {
		return false
	}
	if err := runAuthCommand(args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "irecall-web auth: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
	return true
}

func runAuthCommand(args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) == 0 {
		return errors.New("missing auth subcommand: use issue-token, rotate-token, revoke-token, or token-status")
	}
	subcommand := args[0]
	fs := flag.NewFlagSet("irecall-web auth "+subcommand, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	opts := authCLIOptions{}
	fs.StringVar(&opts.DataPath, "data-path", "", "store database, config, and logs under this root path")
	fs.BoolVar(&opts.PasswordStdin, "password-stdin", false, "read the web password from stdin")
	fs.StringVar(&opts.WriteTokenFile, "write-token-file", "", "write the issued token to a file with mode 0600")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	switch subcommand {
	case "issue-token", "rotate-token":
		password, err := readAuthCommandPassword(stdin, opts.PasswordStdin)
		if err != nil {
			return err
		}
		runtimeApp, err := openAuthCommandApp(opts.DataPath)
		if err != nil {
			return err
		}
		defer runtimeApp.Shutdown(nil)
		result, err := runtimeApp.CreateAPITokenWithPassword(password)
		if err != nil {
			return err
		}
		if opts.WriteTokenFile != "" {
			if err := writeTokenFile(opts.WriteTokenFile, result.Token); err != nil {
				return err
			}
			fmt.Fprintf(stdout, "token written: %s\n", opts.WriteTokenFile)
		} else {
			fmt.Fprintln(stdout, result.Token)
		}
		fmt.Fprintf(stdout, "token prefix: %s\n", result.TokenPrefix)
		return nil
	case "revoke-token":
		password, err := readAuthCommandPassword(stdin, opts.PasswordStdin)
		if err != nil {
			return err
		}
		runtimeApp, err := openAuthCommandApp(opts.DataPath)
		if err != nil {
			return err
		}
		defer runtimeApp.Shutdown(nil)
		if err := runtimeApp.RevokeAPITokenWithPassword(password); err != nil {
			return err
		}
		fmt.Fprintln(stdout, "token revoked")
		return nil
	case "token-status":
		runtimeApp, err := openAuthCommandApp(opts.DataPath)
		if err != nil {
			return err
		}
		defer runtimeApp.Shutdown(nil)
		status, err := runtimeApp.GetAPITokenStatus()
		if err != nil {
			return err
		}
		if !status.HasToken {
			fmt.Fprintln(stdout, "token: not configured")
			return nil
		}
		fmt.Fprintf(stdout, "token: configured\ntoken prefix: %s\n", status.TokenPrefix)
		return nil
	default:
		return fmt.Errorf("unknown auth subcommand %q", subcommand)
	}
}

func openAuthCommandApp(dataPath string) (*irecallapp.App, error) {
	if strings.TrimSpace(dataPath) != "" {
		config.SetRootPath(dataPath)
	} else if _, err := config.ApplyPreferredRootPath(); err != nil {
		return nil, fmt.Errorf("cannot load preferred data root: %w", err)
	}
	if err := config.EnsureDirs(); err != nil {
		return nil, fmt.Errorf("cannot create data directories: %w", err)
	}
	runtimeApp, err := irecallapp.NewApp(config.RootPath())
	if err != nil {
		return nil, err
	}
	return runtimeApp, nil
}

func readAuthCommandPassword(stdin io.Reader, passwordStdin bool) (string, error) {
	if !passwordStdin {
		return "", errors.New("--password-stdin is required for non-interactive auth token commands")
	}
	data, err := io.ReadAll(io.LimitReader(stdin, 64*1024))
	if err != nil {
		return "", fmt.Errorf("read password: %w", err)
	}
	password := strings.TrimRight(string(data), "\r\n")
	if strings.TrimSpace(password) == "" {
		return "", errors.New("password from stdin is empty")
	}
	return password, nil
}

func writeTokenFile(path, token string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("token file path is empty")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create token file directory: %w", err)
	}
	return os.WriteFile(path, []byte(token+"\n"), 0o600)
}
