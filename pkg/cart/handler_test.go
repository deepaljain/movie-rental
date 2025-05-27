package cart

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"movie-rental/pkg/movies"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(db *sql.DB) *gin.Engine {
    router := gin.Default()
    router.POST("/cart", AddToCartHandler(db))
    return router
}

func TestAddToCartHandler_Success(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    mock.ExpectExec("INSERT INTO cart").
        WithArgs(1, 2).
        WillReturnResult(sqlmock.NewResult(1, 1))

    router := setupRouter(db)

	buf := new(bytes.Buffer)
    _ = json.NewEncoder(buf).Encode(AddToCartRequest{UserID: 1, MovieID: 2})

    req, _ := http.NewRequest("POST", "/cart", buf)
    req.Header.Set("Content-Type", "application/json")
    recorder := httptest.NewRecorder()

    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)
    assert.Contains(t, recorder.Body.String(), "Movie added to cart")
}

func TestAddToCartHandler_BadRequest(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db, _, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    router := setupRouter(db)

    req, _ := http.NewRequest("POST", "/cart", bytes.NewBuffer([]byte(`invalid json`)))
    req.Header.Set("Content-Type", "application/json")
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAddToCartHandler_DBError(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    mock.ExpectExec("INSERT INTO cart").
        WithArgs(1, 2).
        WillReturnError(sql.ErrConnDone)

    router := setupRouter(db)

    buf := new(bytes.Buffer)
    _ = json.NewEncoder(buf).Encode(AddToCartRequest{UserID: 1, MovieID: 2})

    req, _ := http.NewRequest("POST", "/cart", buf)
    req.Header.Set("Content-Type", "application/json")
    recorder := httptest.NewRecorder()

    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestViewCartHandler_Success(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    rows := sqlmock.NewRows([]string{"movie_id", "title", "year", "plot", "genre", "imdbid", "actors"}).
        AddRow(1, "Movie 1", 2020, "Plot 1", "Action", "tt1234567", "Actor A, Actor B").
        AddRow(2, "Movie 2", 2021, "Plot 2", "Drama", "tt7654321", "Actor C, Actor D")

    mock.ExpectQuery(`SELECT m.movie_id, m.title, m.year, m.plot, m.genre, m.imdbid, m.actors FROM cart c JOIN movies m ON c.movie_id = m.movie_id WHERE c.user_id = \$1`).
        WithArgs("1").
        WillReturnRows(rows)

    router := gin.Default()
    router.GET("/cart/:user_id", ViewCartHandler(db))

    req, _ := http.NewRequest("GET", "/cart/1", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)

    var resp struct {
        Movies []movies.Movie
    }
    err = json.NewDecoder(recorder.Body).Decode(&resp)
    assert.NoError(t, err)
    assert.Len(t, resp.Movies, 2)
    assert.Equal(t, "Movie 1", resp.Movies[0].Title)
    assert.Equal(t, "Movie 2", resp.Movies[1].Title)
}

func TestViewCartHandler_DBError(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    mock.ExpectQuery(`SELECT m.movie_id, m.title, m.year, m.plot, m.genre, m.imdbid, m.actors FROM cart c JOIN movies m ON c.movie_id = m.movie_id WHERE c.user_id = \$1`).
        WithArgs("1").
        WillReturnError(errors.New("db error"))

    router := gin.Default()
    router.GET("/cart/:user_id", ViewCartHandler(db))

    req, _ := http.NewRequest("GET", "/cart/1", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusInternalServerError, recorder.Code)
    assert.Contains(t, recorder.Body.String(), "db error")
}

func TestViewCartHandler_EmptyCart(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    // Mock empty result set
    rows := sqlmock.NewRows([]string{"movie_id", "title", "year", "plot", "genre", "imdbid", "actors"})
    mock.ExpectQuery(`SELECT m.movie_id, m.title, m.year, m.plot, m.genre, m.imdbid, m.actors FROM cart c JOIN movies m ON c.movie_id = m.movie_id WHERE c.user_id = \$1`).
        WithArgs("1").
        WillReturnRows(rows)

    router := gin.Default()
    router.GET("/cart/:user_id", ViewCartHandler(db))

    req, _ := http.NewRequest("GET", "/cart/1", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)

    var resp struct {
        Movies []movies.Movie
    }
    err = json.NewDecoder(recorder.Body).Decode(&resp)
    assert.NoError(t, err)
    assert.Len(t, resp.Movies, 0)
}