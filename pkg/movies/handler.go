package movies

import (
    "database/sql"
    "net/http"
    "github.com/gin-gonic/gin"
)

func ListMoviesHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        rows, err := db.Query("SELECT * FROM movies")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer rows.Close()

        var movies []Movie
        for rows.Next() {
            var m Movie
            if err := rows.Scan(&m.MovieID, &m.Title, &m.Year, &m.Plot, &m.Genre, &m.ImdbID); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            movies = append(movies, m)
        }
        c.JSON(http.StatusOK, movies)
    }
}