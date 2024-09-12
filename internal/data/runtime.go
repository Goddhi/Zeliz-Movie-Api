package data

import (
	"fmt"
	"strconv"
)

// Declare a custom Runtime type, which has the underlying type int32 (the same as our
// Movie struct field).

type Runtime int32


func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	quotedJSONVALUE := strconv.Quote(jsonValue)


	// convert the quoted string value to a byte slice and return it
	return []byte(quotedJSONVALUE), nil

}