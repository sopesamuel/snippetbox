package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"snippetbox.project.sope/internal/models"
	"snippetbox.project.sope/internal/validator"
)

type snippetcreateForm struct {
	Title string
	Content string
	Expires int
	validator.Validator
}


func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl", data)

}


func (app *application) snippetView(w http.ResponseWriter, r *http.Request){
	id, err:= strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1{
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord){
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return 
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", data)
}


func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request){
	data := app.newTemplateData(r)
	data.Form = snippetcreateForm{}
	app.render(w, r, http.StatusOK, "create.tmpl", data)

}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request){

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		expires = 0
	}

	form := snippetcreateForm{
		Title: r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}
	

	form.CheckField(validator.NotBlank(form.Title), "Title", "Field cannot be empty!")
	form.CheckField(validator.MaxChar(form.Title, 100), "Title", "Value characters cannot exceed 100!")
	form.CheckField(validator.NotBlank(form.Content), "Content", "Field cannot be empty!")
	form.CheckField(validator.PermittedValue(form.Expires, 1,7, 365), "Expires", "Invalid value, expiry date value has to be 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil{
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%v", id), http.StatusSeeOther)

}
