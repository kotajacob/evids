package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/:artist", app.artist)
	router.HandlerFunc(http.MethodGet, "/:artist/:file", app.artist)

	return app.logRequest(router)
}
