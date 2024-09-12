package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"fmt"
	"io"

	"github.com/julienschmidt/httprouter"
)


// Retrieve the "id" URL parameter from the current request context, then convert it to
// an integer and return it. If the operation isn't successful, return 0 and an error.
func (app *application) readIDParam(r *http.Request) (int64, error) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

// implementing envelope responses
type envelope map[string]interface{}

//  writeJSON helper method
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	
	js, err := json.MarshalIndent(data, "", "\t") // Here we use no line prefix ("") and tab indents ("\t") for each element.
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	
// We loop through the header map and add each header to the http.ResponseWriter header map
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)


	return nil

}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	//decode the request body into the target destination
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
	
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d", syntaxError.Offset)

		case errors.As(err, io.ErrUnexpectedEOF):
			return errors.New("body contains baly-formed JSON")
		
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		// // For anything else, return the error message as-is.
		
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
			
		
		default: 
			return err
		}

	}
	
	return nil

}



