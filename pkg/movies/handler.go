package movies

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func ListMoviesHandler(repo MovieRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        genre := c.Query("genre")
        actor := c.Query("actor")
        year := c.Query("year")

        movies, err := repo.ListMovies(c.Request.Context(), genre, actor, year)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, movies)
    }
}

func GetMovieByIDHandler(repo MovieRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")

        movie, err := repo.GetMovieByID(c.Request.Context(), id)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if movie == nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
            return
        }

        c.JSON(http.StatusOK, movie)
    }
}