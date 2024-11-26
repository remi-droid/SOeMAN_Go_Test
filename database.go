package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const Dsn = "postgres://upload-service:password@postgres:5432/main"

var db *gorm.DB

func InitDB() error {
	var err error

	// Connexion à la base de données
	db, err = gorm.Open(postgres.Open(Dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("erreur lors de la connexion à la base de données: %w", err)
	}

	// Migration des tables vers la base de données
	err = db.AutoMigrate(&Document{})
	if err != nil {
		return fmt.Errorf("erreur lors de la migration: %w", err)
	}

	fmt.Println("Initialisation de la base de données réussie.")
	return nil
}
