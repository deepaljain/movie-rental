package movies

import (
	"database/sql"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

func ListMoviesHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        genre := c.Query("genre")
        actor := c.Query("actor")
        year := c.Query("year")

        query := "SELECT * FROM movies WHERE 1=1"
        var args []interface{}
        idx := 1

        if genre != "" {
            query += fmt.Sprintf(" AND genre ILIKE '%%' || $%d || '%%'", idx)
            args = append(args, genre)
            idx++
        }
        if actor != "" {
            query += fmt.Sprintf(" AND actors ILIKE '%%' || $%d || '%%'", idx)
            args = append(args, actor)
            idx++
        }
        if year != "" {
            query += fmt.Sprintf(" AND year = $%d", idx)
            args = append(args, year)
            idx++
        }

        rows, err := db.Query(query, args...)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer rows.Close()

        var movies []Movie
        for rows.Next() {
            var m Movie
            if err := rows.Scan(&m.MovieID, &m.Title, &m.Year, &m.Plot, &m.Genre, &m.ImdbID, &m.Actors); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            movies = append(movies, m)
        }
        c.JSON(http.StatusOK, movies)
    }
}

func GetMovieByIDHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        var m Movie
        err := db.QueryRow(
            "SELECT * FROM movies WHERE movie_id = $1", id,
        ).Scan(&m.MovieID, &m.Title, &m.Year, &m.Plot, &m.Genre, &m.ImdbID, &m.Actors)
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
            return
        } else if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, m)
    }
}