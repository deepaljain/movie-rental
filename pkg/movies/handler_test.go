package movies

import (
	"errors"
    "context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockMovieRepository struct {
    ListMoviesFunc   func(genre, actor, year string) ([]Movie, error)
    GetMovieByIDFunc func(id string) (*Movie, error)
}

func (m *mockMovieRepository) ListMovies(_ctx context.Context, genre, actor, year string) ([]Movie, error) {
    return m.ListMoviesFunc(genre, actor, year)
}
func (m *mockMovieRepository) GetMovieByID(_ctx context.Context, id string) (*Movie, error) {
    return m.GetMovieByIDFunc(id)
}

func setupRouter(repo MovieRepository) *gin.Engine {
    router := gin.Default()
    router.GET("/movies", ListMoviesHandler(repo))
    router.GET("/movies/:id", GetMovieByIDHandler(repo))
    return router
}

func TestListMoviesHandler_ListAll(t *testing.T) {
    gin.SetMode(gin.TestMode)
    repo := &mockMovieRepository{
        ListMoviesFunc: func(genre, actor, year string) ([]Movie, error) {
            return []Movie{
                {MovieID: 1, Title: "Movie 1"},
                {MovieID: 2, Title: "Movie 2"},
            }, nil
        },
    }
    router := setupRouter(repo)

    req, _ := http.NewRequest("GET", "/movies", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestListMoviesHandler_FilterByGenre(t *testing.T) {
    gin.SetMode(gin.TestMode)
    repo := &mockMovieRepository{
        ListMoviesFunc: func(genre, actor, year string) ([]Movie, error) {
            assert.Equal(t, "Action", genre)
            return []Movie{{MovieID: 1, Title: "Movie 1", Genre: "Action"}}, nil
        },
    }
    router := setupRouter(repo)

    req, _ := http.NewRequest("GET", "/movies?genre=Action", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestListMoviesHandler_DBError(t *testing.T) {
    gin.SetMode(gin.TestMode)
    repo := &mockMovieRepository{
        ListMoviesFunc: func(genre, actor, year string) ([]Movie, error) {
            return nil, errors.New("db error")
        },
    }
    router := setupRouter(repo)

    req, _ := http.NewRequest("GET", "/movies", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusInternalServerError, recorder.Code)
    assert.Contains(t, recorder.Body.String(), "db error")
}

func TestGetMovieByIDHandler_Found(t *testing.T) {
    gin.SetMode(gin.TestMode)
    repo := &mockMovieRepository{
        GetMovieByIDFunc: func(id string) (*Movie, error) {
            assert.Equal(t, "1", id)
            return &Movie{MovieID: 1, Title: "Movie 1"}, nil
        },
    }
    router := setupRouter(repo)

    req, _ := http.NewRequest("GET", "/movies/1", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetMovieByIDHandler_NotFound(t *testing.T) {
    gin.SetMode(gin.TestMode)
    repo := &mockMovieRepository{
        GetMovieByIDFunc: func(id string) (*Movie, error) {
            return nil, nil
        },
    }
    router := setupRouter(repo)

    req, _ := http.NewRequest("GET", "/movies/99", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestGetMovieByIDHandler_DBError(t *testing.T) {
    gin.SetMode(gin.TestMode)
    repo := &mockMovieRepository{
        GetMovieByIDFunc: func(id string) (*Movie, error) {
            return nil, errors.New("db failure")
        },
    }
    router := setupRouter(repo)

    req, _ := http.NewRequest("GET", "/movies/1", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusInternalServerError, recorder.Code)
    assert.Contains(t, recorder.Body.String(), "db failure")
}