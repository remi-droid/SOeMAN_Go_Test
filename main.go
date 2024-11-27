package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	InitDB()

	// Initialisation du serveur
	router := gin.Default()

	// Routes

	router.GET("/dl/:filename", DownloadDocumentHandler)
	router.POST("/ul", UploadDocumentHandler)

	// The servers runs on port 80
	router.Run(":80")

}
