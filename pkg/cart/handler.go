package cart

import (
    "database/sql"
    "net/http"
    "github.com/gin-gonic/gin"
)

func AddToCartHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req AddToCartRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
            return
        }
        _, err := db.Exec("INSERT INTO cart (user_id, movie_id) VALUES ($1, $2)", req.UserID, req.MovieID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "Movie added to cart"})
    }
}