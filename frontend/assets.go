package frontend

import "embed"

// Assets bundles the built frontend for both the desktop and web runtimes.
//
//go:embed all:dist
var Assets embed.FS
