package cart

import (
	"database/sql"
	"movie-rental/pkg/movies"
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

func ViewCartHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Param("user_id")
        rows, err := db.Query(`
            SELECT m.movie_id, m.title, m.year, m.plot, m.genre, m.imdbid, m.actors
            FROM cart c
            JOIN movies m ON c.movie_id = m.movie_id
            WHERE c.user_id = $1
        `, userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer rows.Close()

        var moviesList []movies.Movie
        for rows.Next() {
            var m movies.Movie
            if err := rows.Scan(&m.MovieID, &m.Title, &m.Year, &m.Plot, &m.Genre, &m.ImdbID, &m.Actors); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            moviesList = append(moviesList, m)
        }
		if len(moviesList) == 0 {
			c.JSON(http.StatusOK, gin.H{"movies": []movies.Movie{}})	
		} else {
			c.JSON(http.StatusOK, gin.H{"movies": moviesList})
		}
    }
}