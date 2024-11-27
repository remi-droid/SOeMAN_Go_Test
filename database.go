package main

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Document struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"unique;	not null" json:"name"`
	Url        string    `gorm:"unique; not null" json:"url"`
	UploadDate time.Time `gorm:"autoUpdateTime" json:"uploaded_at"`
}

const Dsn = "postgres://upload-service:password@postgres:5432/main"

var Database *gorm.DB

func OpenDB() error {
	var err error

	// Connection to the database
	Database, err = gorm.Open(postgres.Open(Dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("an error occured during database connection: %w", err)
	}

	err = Database.AutoMigrate(&Document{})
	if err != nil {
		return fmt.Errorf("erreur lors de la migration: %w", err)
	}

	return nil
}

func ClearDatabase() {
	// Delete all the documents in the database
	Database.Unscoped().Delete(&Document{})
}
