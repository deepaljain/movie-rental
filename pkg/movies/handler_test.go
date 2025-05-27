package movies

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestListMoviesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create sqlmock database connection and mock object
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock rows to return
	rows := sqlmock.NewRows([]string{"movie_id", "title", "year", "plot", "genre", "imdbid"}).
		AddRow(1, "Movie 1", 2020, "Plot 1", "Action", "tt1234567").
		AddRow(2, "Movie 2", 2021, "Plot 2", "Drama", "tt7654321")

	mock.ExpectQuery("SELECT \\* FROM movies").WillReturnRows(rows)

	// Create a test router and register the handler
	router := gin.Default()
	router.GET("/movies", ListMoviesHandler(db))

	// Create a test HTTP request
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