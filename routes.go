package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/:path", app.path)
	router.HandlerFunc(http.MethodGet, "/:path/:subpath", app.path)

	return app.logRequest(router)
}

// home is an http.HandlerFunc which displays the home page.
// The home page is a list of all "released" videos, which are loaded on each
// request from "released.txt".
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	released, err := Released(filepath.Join(app.dir, "released.txt"))
	if err != nil {
		app.errLog.Println(err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	tsName := "home.tmpl"
	ts, ok := app.templateCache[tsName]
	if !ok {
		app.errLog.Println(fmt.Errorf(
			"the template %s is missing",
			tsName,
		))
		http.NotFound(w, r)
		return
	}
	err = ts.ExecuteTemplate(w, "base", released)
	if err != nil {
		app.errLog.Println(err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
}

// path is an http.HandlerFunc which passes the request to either artist,
// video, or file depending on if the request is for a file, video file, or
// directory.
func (app *application) path(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(app.dir, filepath.Clean(r.URL.Path))
	info, err := os.Stat(path)
	if err != nil {
		app.errLog.Println(err)
		http.NotFound(w, r)
		return
	}
	if !info.IsDir() {
		q := r.URL.Query()
		_, direct := q["direct"]
		if strings.HasSuffix(path, ".mp4") && !direct {
			app.video(w, r)
			return
		}
		app.file(w, r)
		return
	}

	app.artist(w, r)
}

// ArtistPage is the datastructure used in the artist handler for the artist
// template.
type ArtistPage struct {
	Artist  string
	Entries []DirEntry
}

// artist is an http.HandlerFunc which displays a page for the requested artist.
// The artist's page is a listing of all their videos, including unreleased
// videos.
func (app *application) artist(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(app.dir, filepath.Clean(r.URL.Path))
	artist := strings.TrimPrefix(path, filepath.Clean(app.dir)+"/")
	entries, err := ListVideos(path)
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

// file is an http.HandlerFunc for files.
func (app *application) file(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(app.dir, filepath.Clean(r.URL.Path))
	info, err := os.Stat(path)
	if err != nil {
		app.errLog.Println(err)
		http.NotFound(w, r)
		return
	}
	f, err := os.Open(path)
	if err != nil {
		app.errLog.Println(err)
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, path, info.ModTime(), f)
	return
}

// VideoPage is the datastructure used in the video handler for the video
// template.
type VideoPage struct {
	Artist string
	Video  string
}

// video is an http.HandlerFunc for videos.
// The video is displayed in browser with a helpful player and download link.
func (app *application) video(w http.ResponseWriter, r *http.Request) {
	tsName := "video.tmpl"
	ts, ok := app.templateCache[tsName]
	if !ok {
		app.errLog.Println(fmt.Errorf(
			"the template %s is missing",
			tsName,
		))
		http.NotFound(w, r)
		return
	}
	err := ts.ExecuteTemplate(w, "base", VideoPage{
		Artist: strings.TrimPrefix(
			filepath.Dir(filepath.Clean(r.URL.Path)),
			"/",
		),
		Video: filepath.Clean(r.URL.Path),
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
