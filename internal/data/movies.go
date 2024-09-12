package data

import (
	"time"
)

type Movie struct {
	ID			int64  `json:"id"`   // Unique integer ID for the movie
	CreateAt	time.Time `json:"-"` // Timestamp for when the movie is added to our database
	Title		string  `json:"title"` // Movie title
	Year 		int32  `json:"year,omitempty"`// Movie release year
	Runtime		Runtime `json:"runtime,omitempty"`// Movie runtime (in minutes)
	Genres		[]string `json:"genres,omitempty"`//Slice of genres for the movied (romance, comedy, etc)
	Version		int32 `json:"version"`// time the movie information is updated

}

