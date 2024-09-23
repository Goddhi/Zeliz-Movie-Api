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

func (app *application) listMovieHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title 			string
		Genres			[]string
		data.Filters

	}

	v := validator.New()

	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})


	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).
	input.Filters.Sort = app.readString(qs, "sort", "id")
	/// sorting based on ascending and descending(-) order
	input.Filters.SortStatelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	// evaluate the validation checks on the filters structs and send a response if it contains an error, if no error it sends the field
	
	if data.ValidateFilters(v, input.Filters); !v.Valid()  {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	
}

