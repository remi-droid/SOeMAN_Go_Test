package main

import (
	"fmt"
	"os"
	"time"
)

const filePath = "files/"

type Document struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"unique;	not null" json:"name"`
	Url        string    `gorm:"unique; not null" json:"url"`
	UploadDate time.Time `gorm:"autoUpdateTime" json:"uploaded_at"`
}

func main() {

	InitDB()

	filepath := filePath + "test.txt"
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println(string(data))

	// r := gin.Default()

	// // Routes

	// // The servers runs on port 8080
	// r.Run(":8080")

}
