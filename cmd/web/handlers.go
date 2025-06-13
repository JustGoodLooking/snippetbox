package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"snippetbox.justgoodlooking.com/internal/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippet.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	for _, snippet := range snippets {
		data, _ := json.MarshalIndent(snippet, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippet.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrorNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Write the snippet data as a plain-text HTTP response body.
	//fmt.Fprintf(w, "%+v", snippet)
	data, _ := json.MarshalIndent(snippet, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippet.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
