package main

import (
	"log"
	"os"
	"praktikum4-crud/config"
	"praktikum4-crud/database"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal load .env")
	}

	database.ConnectDB()
	defer database.DB.Close()

	app := config.NewApp()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
