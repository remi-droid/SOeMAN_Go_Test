package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListDocumentsHandler answers the list of all the document entries in the database
func ListDocumentsHandler(c *gin.Context) {
	var documents []Document

	// Get all the documents in the database
	result := Database.Find(&documents)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents"})
		return
	}

	// The database is empty and contains no documents
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Database empty, no documents to list."})
	}
	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

// UploadDocumentMinioHandler is the handler used to upload a document passed in a request to the minio storage
func UploadDocumentMinioHandler(c *gin.Context) {
	file, err := c.FormFile("document")

	// No document is attached to the request
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No document found in the 'document' field of the request's body"})
		return
	}

	// Error during the opening of the document
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open the document"})
		return
	}
	defer src.Close()

	// Error during the reading of the document
	fileData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read the document"})
		return
	}

	// Error during the presence check in the database
	fileIsPresent, err := DocumentIsPresent(file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "a problem occured during database request."})
		return
	}

	// We dont upload the document if a document with this name is already present
	if fileIsPresent {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "A document with this name was uploaded previously, please use another name."})
		return
	}

	// Error during the insertion in the postgre database
	if err := InsertInDatabase(file.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occured during insertion into the database."})
		return
	}

	// Error during the upload in the minio storage
	err = UploadDocument(file.Filename, fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload document to storage."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document " + file.Filename + " uploaded successfully", "file_name": file.Filename})
}

// DownloadDocumentMinioHandler looks for a document from its filename in the minio storage and put it in the answer to be downloaded
func DownloadDocumentMinioHandler(c *gin.Context) {

	filename := c.Param("filename")

	// Check if the document is present in the database
	fileIsPresent, err := DocumentIsPresent(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "a problem occured during database request"})
		return
	}

	// If the file is not present in the database we send an error message
	if !fileIsPresent {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "file " + filename + " is not present in the database"})
		return
	}

	// We try to download the document from the minio storage
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

// ClearDataHandler erases all the data stored by the program in the postgre database and the minio storage.
func ClearDataHandler(c *gin.Context) {

	// First we clear the database of all the documents entries, then we clear the minio storage.
	ClearDatabase()
	deletedCount, err := ClearMinioStorage()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear the minio storage :", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Database and bucket cleared successfully! " + string(deletedCount) + " documents deleted."})
}
