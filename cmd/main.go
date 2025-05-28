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
	db, err := sql.Open("postgres", "postgres://deepaljain:postgres@localhost:5432/movierental?sslmode=disable")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer db.Close()

	movieRepo := movies.NewMovieRepository(db)
	cartRepo := cart.NewRepository(db)

	router := gin.Default()
	router.GET("/hello", hello.HelloHandler)
	router.GET("/movies", movies.ListMoviesHandler(movieRepo))
	router.GET("/movies/:id", movies.GetMovieByIDHandler(movieRepo))
	router.POST("/cart", cart.AddToCartHandler(cartRepo))
	router.GET("/cart/:user_id", cart.ViewCartHandler(cartRepo))
	router.Run(":8080")
}