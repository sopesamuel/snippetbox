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
	Title string `form:"title"`
	Content string `form:"content"`
	Expires int `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name string `form:"name"`
	Email string `form:"email"`
	Password string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email string `form:"email"`
	Password string `form:"password"`
	validator.Validator `form:"-"`
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

	var form snippetcreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
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

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%v", id), http.StatusSeeOther)

}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request){
			data := app.newTemplateData(r)
			data.Form = userSignupForm{}
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request){

	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil{
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "Name", "Field cannot be empty!")
	form.CheckField(validator.NotBlank(form.Email), "Email", "Field cannot be empty!")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "Email", "Invalid email format!")
	form.CheckField(validator.NotBlank(form.Password), "Password", "Field cannot be empty!")
	form.CheckField(validator.MinChars(form.Password, 8), "Password", "Password must be 8 characters long!")

	if !form.Valid(){
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusOK, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {

		if errors.Is(err, models.ErrDuplicateEmail){
			form.AddField("Email", "This email is already in use!")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
			return
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful!, please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}



func (app *application) userLogin(w http.ResponseWriter, r *http.Request){
	
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}

	app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
}
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Authenticate and login the user...")
}
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Logout the user...")
}
