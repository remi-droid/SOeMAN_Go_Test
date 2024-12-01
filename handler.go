package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListDocumentsHandler(c *gin.Context) {
	var documents []Document

	// Get all the documents in the database
	result := Database.Find(&documents)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Database empty, no documents to list."})
	}
	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

func UploadDocumentMinioHandler(c *gin.Context) {
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No document found in the 'document' field of the request's body"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open the document"})
		return
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read the document"})
		return
	}

	fileIsPresent, err := DocumentIsPresent(file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "a problem occured during database request"})
		return
	}

	if fileIsPresent {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "A document named with this name was uploaded previously, please use another name"})
		return
	}

	if err := InsertInDatabase(file.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error during insertion into the database"})
		return
	}

	err = UploadDocument(file.Filename, fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload document to storage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document " + file.Filename + " uploaded successfully", "file_name": file.Filename})
}

func DownloadDocumentMinioHandler(c *gin.Context) {

	filename := c.Param("filename")

	fileIsPresent, err := DocumentIsPresent(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "a problem occured during database request"})
		return
	}
	if !fileIsPresent {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "file " + filename + " is not present in the database"})
		return
	}

	document, err := DownloadDocument(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve file from MinIO", "details": err.Error()})
		return
	}
	defer document.Close()

	// Configure headers of the response to download the document
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "application/octet-stream")

	// Put the file in the response
	if _, err := io.Copy(c.Writer, document); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to stream document", "details": err.Error()})
		return
	}
}

func ClearDataHandler(c *gin.Context) {

	ClearDatabase()
	deletedCount, err := ClearMinioStorage()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear the minio storage :", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database and bucket cleared successfully! " + string(deletedCount) + " documents deleted."})
}
