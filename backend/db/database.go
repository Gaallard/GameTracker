package db

import (
	"gametracker/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func ConnectDB() {
	// Get database configuration from environment variables
	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "root")
	dbName := getEnv("DB_NAME", "gametracker")

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	
	var err error
	maxRetries := 30
	retryDelay := 10 * time.Second

	for i := 0; i < maxRetries; i++ {
		var tempDB *gorm.DB
		tempDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			DB = tempDB
			log.Println("Conexión a la base de datos establecida exitosamente")
			break
		}
		
		log.Printf("Intento %d/%d: No se pudo conectar a la base de datos: %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			log.Printf("Reintentando en %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos después de", maxRetries, "intentos:", err)
	}

	err = DB.AutoMigrate(&models.Game{})
	if err != nil {
		log.Fatal("Error en migración de modelos:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
