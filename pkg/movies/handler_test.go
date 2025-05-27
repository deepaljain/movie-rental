package movies

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()
	router.GET("/movies", ListMoviesHandler(db))
    router.GET("/movies/:id", GetMovieByIDHandler(db))
	return router
}

func newMovieRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"movie_id", "title", "year", "plot", "genre", "imdbid", "actors"})
}

func TestListMoviesHandler_ListAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := newMovieRows().
		AddRow(1, "Movie 1", 2020, "Plot 1", "Action", "tt1234567", "Actor A").
		AddRow(2, "Movie 2", 2021, "Plot 2", "Drama", "tt7654321", "Actor B")

	mock.ExpectQuery(`SELECT \* FROM movies WHERE 1=1`).WillReturnRows(rows)

	router := setupRouter(db)

	req, _ := http.NewRequest("GET", "/movies", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var movies []Movie
	err = json.NewDecoder(recorder.Body).Decode(&movies)
	assert.NoError(t, err)
	assert.Len(t, movies, 2)
	assert.Equal(t, "Movie 1", movies[0].Title)
	assert.Equal(t, "Movie 2", movies[1].Title)
}

func TestListMoviesHandler_FilterByGenre(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    router := setupRouter(db)

    rows := newMovieRows().
        AddRow(1, "Movie 1", 2020, "Plot 1", "Action", "tt1234567", "Actor A, Actor B")
    mock.ExpectQuery(`SELECT \* FROM movies WHERE 1=1 AND genre ILIKE '%' \|\| \$1 \|\| '%'`).
        WithArgs("Action").
        WillReturnRows(rows)

    req, _ := http.NewRequest("GET", "/movies?genre=Action", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)

    var movies []Movie
    err = json.NewDecoder(recorder.Body).Decode(&movies)
    assert.NoError(t, err)
    assert.Len(t, movies, 1)
    assert.Equal(t, "Movie 1", movies[0].Title)
}

func TestListMoviesHandler_FilterByActor(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    router := setupRouter(db)

    rows := newMovieRows().
        AddRow(2, "Movie 2", 2021, "Plot 2", "Drama", "tt7654321", "Actor C, Actor D")
    mock.ExpectQuery(`SELECT \* FROM movies WHERE 1=1 AND actors ILIKE '%' \|\| \$1 \|\| '%'`).
        WithArgs("Actor C").
        WillReturnRows(rows)

    req, _ := http.NewRequest("GET", "/movies?actor=Actor C", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)

    var movies []Movie
    err = json.NewDecoder(recorder.Body).Decode(&movies)
    assert.NoError(t, err)
    assert.Len(t, movies, 1)
    assert.Equal(t, "Movie 2", movies[0].Title)
}

func TestListMoviesHandler_FilterByYear(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    router := setupRouter(db)

    rows := newMovieRows().
        AddRow(3, "Movie 3", 2022, "Plot 3", "Comedy", "tt9999999", "Actor E, Actor F")
    mock.ExpectQuery(`SELECT \* FROM movies WHERE 1=1 AND year = \$1`).
        WithArgs("2022").
        WillReturnRows(rows)

    req, _ := http.NewRequest("GET", "/movies?year=2022", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)

    var movies []Movie
    err = json.NewDecoder(recorder.Body).Decode(&movies)
    assert.NoError(t, err)
    assert.Len(t, movies, 1)
    assert.Equal(t, "Movie 3", movies[0].Title)
}

func TestListMoviesHandler_FilterByActorAndYear(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    router := setupRouter(db)

    rows := newMovieRows().
        AddRow(4, "Movie 4", 2023, "Plot 4", "Thriller", "tt8888888", "Actor G, Actor H")
    mock.ExpectQuery(`SELECT \* FROM movies WHERE 1=1 AND actors ILIKE '%' \|\| \$1 \|\| '%' AND year = \$2`).
        WithArgs("Actor G", "2023").
        WillReturnRows(rows)

    req, _ := http.NewRequest("GET", "/movies?actor=Actor G&year=2023", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)

    var movies []Movie
    err = json.NewDecoder(recorder.Body).Decode(&movies)
    assert.NoError(t, err)
    assert.Len(t, movies, 1)
    assert.Equal(t, "Movie 4", movies[0].Title)
}

func TestGetMovieByIDHandler_Found(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    router := setupRouter(db)

    row := newMovieRows().
        AddRow(1, "Movie 1", 2020, "Plot 1", "Action", "tt1234567", "Actor A, Actor B")
    mock.ExpectQuery(`SELECT \* FROM movies WHERE movie_id = \$1`).
        WithArgs("1").WillReturnRows(row)

    req, _ := http.NewRequest("GET", "/movies/1", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)

    var movie Movie
    err = json.NewDecoder(recorder.Body).Decode(&movie)
    assert.NoError(t, err)
    assert.Equal(t, "Movie 1", movie.Title)
}

func TestGetMovieByIDHandler_NotFound(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    router := setupRouter(db)

    mock.ExpectQuery(`SELECT \* FROM movies WHERE movie_id = \$1`).
        WithArgs("99").WillReturnRows(newMovieRows())

    req, _ := http.NewRequest("GET", "/movies/99", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusNotFound, recorder.Code)
}