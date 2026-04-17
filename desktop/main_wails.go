//go:build wails

package main

import (
	"log"

	"github.com/gigol/irecall/app"
	frontendassets "github.com/gigol/irecall/frontend"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	windowsoptions "github.com/wailsapp/wails/v2/pkg/options/windows"
)

func main() {
	runtimeApp, err := app.NewApp("")
	if err != nil {
		log.Fatal(err)
	}

	err = wails.Run(&options.App{
		Title:            "iRecall",
		Width:            1360,
		Height:           860,
		MinWidth:         1100,
		MinHeight:        760,
		DisableResize:    false,
		AssetServer:      &assetserver.Options{Assets: frontendassets.Assets},
		OnStartup:        runtimeApp.Startup,
		OnShutdown:       runtimeApp.Shutdown,
		BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 1},
		Bind: []interface{}{
			runtimeApp,
		},
		Windows: &windowsoptions.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			DisableFramelessWindowDecorations: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
