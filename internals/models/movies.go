package models

import (
	"database/sql"
	"errors"
	"time"
)

type MovieModelInterface interface {
	Insert(title string, synopsis string, rating float64) (int, error)
	Get(id int) (*Movie, error)
	Latest() ([]*Movie, error)
}

type Movie struct {
	ID       int
	Title    string
	Synopsis string
	Created  time.Time
	Rating   float64
}

type MovieModel struct {
	DB *sql.DB
}

// Insert into the Movie table

func (m *MovieModel) Insert(title string, synopsis string, rating float64) (int, error) {
	stmt := `INSERT INTO movies (title, synopsis, created, rating) VALUES(?, ?, UTC_TIMESTAMP(), ?))`
	result, err := m.DB.Exec(stmt, title, synopsis, rating)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Get a movie by id

func (m *MovieModel) Get(id int) (*Movie, error) {
	s := &Movie{}
	err := m.DB.QueryRow(`SELECT id, title, synopsis, created, rating FROM movies
								WHERE id = ?`, id).Scan(&s.ID, &s.Title, &s.Synopsis, &s.Created, &s.Rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no movie found")
		} else {
			return nil, err
		}
	}
	return s, nil
}

// This will return the 10 most recently created movies.

func (m *MovieModel) Latest() ([]*Movie, error) {
	stmt := `SELECT id, title, synopsis, created, rating FROM snippets ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*Movie{}
	for rows.Next() {
		s := &Movie{}
		err = rows.Scan(&s.ID, &s.Title, &s.Synopsis, &s.Created, &s.Rating)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
