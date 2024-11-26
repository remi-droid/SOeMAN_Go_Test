package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const filePath = "files/"

type Document struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Url        string    `json:"url"`
	UploadDate time.Time `json:"uploaded_at"`
}

// The DSN to connect to postgres is: "postgres://upload-service:password@postgres:5432/main".

func main() {

	filepath := filePath + "test.txt"
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println(string(data))

	r := gin.Default()

	// Routes

	// The servers runs on port 8080
	r.Run(":8080")

}
