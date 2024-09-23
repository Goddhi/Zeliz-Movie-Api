package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/goddhi/zeliz-movie/internal/validator"
	"github.com/lib/pq"
)

//MovieModel struct type which wraps a sql.DB connection pool.
type MovieModel struct {
	DB *sql.DB
}


type Movie struct {
	ID			int64  `json:"id"`   // Unique integer ID for the movie
	CreateAt	time.Time `json:"-"` // Timestamp for when the movie is added to our database
	Title		string  `json:"title"` // Movie title
	Year 		int32  `json:"year,omitempty"`// Movie release year
	Runtime		Runtime `json:"runtime,omitempty"`// Movie runtime (in minutes)
	Genres		[]string `json:"genres,omitempty"`//Slice of genres for the movied (romance, comedy, etc)
	Version		int32 `json:"version"`// time the movie information is updated

}


func (m MovieModel) Insert(movie *Movie) error {

	query := `
				INSERT INTO movies (title, year, runtime, genres)
				VALUES ($1, $2, $3, $4)
				RETURNING id, created_at, version`

			// an args slice containing the values for the placeholder parameters from
			// the movie struct.
	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreateAt, &movie.Version)
}


func (m *MovieModel) Get(id int64) (*Movie, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
				SELECT id, created_at, title, year, runtime, genres, version
				FROM movies
				WHERE id = $1`

	// this holds the data returned by the query
	var movie Movie

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&[]byte{},
		&movie.ID,
		&movie.CreateAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (m *MovieModel) Update(movie *Movie) error {
	query := `
			UPDATE movies
			SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
			WHERE id = $5 AND version = $6
			RETURNING version`


	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,

}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
// Execute the SQL query. If no matching row could be found, we know the movie
// version has changed (or the record has been deleted) and we return our custom
// ErrEditConflict error.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return err
		default:
			return err
		}
	}
	return nil
}

func (m *MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM movies
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}


	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}


func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error) {
	// CSQL query to retrieve all movie records.
	query := `
				SELECT id, created_at, title, year, runtime, genres, version
				FROM movies
				WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
				AND (genres @> $2 OR $2 = '{}')
				ORDER BY id`



	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Pass the title and genres as the placeholder parameter values.
	rows, err := m.DB.QueryContext(ctx, query, title, pq.Array(genres))
	if err != nil {
		return nil, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
// before GetAll() returns.
	defer rows.Close()

	movies := []*Movie{}

	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var movie Movie

		err := rows.Scan(
			&movie.ID,
			&movie.CreateAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)

		if err != nil {
			return nil, err
		}

		// Add the Movie struct to the slice.
		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}


func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}


