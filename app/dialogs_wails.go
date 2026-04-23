//go:build wails

package app

import (
	"errors"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) SelectQuoteImportFile() (string, error) {
	if a.ctx == nil {
		return "", errors.New("desktop runtime is not available")
	}
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Import iRecall Quotes",
		Filters: []runtime.FileFilter{
			{DisplayName: "iRecall Share Files", Pattern: "*.json"},
			{DisplayName: "All Files", Pattern: "*"},
		},
	})
}

func (a *App) SelectQuoteExportFile() (string, error) {
	if a.ctx == nil {
		return "", errors.New("desktop runtime is not available")
	}
	return runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Export iRecall Quotes",
		DefaultFilename: "irecall-share.json",
		DefaultDirectory: func() string {
			if a.paths.RootDir == "" {
				return ""
			}
			return filepath.Join(a.paths.RootDir, "exports")
		}(),
		Filters: []runtime.FileFilter{
			{DisplayName: "JSON Files", Pattern: "*.json"},
			{DisplayName: "All Files", Pattern: "*"},
		},
	})
}

func (a *App) SelectRootDir() (string, error) {
	if a.ctx == nil {
		return "", errors.New("desktop runtime is not available")
	}
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Choose iRecall Root Directory",
		DefaultDirectory: func() string {
			if a.paths.RootDir != "" {
				return a.paths.RootDir
			}
			return a.paths.DataDir
		}(),
	})
}
