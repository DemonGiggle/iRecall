//go:build !wails

package app

import "errors"

func (a *App) SelectQuoteImportFile() (string, error) {
	return "", errors.New("file dialogs require the Wails desktop build")
}

func (a *App) SelectQuoteExportFile() (string, error) {
	return "", errors.New("file dialogs require the Wails desktop build")
}

func (a *App) SelectRootDir() (string, error) {
	return "", errors.New("file dialogs require the Wails desktop build")
}
