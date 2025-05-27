package main

import (
	"fmt"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"movie-rental/pkg/hello"
	"movie-rental/pkg/movies"
)

func main() {
	db, err := sql.Open("postgres", "postgres://deepaljain:postgres@localhost:5432/movie?sslmode=disable")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer db.Close()

	r := gin.Default()
	r.GET("/hello", hello.HelloHandler)
	r.GET("/movies", movies.ListMoviesHandler(db))
	r.GET("/movies/filter", movies.FilterMoviesHandler(db))
	r.GET("/movies/:id", movies.GetMovieByIDHandler(db))
	r.Run(":8080")
}