package main

import (
	"log"

	"github.com/eciccone/rh/database"
	"github.com/eciccone/rh/router"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if godotenv.Load() != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.Open()
	if err != nil {
		log.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

	r := router.New()
	r.BuildRoutes(db)
	r.Run(":8080")
}
