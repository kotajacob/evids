package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (app *application) artist(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(app.dir, filepath.Clean(r.URL.Path))
	info, err := os.Stat(path)
	if err != nil {
		app.errLog.Println(err)
		http.NotFound(w, r)
		return
	}
	if !info.IsDir() {
		f, err := os.Open(path)
		if err != nil {
			app.errLog.Println(err)
			http.NotFound(w, r)
			return
		}
		http.ServeContent(w, r, path, info.ModTime(), f)
		return
	}

	artist := strings.TrimPrefix(path, filepath.Clean(app.dir)+"/")
	entries, err := ListDir(path)
	if err != nil {
		app.errLog.Println(err)
		http.NotFound(w, r)
		return
	}

	tsName := "artist.tmpl"
	ts, ok := app.templateCache[tsName]
	if !ok {
		app.errLog.Println(fmt.Errorf(
			"the template %s is missing",
			tsName,
		))
		http.NotFound(w, r)
		return
	}
	err = ts.ExecuteTemplate(w, "base", ArtistPage{
		Artist:  artist,
		Entries: entries,
	})
	if err != nil {
		app.errLog.Println(err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
}
