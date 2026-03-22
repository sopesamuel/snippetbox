package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)


func (app *application) newTemplateData(r *http.Request)templateData{
	return templateData{
		CurrentYear : time.Now().Year(),
		Flash : app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken: nosurf.Token(r),
	}
}


func (app *application) decodePostForm(r *http.Request, dst any) error {

	err := r.ParseForm()
	if err != nil{
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		var EmptyDecodeError form.DecodeErrors
		if errors.As(err, &EmptyDecodeError) {
			return nil
		}

		return err
	}

	return nil
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData){

	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("Page %s not available!", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil{
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
	
}

//Error from our server
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error){

	var (
		method = r.Method
		uri = r.URL.RequestURI()
		trace = string(debug.Stack())
	)
	
	app.logger.Error(err.Error(), "method", method, "uri", uri )

	if app.debug {
		body := fmt.Sprintf("%s\n%s", err, trace)
		http.Error(w, body, http.StatusInternalServerError)
		return
	}

	//Status text receives the number and converts it to human text version i.e bad request
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}


//Error from the client or user
func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int){
	http.Error(w, http.StatusText(status), status)
}

func (app *application) isAuthenticated(r *http.Request) bool{
	isAuthenticated , ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}