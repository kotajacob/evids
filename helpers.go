package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("home"))
}

func (app *application) artist(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(app.content, filepath.Clean(r.URL.Path))
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

	artist := strings.TrimPrefix(path, filepath.Clean(app.content))
	entries, err := ListDir(path)
	if err != nil {
		app.errLog.Println(err)
		http.NotFound(w, r)
		return
	}
	err = app.ts.Execute(w, ArtistPath{
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
