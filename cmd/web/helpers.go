package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serveError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%e\n%s", err, debug.Stack())
	app.errLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientErr(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientErr(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {

	tmpl, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)

		app.serveError(w, err)

		return
	}

	buf := new(bytes.Buffer)

	if err := tmpl.ExecuteTemplate(buf, "base", data); err != nil {
		app.serveError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)

}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		Flash:       app.sessionManager.PopString(r.Context(), "flash"),
	}
}
