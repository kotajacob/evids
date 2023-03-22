package main

import (
	"embed"
	"io/fs"
	"path/filepath"
	"text/template"
)

//go:embed "ui"
var EmbededFiles embed.FS

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(EmbededFiles, "ui/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		files := []string{
			"ui/base.tmpl",
			page,
		}

		ts, err := template.ParseFS(EmbededFiles, files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
