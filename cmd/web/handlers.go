package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/swaggest/form/v5"
	"net/http"
	"snippetbox.vadimpush.net/internal/models"
	"snippetbox.vadimpush.net/internal/validator"
	"strconv"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serveError(w, err)
		return
	}

	data := app.newTemplateData(r)

	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serveError(w, err)
		}

		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.html", data)

}
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		//Expires: 365,
	}

	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()

	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)

	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}
	return nil
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	var form snippetCreateForm

	err := app.decodePostForm(r, &form)

	if err != nil {
		app.clientErr(w, http.StatusBadRequest)
		return
	}

	err = app.formDecoder.Decode(&form, r.PostForm)

	if err != nil {
		app.clientErr(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", blankMsg)

	form.CheckField(validator.MaxChars(form.Title, 100), "title", longCharMsg)

	form.CheckField(validator.NotBlank(form.Content), "content", blankMsg)

	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", expiresMsg)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.serveError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
