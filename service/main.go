package main

import (
	"log/slog"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/gin-gonic/gin"

	"github.com/GRO4T/edgy/videos"
)

func getMongoClient() (*mongo.Client, error) {
	return mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
}

func main() {
	r := gin.Default()

	mongoClient, err := getMongoClient()
	if err != nil {
		slog.Error("failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}

	videos.RegisterRoutes("/api/v1", r, mongoClient)

	if err := r.Run(); err != nil {
		slog.Error("failed to run server", "error", err)
		os.Exit(1)
	}
}
