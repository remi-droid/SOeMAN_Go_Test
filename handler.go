package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const uploadedFilesPath = "uploaded_files/"

func UploadDocumentHandler(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found in the 'file' field of the request's body."})
		return
	}

	filePath := filepath.Join(uploadedFilesPath, file.Filename)

	// If a file with this name already exist we dont upload the new one
	if _, err := os.Stat(filePath); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "File with this name already exists"})
		return
	} else if !os.IsNotExist(err) {
		// Another error during file analysis occured
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking file existence"})
		return
	}

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file": file.Filename})
}

func DownloadDocumentHandler(c *gin.Context) {

	filename := c.Param("filename")
	filePath := filepath.Join(uploadedFilesPath, filename)

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
