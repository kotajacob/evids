package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
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
	err := ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errLog.Println(err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
}
