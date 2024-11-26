package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Document struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Url        string    `json:"url"`
	UploadDate time.Time `json:"uploaded_at"`
}

// The DSN to connect to postgres is: "postgres://upload-service:password@postgres:5432/main".

func main() {

	r := gin.Default()

	// Routes

	// The servers runs on port 8080
	r.Run(":8080")

}
