package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

const DownloadRoute = "/dl/"

func main() {

	err := OpenDB()
	if err != nil {
		log.Fatal("Program exit : An error occured during database initialization -> ", err)
	}

	// Initialisation du serveur
	router := gin.Default()

	// Routes
	router.GET("/list", ListDocumentsHandler)
	router.GET(DownloadRoute+":filename", DownloadDocumentHandler)
	router.POST("/ul", UploadDocumentHandler)

	// The servers runs on port 80
	router.Run(":80")
}
