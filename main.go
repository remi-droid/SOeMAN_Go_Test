package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	InitDB()

	r := gin.Default()

	// Routes

	// The servers runs on port 8080
	r.Run(":8080")

}
