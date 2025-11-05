package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func ConnectMongo() {
    uri := os.Getenv("MONGODB_URI")
    if uri == "" {
        uri = "mongodb://localhost:27017"
    }

    dbname := os.Getenv("MONGO_DBNAME")
    if dbname == "" {
        dbname = "alumni_mongo"
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatalf("❌ Gagal koneksi MongoDB: %v", err)
    }

    if err := client.Ping(ctx, nil); err != nil {
        log.Fatalf("❌ MongoDB tidak merespons: %v", err)
    }

    MongoClient = client
    MongoDB = client.Database(dbname)
    fmt.Println("✅ Terhubung ke MongoDB:", dbname)
}