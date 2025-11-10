package main

import (
	"gametracker/db"
	"gametracker/routes"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// Configure CORS
	allowedOrigins := []string{
		"http://localhost",
		"http://localhost:80",
		"http://localhost:3000",
		"http://localhost:8080",
		"http://localhost:5173",
	}

	// Add origins from environment variable (comma-separated)
	if corsOrigins := os.Getenv("CORS_ORIGINS"); corsOrigins != "" {
		origins := strings.Split(corsOrigins, ",")
		for _, origin := range origins {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				allowedOrigins = append(allowedOrigins, origin)
			}
		}
	}

	log.Printf("Allowed CORS origins: %v", allowedOrigins)

	r := gin.Default()

	// Configure CORS with flexible origin handling
	corsConfig := cors.Config{
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
		ExposeHeaders: []string{"Content-Length", "Content-Type"},
		MaxAge:        12 * 3600, // 12 hours
	}

	// If CORS_ORIGINS is set, use specific origins with credentials
	if os.Getenv("CORS_ORIGINS") != "" {
		corsConfig.AllowOrigins = allowedOrigins
		corsConfig.AllowCredentials = true
		log.Printf("CORS: Using specific origins with credentials")
	} else {
		// Allow all origins (for development or when frontend uses proxy)
		// Note: Cannot use AllowCredentials with AllowAllOrigins
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowCredentials = false
		log.Printf("CORS: Allowing all origins (no credentials)")
	}

	r.Use(cors.New(corsConfig))

	db.ConnectDB()
	routes.SetupGameRoutes(r)
	routes.SetupAuthRoutes(r)

	// Get port from environment
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080" // default
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
