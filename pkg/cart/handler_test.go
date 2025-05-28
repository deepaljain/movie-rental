package cart

import (
    "bytes"
    "encoding/json"
    "errors"
    "movie-rental/pkg/movies"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

// Mock Repository for handler tests
type mockRepository struct {
    AddToCartFunc   func(userID, movieID int) error
    GetCartItemsFunc func(userID string) ([]movies.Movie, error)
}

func (m *mockRepository) AddToCart(userID, movieID int) error {
    return m.AddToCartFunc(userID, movieID)
}
func (m *mockRepository) GetCartItems(userID string) ([]movies.Movie, error) {
    return m.GetCartItemsFunc(userID)
}

func setupRouter(repo Repository) *gin.Engine {
    router := gin.Default()
    router.POST("/cart", AddToCartHandler(repo))
    router.GET("/cart/:user_id", ViewCartHandler(repo))
    return router
}

func TestAddToCartHandler_Success(t *testing.T) {
    gin.SetMode(gin.TestMode)
    repo := &mockRepository{
        AddToCartFunc: func(userID, movieID int) error {
            assert.Equal(t, 1, userID)
            assert.Equal(t, 2, movieID)
            return nil
        },
    }
    router := setupRouter(repo)

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
    repo := &mockRepository{}
    router := setupRouter(repo)

    req, _ := http.NewRequest("POST", "/cart", bytes.NewBuffer([]byte(`invalid json`)))
    req.Header.Set("Content-Type", "application/json")
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAddToCartHandler_DBError(t *testing.T) {
    gin.SetMode(gin.TestMode)
    repo := &mockRepository{
        AddToCartFunc: func(userID, movieID int) error {
            return errors.New("db error")
        },
    }
    router := setupRouter(repo)

    buf := new(bytes.Buffer)
    _ = json.NewEncoder(buf).Encode(AddToCartRequest{UserID: 1, MovieID: 2})

    req, _ := http.NewRequest("POST", "/cart", buf)
    req.Header.Set("Content-Type", "application/json")
    recorder := httptest.NewRecorder()

    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusInternalServerError, recorder.Code)
    assert.Contains(t, recorder.Body.String(), "db error")
}

func TestViewCartHandler_Success(t *testing.T) {
    gin.SetMode(gin.TestMode)
    repo := &mockRepository{
        GetCartItemsFunc: func(userID string) ([]movies.Movie, error) {
            assert.Equal(t, "1", userID)
            return []movies.Movie{
                {MovieID: 1, Title: "Movie 1"},
                {MovieID: 2, Title: "Movie 2"},
            }, nil
        },
    }
    router := setupRouter(repo)

    req, _ := http.NewRequest("GET", "/cart/1", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)

    var resp struct {
        Movies []movies.Movie `json:"movies"`
    }
    err := json.NewDecoder(recorder.Body).Decode(&resp)
    assert.NoError(t, err)
    assert.Len(t, resp.Movies, 2)
    assert.Equal(t, "Movie 1", resp.Movies[0].Title)
    assert.Equal(t, "Movie 2", resp.Movies[1].Title)
}

func TestViewCartHandler_DBError(t *testing.T) {
    gin.SetMode(gin.TestMode)
    repo := &mockRepository{
        GetCartItemsFunc: func(userID string) ([]movies.Movie, error) {
            return nil, errors.New("db error")
        },
    }
    router := setupRouter(repo)

    req, _ := http.NewRequest("GET", "/cart/1", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusInternalServerError, recorder.Code)
    assert.Contains(t, recorder.Body.String(), "db error")
}

func TestViewCartHandler_EmptyCart(t *testing.T) {
    gin.SetMode(gin.TestMode)
    repo := &mockRepository{
        GetCartItemsFunc: func(userID string) ([]movies.Movie, error) {
            return []movies.Movie{}, nil
        },
    }
    router := setupRouter(repo)

    req, _ := http.NewRequest("GET", "/cart/1", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)

    var resp struct {
        Movies []movies.Movie `json:"movies"`
    }
    err := json.NewDecoder(recorder.Body).Decode(&resp)
    assert.NoError(t, err)
    assert.Len(t, resp.Movies, 0)
}