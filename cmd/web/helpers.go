package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	buf := new(bytes.Buffer)

	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	// Execute the template set, passing in any dynamic data.
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
