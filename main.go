package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

// DownloadRoute represents the endpoint to call to download a file
const DownloadRoute = "/dl/"

func main() {

	// Connection to the postgre database
	err := OpenDB()
	if err != nil {
		log.Fatal("Program exit : An error occured during database initialization -> ", err)
	}

	// Connection to the minio storage
	err = InitMinioStorage()
	if err != nil {
		log.Fatal("Program exit : An error occured during minio storage initialization -> ", err)
	}

	// Server initialization
	router := gin.Default()

	// Endpoints
	router.GET("/list", ListDocumentsHandler)
	router.GET(DownloadRoute+":filename", DownloadDocumentMinioHandler)
	router.POST("/ul", UploadDocumentMinioHandler)
	router.DELETE("/clear", ClearDataHandler)

	// The servers runs on port 80
	router.Run(":80")
}
