package hello

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

func TestHelloHandler(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.Default()
    router.GET("/hello", HelloHandler)

    req, _ := http.NewRequest("GET", "/hello", nil)
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    if recorder.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", recorder.Code)
    }
}