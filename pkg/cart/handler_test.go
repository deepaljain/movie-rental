package cart

import (
    "bytes"
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