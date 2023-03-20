package main

import (
	"embed"
)

//go:embed "dir.tmpl"
var EmbededFiles embed.FS
