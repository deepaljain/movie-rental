package main

import (
	"fmt"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"movie-rental/pkg/hello"
	"movie-rental/pkg/movies"
	"movie-rental/pkg/cart"
)

func main() {
	db, err := sql.Open("postgres", "postgres://deepaljain:postgres@localhost:5432/movie?sslmode=disable")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer db.Close()

	router := gin.Default()
	router.GET("/hello", hello.HelloHandler)
	router.GET("/movies", movies.ListMoviesHandler(db))
	router.GET("/movies/filter", movies.FilterMoviesHandler(db))
	router.GET("/movies/:id", movies.GetMovieByIDHandler(db))
	router.POST("/cart", cart.AddToCartHandler(db))
	router.Run(":8080")
}