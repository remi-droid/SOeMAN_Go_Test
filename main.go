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

	err = InitMinioStorage()
	if err != nil {
		log.Fatal("Program exit : An error occured during minio storage initialization -> ", err)
	}

	// Server initialization
	router := gin.Default()

	// Endpoints
	router.GET("/list", ListDocumentsHandler)
	router.GET(DownloadRoute+":filename", DownloadDocumentHandler)
	router.POST("/ul", UploadDocumentMinioHandler)

	// The servers runs on port 80
	router.Run(":80")
}
