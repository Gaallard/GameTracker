package main

import (
	"gametracker/db"
	"gametracker/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	// Set Gin mode based on environment
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug" // default
	}
	gin.SetMode(ginMode)

	// Set log level
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info" // default
	}
	log.Printf("Starting GameTracker in %s mode with log level: %s", ginMode, logLevel)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://localhost:80", "http://localhost:3000", "http://localhost:8080", "http://localhost:5173"}, // URLs del frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	db.ConnectDB()
	routes.SetupGameRoutes(r)

	// Get port from environment
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080" // default
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}