package main

import (
	"fmt"
	"net/http"
	"errors"

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

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new movie in the URL.

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))
	
	// Write a JSON response with a 201 Created status code, the movie data in the
	// response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}


func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

		id, err := app.readIDParam(r)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		movie, err := app.models.Movies.Get(id)
		if err != nil {
			switch  {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)	
			default:
				app.serverErrorResponse(w, r, err)		
			}
			return
		}

		err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
}


func (app *application) updtaeMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Declare an input struct to hold the expected data from the client.
var input struct {
	Title *string `json:"title"`
	Year *int32 `json:"year"`
	Runtime *data.Runtime `json:"runtime"`
	Genres *[]string `json:"genres"`
	}
	
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return 
	}

	// path update implementation
	// If the input.Title value is nil then we know that no corresponding "title" key/
	// value pair was provided in the JSON request body. So we move on and leave the
	// movie record unchanged. Otherwise, we update the movie record with the new title
	// value. Importantly, because input.Title is a now a pointer to a string, we need
	// to dereference the pointer using the * operator to get the underlying value
	// before assigning it to our movie record.


	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}


	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		movie.Genres = *input.Genres
	}

	// Copy the values from the request body to the appropriate fields of the movie
	// record

	// movie.Title = input.Title
	// movie.Year 	= input.Year
	// movie.Runtime = input.Runtime
	// movie.Genres  = input.Genres

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return 
	}

	// Pass the updated movie record to our new Update() method.
	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.EditConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r,)
		return 
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err!= nil {
		app.serverErrorResponse(w, r, err)
	}
}

