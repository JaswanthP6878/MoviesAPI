package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"greenlight.jaswanthp.com/internal/data"
)

func (app *application) readParamID(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, err
	}

	return id, nil

}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	// declaring a variable with the required struct type
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.decodeJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	fmt.Printf("%+v\n", input)

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	data := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Godfather",
		Year:      1973,
		Runtime:   120,
		Genres:    []string{"gangster"},
		Version:   1,
	}

	err = app.writeJson(w, http.StatusOK, envolope{"movie": data}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
