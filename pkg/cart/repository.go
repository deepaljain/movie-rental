package cart

import (
	"database/sql"
	"movie-rental/pkg/movies"
)

type Repository interface {
	AddToCart(userID, movieID int) error
	GetCartItems(userID string) ([]movies.Movie, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) AddToCart(userID, movieID int) error {
	_, err := r.db.Exec("INSERT INTO cart (user_id, movie_id) VALUES ($1, $2)", userID, movieID)
	return err
}

func (r *repository) GetCartItems(userID string) ([]movies.Movie, error) {
	rows, err := r.db.Query(`
		SELECT m.movie_id, m.title, m.year, m.plot, m.genre, m.imdbid, m.actors
		FROM cart c
		JOIN movies m ON c.movie_id = m.movie_id
		WHERE c.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []movies.Movie
	for rows.Next() {
		var m movies.Movie
		if err := rows.Scan(&m.MovieID, &m.Title, &m.Year, &m.Plot, &m.Genre, &m.ImdbID, &m.Actors); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, nil
}