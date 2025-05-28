package movies

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func newTestMovieRows() *sqlmock.Rows {
    return sqlmock.NewRows([]string{"movie_id", "title", "year", "plot", "genre", "imdbid", "actors"})
}

func TestListMovies_All(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    rows := newTestMovieRows().
        AddRow(1, "Movie 1", 2020, "Plot 1", "Action", "tt1234567", "Actor A").
        AddRow(2, "Movie 2", 2021, "Plot 2", "Drama", "tt7654321", "Actor B")
    mock.ExpectQuery(`SELECT \* FROM movies WHERE 1=1`).WillReturnRows(rows)

    repo := NewMovieRepository(db)
    movies, err := repo.ListMovies(context.Background(), "", "", "")
    assert.NoError(t, err)
    assert.Len(t, movies, 2)
    assert.Equal(t, "Movie 1", movies[0].Title)
    assert.Equal(t, "Movie 2", movies[1].Title)
}

func TestListMovies_WithGenre(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    rows := newTestMovieRows().
        AddRow(1, "Movie 1", 2020, "Plot 1", "Action", "tt1234567", "Actor A")
    mock.ExpectQuery(`SELECT \* FROM movies WHERE 1=1 AND genre ILIKE '%' \|\| \$1 \|\| '%'`).
        WithArgs("Action").
        WillReturnRows(rows)

    repo := NewMovieRepository(db)
    movies, err := repo.ListMovies(context.Background(), "Action", "", "")
    assert.NoError(t, err)
    assert.Len(t, movies, 1)
    assert.Equal(t, "Action", movies[0].Genre)
}

func TestListMovies_WithActor(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    rows := newTestMovieRows().
        AddRow(2, "Movie 2", 2021, "Plot 2", "Drama", "tt7654321", "Actor B")
    mock.ExpectQuery(`SELECT \* FROM movies WHERE 1=1 AND actors ILIKE '%' \|\| \$1 \|\| '%'`).
        WithArgs("Actor B").
        WillReturnRows(rows)

    repo := NewMovieRepository(db)
    movies, err := repo.ListMovies(context.Background(), "", "Actor B", "")
    assert.NoError(t, err)
    assert.Len(t, movies, 1)
    assert.Equal(t, "Movie 2", movies[0].Title)
}

func TestListMovies_WithYear(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    rows := newTestMovieRows().
        AddRow(3, "Movie 3", 2022, "Plot 3", "Comedy", "tt1111111", "Actor C")
    mock.ExpectQuery(`SELECT \* FROM movies WHERE 1=1 AND year = \$1`).
        WithArgs("2022").
        WillReturnRows(rows)

    repo := NewMovieRepository(db)
    movies, err := repo.ListMovies(context.Background(), "", "", "2022")
    assert.NoError(t, err)
    assert.Len(t, movies, 1)
    assert.Equal(t, 2022, movies[0].Year)
}

func TestListMovies_WithActorAndYear(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    rows := newTestMovieRows().
        AddRow(4, "Movie 4", 2023, "Plot 4", "Thriller", "tt2222222", "Actor D")
    mock.ExpectQuery(`SELECT \* FROM movies WHERE 1=1 AND actors ILIKE '%' \|\| \$1 \|\| '%' AND year = \$2`).
        WithArgs("Actor D", "2023").
        WillReturnRows(rows)

    repo := NewMovieRepository(db)
    movies, err := repo.ListMovies(context.Background(), "", "Actor D", "2023")
    assert.NoError(t, err)
    assert.Len(t, movies, 1)
    assert.Equal(t, "Movie 4", movies[0].Title)
    assert.Equal(t, 2023, movies[0].Year)
}

func TestListMovies_DBError(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    mock.ExpectQuery(`SELECT \* FROM movies WHERE 1=1`).WillReturnError(errors.New("db error"))

    repo := NewMovieRepository(db)
    movies, err := repo.ListMovies(context.Background(), "", "", "")
    assert.Error(t, err)
    assert.Nil(t, movies)
}

func TestListMovies_RowScanError(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    rows := newTestMovieRows().
        AddRow("not-an-int", "Title", 2020, "Plot", "Genre", "imdbid", "Actors")
    mock.ExpectQuery(`SELECT \* FROM movies WHERE 1=1`).WillReturnRows(rows)

    repo := NewMovieRepository(db)
    movies, err := repo.ListMovies(context.Background(), "", "", "")
    assert.Error(t, err)
    assert.Nil(t, movies)
}

func TestGetMovieByID_Found(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    row := newTestMovieRows().
        AddRow(1, "Movie 1", 2020, "Plot 1", "Action", "tt1234567", "Actor A")
    mock.ExpectQuery(`SELECT \* FROM movies WHERE movie_id = \$1`).
        WithArgs("1").
        WillReturnRows(row)

    repo := NewMovieRepository(db)
    movie, err := repo.GetMovieByID(context.Background(), "1")
    assert.NoError(t, err)
    assert.NotNil(t, movie)
    assert.Equal(t, "Movie 1", movie.Title)
}

func TestGetMovieByID_NotFound(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    mock.ExpectQuery(`SELECT \* FROM movies WHERE movie_id = \$1`).
        WithArgs("99").
        WillReturnRows(newTestMovieRows())

    repo := NewMovieRepository(db)
    movie, err := repo.GetMovieByID(context.Background(), "99")
    assert.NoError(t, err)
    assert.Nil(t, movie)
}

func TestGetMovieByID_DBError(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    mock.ExpectQuery(`SELECT \* FROM movies WHERE movie_id = \$1`).
        WithArgs("1").
        WillReturnError(errors.New("db failure"))

    repo := NewMovieRepository(db)
    movie, err := repo.GetMovieByID(context.Background(), "1")
    assert.Error(t, err)
    assert.Nil(t, movie)
}