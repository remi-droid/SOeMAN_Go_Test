package main

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Document represents a document entry is the postgre database
type Document struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"unique;	not null" json:"name"`
	Url        string    `gorm:"unique; not null" json:"url"`
	UploadDate time.Time `gorm:"autoUpdateTime" json:"uploaded_at"`
}

// Dsn is the connection string used to connect to the postgre database
const Dsn = "postgres://upload-service:password@postgres:5432/main"

// Database is the variable used in the program to represent the database
var Database *gorm.DB

// OpenDB creates a connection to the database from the connection string above
func OpenDB() error {
	var err error

	// Connection to the database
	Database, err = gorm.Open(postgres.Open(Dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("an error occured during database connection: %w", err)
	}

	// Migration of the Document model to the database
	err = Database.AutoMigrate(&Document{})
	if err != nil {
		return fmt.Errorf("an error occured during the migration: %w", err)
	}

	return nil
}

// ClearDatabase delete all the documents in the postgre database
func ClearDatabase() {
	Database.Where("1 = 1").Delete(&Document{})
}

// DocumentIsPresent checks if a document exists in the database from a filename passed in argument
func DocumentIsPresent(filename string) (bool, error) {

	result := Database.First(&Document{}, "name = ?", filename)

	// If the requests encountered an error we check if the error was the absence of document or something else
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil

}

// InsertInDatabase takes a filename as parameter and creates a document entry in the database
func InsertInDatabase(filename string) error {
	result := Database.Create(&Document{Name: filename, Url: "http://localhost" + DownloadRoute + filename})
	return result.Error
}
