package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const Dsn = "postgres://upload-service:password@postgres:5432/main"

var Database *gorm.DB

func InitDB() error {
	var err error

	// Connection to the database
	Database, err = gorm.Open(postgres.Open(Dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("an error occured during database connection: %w", err)
	}

	// Migration of all the tables to the database
	err = Database.AutoMigrate(&Document{})
	if err != nil {
		return fmt.Errorf("erreur lors de la migration: %w", err)
	}

	fmt.Println("Initialisation de la base de données réussie.")
	return nil
}
