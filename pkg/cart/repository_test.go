package cart

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddToCart_Success(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    mock.ExpectExec("INSERT INTO cart").
        WithArgs(1, 2).
        WillReturnResult(sqlmock.NewResult(1, 1))

    repo := NewRepository(db)
    err = repo.AddToCart(1, 2)
    assert.NoError(t, err)
}

func TestAddToCart_DBError(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    mock.ExpectExec("INSERT INTO cart").
        WithArgs(1, 2).
        WillReturnError(errors.New("db error"))

    repo := NewRepository(db)
    err = repo.AddToCart(1, 2)
    assert.Error(t, err)
}

func TestGetCartItems_Success(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    rows := sqlmock.NewRows([]string{"movie_id", "title", "year", "plot", "genre", "imdbid", "actors"}).
        AddRow(1, "Movie 1", 2020, "Plot 1", "Action", "tt1234567", "Actor A, Actor B").
        AddRow(2, "Movie 2", 2021, "Plot 2", "Drama", "tt7654321", "Actor C, Actor D")

    mock.ExpectQuery(`SELECT m.movie_id, m.title, m.year, m.plot, m.genre, m.imdbid, m.actors FROM cart c JOIN movies m ON c.movie_id = m.movie_id WHERE c.user_id = \$1`).
        WithArgs("1").
        WillReturnRows(rows)

    repo := NewRepository(db)
    moviesList, err := repo.GetCartItems("1")
    assert.NoError(t, err)
    assert.Len(t, moviesList, 2)
    assert.Equal(t, "Movie 1", moviesList[0].Title)
    assert.Equal(t, "Movie 2", moviesList[1].Title)
}

func TestGetCartItems_DBError(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    mock.ExpectQuery(`SELECT m.movie_id, m.title, m.year, m.plot, m.genre, m.imdbid, m.actors FROM cart c JOIN movies m ON c.movie_id = m.movie_id WHERE c.user_id = \$1`).
        WithArgs("1").
        WillReturnError(errors.New("db error"))

    repo := NewRepository(db)
    moviesList, err := repo.GetCartItems("1")
    assert.Error(t, err)
    assert.Nil(t, moviesList)
}

func TestGetCartItems_RowScanError(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    rows := sqlmock.NewRows([]string{"movie_id", "title", "year", "plot", "genre", "imdbid", "actors"}).
        AddRow("not-an-int", "Title", 2020, "Plot", "Genre", "imdbid", "Actors")
    mock.ExpectQuery(`SELECT m.movie_id, m.title, m.year, m.plot, m.genre, m.imdbid, m.actors FROM cart c JOIN movies m ON c.movie_id = m.movie_id WHERE c.user_id = \$1`).
        WithArgs("1").
        WillReturnRows(rows)

    repo := NewRepository(db)
    moviesList, err := repo.GetCartItems("1")
    assert.Error(t, err)
    assert.Nil(t, moviesList)
}