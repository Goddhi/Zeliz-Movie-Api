package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/goddhi/zeliz-movie/internal/data"
	"github.com/goddhi/zeliz-movie/internal/validator"
)


func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	
	var input struct {
		Title 	string			`json:"title"`
		Year	int32			`json:"year"`
		Runtime data.Runtime	`json:"runtime"`
		Genres  []string		`json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title: 		input.Title,
		Year: 		input.Year,
		Runtime: 	input.Runtime,
		Genres: 	input.Genres,
	}

	// initialized a validator instance from the validator 
	v := validator.New()
	
	// Call the ValidateMovie() function and return a response containing the errors if
	// any of the checks fail.

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
 	// the Valid() method to see if any of the checks failed. If they did, then use
// the failedValidationResponse() helper to send a response to the client,

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

		id, err := app.readIDParam(r)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		movie := data.Movie{
			ID:			id,
			CreateAt: 	time.Now(),
			Title: 		"The Walking Dead",
			Runtime: 	102,
			Genres: 	[]string{"drama", "romance", "war"},
			Version: 	1,
		}

		err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
}

