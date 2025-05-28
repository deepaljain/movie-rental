package movies

import (
    "context"
    "database/sql"
    "fmt"
)

type MovieRepository interface {
    ListMovies(ctx context.Context, genre, actor, year string) ([]Movie, error)
    GetMovieByID(ctx context.Context, id string) (*Movie, error)
}

type movieRepository struct {
    db *sql.DB
}

func NewMovieRepository(db *sql.DB) MovieRepository {
    return &movieRepository{db: db}
}

func (r *movieRepository) ListMovies(ctx context.Context, genre, actor, year string) ([]Movie, error) {
    query := "SELECT * FROM movies WHERE 1=1"
    var args []interface{}
    idx := 1

    if genre != "" {
        query += fmt.Sprintf(" AND genre ILIKE '%%' || $%d || '%%'", idx)
        args = append(args, genre)
        idx++
    }
    if actor != "" {
        query += fmt.Sprintf(" AND actors ILIKE '%%' || $%d || '%%'", idx)
        args = append(args, actor)
        idx++
    }
    if year != "" {
        query += fmt.Sprintf(" AND year = $%d", idx)
        args = append(args, year)
        idx++
    }

    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var movies []Movie
    for rows.Next() {
        var m Movie
        if err := rows.Scan(&m.MovieID, &m.Title, &m.Year, &m.Plot, &m.Genre, &m.ImdbID, &m.Actors); err != nil {
            return nil, err
        }
        movies = append(movies, m)
    }

    return movies, nil
}

func (r *movieRepository) GetMovieByID(ctx context.Context, id string) (*Movie, error) {
    var m Movie
    err := r.db.QueryRowContext(ctx,
        "SELECT * FROM movies WHERE movie_id = $1", id,
    ).Scan(&m.MovieID, &m.Title, &m.Year, &m.Plot, &m.Genre, &m.ImdbID, &m.Actors)

    if err == sql.ErrNoRows {
        return nil, nil
    } else if err != nil {
        return nil, err
    }

    return &m, nil
}