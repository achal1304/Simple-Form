package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	//standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/show", http.HandlerFunc(app.showSnippet))

	// mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/create", http.HandlerFunc(app.createSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
