package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const uploadedDocumentsPath = "uploaded_documents/"

func UploadDocumentHandler(c *gin.Context) {

	document, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No document found in the 'document' field of the request's body."})
		return
	}

	documentPath := filepath.Join(uploadedDocumentsPath, document.Filename)

	// If a file with this name already exist we dont upload the new one
	if _, err := os.Stat(documentPath); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "A document with this name already exists"})
		return
	} else if !os.IsNotExist(err) {
		// Another error during file analysis occured
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking document existence"})
		return
	}

	if err := c.SaveUploadedFile(document, documentPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the document"})
		return
	}

	if err := InsertInDatabase(document.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error during insertion into the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file": document.Filename})
}

func DownloadDocumentHandler(c *gin.Context) {

	filename := c.Param("filename")
	filePath := filepath.Join(uploadedDocumentsPath, filename)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File '" + filename + "' not found."})
		return
	}

	// Headers to automatically download the file
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(filePath)
}

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

/********************************************************************************/
/*                      	MINIO STORAGE FUNCTIONS								*/
/********************************************************************************/

func UploadDocumentMinioHandler(c *gin.Context) {
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No document found in the 'document' field of the request's body."})
		return
	}

	// Read the file content
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open the document."})
		return
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read the document"})
		return
	}

	if err := InsertInDatabase(file.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error during insertion into the database"})
		return
	}

	// Upload to MinIO
	err = UploadFile(file.Filename, fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload document to storage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document " + file.Filename + "uploaded successfully", "file_name": file.Filename})
}
