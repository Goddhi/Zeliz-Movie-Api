package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// error that our UnmarshalJSON() method can return if we're unable to parse
// or convert the JSON string successfully.
var ErrInvalidRUntimeFormat = errors.New("invalid runtime format")
// Declare a custom Runtime type, which has the underlying type int32 (the same as our
// Movie struct field).

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	quotedJSONVALUE := strconv.Quote(jsonValue)


	// convert the quoted string value to a byte slice and return it
	return []byte(quotedJSONVALUE), nil

}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {

	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRUntimeFormat
	}
	// // Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedJSONValue, " ")

	// // Sanity check the parts of the string to make sure it was in the expected format.
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRUntimeFormat
	}

	// Otherwise, parse the string containing the number into an int32. Again, if this
	// fails return the ErrInvalidRuntimeFormat error.
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRUntimeFormat


	}

	// Convert the int32 to a Runtime type and assign this to the receiver. Note that we
	// use the * operator to deference the receiver (which is a pointer to a Runtime
	// type) in order to set the underlying value of the pointer.
	*r = Runtime(i)

	return nil


}


