package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"fmt"
	"io"
	"strings"
	"net/url"

	"github.com/julienschmidt/httprouter"
	"github.com/goddhi/zeliz-movie/internal/validator"


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
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB.

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
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()  // This means that if the JSON from the client now includes any
	// field which cannot be mapped to the target destination, the decoder will return
	// an error instead of just ignoring the field.
	err := dec.Decode(dst)	// // Decode the request body to the destination.
	if err != nil {
	// If there is an error during decoding, start the triage...
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError
	switch {

	case errors.As(err, &syntaxError):
	return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

	case errors.Is(err, io.ErrUnexpectedEOF):
	return errors.New("body contains badly-formed JSON")

	case errors.As(err, &unmarshalTypeError):
	if unmarshalTypeError.Field != "" {
	return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
	}
	return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

	case errors.Is(err, io.EOF):
	return errors.New("body must not be empty")


	case strings.HasPrefix(err.Error(), "json: unknow field"):
		fileName := strings.TrimPrefix(err.Error(), "json: unknown field")
		return fmt.Errorf("body contains unknown key %s", fileName)

	case err.Error() == "http: request body too large":
		return fmt.Errorf("body must be larger than %d bytes", maxBytes)
 	
	
	case errors.As(err, &invalidUnmarshalError):
		panic(err)  // this error is likely to occur in the development stage due to not assiging pointers to the struct 
		// For anything else, return the error message as-is.	
	
	default:
	return err

	}
}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
return nil
}

// The readString() helper returns a string value from the query string, or the provided
// default value if no matching key could be found
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	
	s := qs.Get(key) // extacts the value for a given key from the query string, if no key exists this will return the enoty string ""

	if s == "" {
		return defaultValue
	}

	return s 
}


// The readCSV() helper reads a string value from the query string and then splits it
// into a slice on the comma character. If no matching key could be found, it returns
// the provided default value.

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	// Otherwise parse the value into a []string slice and return it.
	return strings.Split(csv, ",")

} 

// The readInt() helper reads a string value from the query string and converts it to an
// integer before returning. If no matching key could be found it returns the provided
// default value. If the value couldn't be converted to an integer, then we record an
// error message in the provided Validator instance.
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be integer vallue")
		return defaultValue
	}

	return i
}
