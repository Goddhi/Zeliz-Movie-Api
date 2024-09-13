package validator



import (
	"regexp"
)


var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,60}[a-zA-Z0-9])?)*$")
)

// Define a new Validator type which contains a map of validation errors.
type Validaor struct {
	Errors map[string]string
}