package main

import (
	"context"
	"log"
	"os"
	"praktikum4-crud/config"
	"praktikum4-crud/database"

	"github.com/joho/godotenv"
)

// @title Praktikum 4 CRUD API Documentation
// @version 1.0
// @description This is a RESTful API documentation for Praktikum 4 CRUD Project using Go Fiber, PostgreSQL, and MongoDB.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/yourusername
// @contact.email support@example.com

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal load .env")
	}

	database.ConnectDB()
	database.ConnectMongo()

	defer database.DB.Close()
	defer database.MongoClient.Disconnect(context.Background())

	app := config.NewApp()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}