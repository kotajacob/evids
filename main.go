package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type application struct {
	content string
	ts      *template.Template

	infoLog *log.Logger
	errLog  *log.Logger
}

type page struct {
	Artist  string
	Entries []entry
}

type entry struct {
	Name string
	Path string
	Size int64
	Time string
}

func (app *application) dir(w http.ResponseWriter, r *http.Request) {
	systempath := filepath.Join(app.content, filepath.Clean(r.URL.Path))
	info, err := os.Stat(systempath)
	if err != nil {
		app.errLog.Println(err)
		http.NotFound(w, r)
		return
	}
	if !info.IsDir() {
		f, err := os.Open(systempath)
		if err != nil {
			app.errLog.Println(err)
			http.NotFound(w, r)
			return
		}
		http.ServeContent(w, r, systempath, info.ModTime(), f)
		return
	}

	dirEntries, err := os.ReadDir(systempath)
	if err != nil {
		app.errLog.Println(err)
		http.NotFound(w, r)
		return
	}

	artist := strings.TrimPrefix(systempath, filepath.Clean(app.content))
	var entries []entry
	for _, e := range dirEntries {
		name := e.Name()
		if len(name) > 28 {
			name = name[:25] + "..."
		}
		if e.IsDir() {
			name += "/"
		}

		info, err := e.Info()
		if err != nil {
			app.errLog.Println(err)
			http.NotFound(w, r)
			return
		}

		entries = append(entries, entry{
			Name: name,
			Path: filepath.Join(artist, e.Name()),
			Size: info.Size(),
			Time: info.ModTime().Format("Jan _2 15:04 2006"),
		})
	}

	err = app.ts.Execute(w, page{
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

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	ui := flag.String("ui", "/etc/evids", "Path to css and ui files")
	content := flag.String("path", "/var/www/html", "Path to serve")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)

	ts := template.Must(template.ParseFS(EmbededFiles, "dir.tmpl"))
	app := &application{
		content: *content,
		ts:      ts,
		infoLog: infoLog,
		errLog:  errLog,
	}

	mux := http.NewServeMux()

	static := http.FileServer(http.Dir(*ui))
	mux.Handle("/static/", http.StripPrefix("/static", static))

	mux.HandleFunc("/", app.dir)

	infoLog.Printf("starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	errLog.Fatal(err)
}
