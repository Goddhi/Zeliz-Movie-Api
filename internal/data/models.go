package data

import (
	"database/sql"
	"errors"
)

//return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict  = errors.New("edit conflict")
)

//Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	Movies MovieModel
}

// For ease of use, NewModels() method which returns a Models struct containing
// the initialized MovieModel
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}