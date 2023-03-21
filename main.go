package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
)

type application struct {
	content string
	ts      *template.Template

	infoLog *log.Logger
	errLog  *log.Logger
}

type ArtistPath struct {
	Artist  string
	Entries []DirEntry
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	content := flag.String("path", "/var/www", "Path to serve")
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

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("starting server on %s", *addr)
	err := srv.ListenAndServe()
	errLog.Fatal(err)
}
