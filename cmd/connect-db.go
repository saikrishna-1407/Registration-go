package main

import (
	"finalreg/config"
	"finalreg/internal/store"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectPostgres(conf *config.Config) (*gorm.DB, error) {

	fmt.Println("Starting PostgreSQL connection process...")

	fmt.Printf("Using DatabaseURI: %s\n", conf.DatabaseURI)

	config := &gorm.Config{}

	fmt.Println("GORM configuration initialized.")

	fmt.Println("Attempting to open connection to PostgreSQL...")

	db, err := gorm.Open(postgres.Open(conf.DatabaseURI), config)
	if err != nil {
		fmt.Printf("Failed to connect to PostgreSQL. Error: %v\n", err)
		return nil, err
	}
	fmt.Println("Successfully connected to PostgreSQL database.")

	fmt.Println("Running database migrations...")
	if err := store.Migrate(db); err != nil {
		fmt.Printf("Database migration failed. Error: %v\n", err)
		return nil, fmt.Errorf("migration failed: %w", err)
	}
	fmt.Println("Database migrations completed successfully.")

	fmt.Println("Returning PostgreSQL database connection.")
	return db, nil

}
