package main

import (
	"github.com/gin-gonic/gin"
	"movie-rental/pkg/hello"
)

func main() {
	r := gin.Default()
	r.GET("/hello", hello.HelloHandler)
	r.Run(":8080")
}