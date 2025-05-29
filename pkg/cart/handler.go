package cart

import (
	"movie-rental/pkg/movies"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddToCartHandler(repo Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AddToCartRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := repo.AddToCart(req.UserID, req.MovieID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Movie added to cart"})
	}
}

func ViewCartHandler(repo Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")

		moviesList, err := repo.GetCartItems(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
        if len(moviesList) == 0 {
            c.JSON(http.StatusOK, gin.H{"movies": []movies.Movie{}})
            return
        }

		c.JSON(http.StatusOK, gin.H{"movies": moviesList})
	}
}