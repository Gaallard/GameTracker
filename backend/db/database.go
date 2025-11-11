package db

import (
	"gametracker/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "root")
	dbName := getEnv("DB_NAME", "gametracker")

	baseDSN := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" +
		"?charset=utf8mb4&parseTime=True&loc=Local"
	fullDSN := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName +
		"?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	maxRetries := 30
	retryDelay := 10 * time.Second

	var serverDB *gorm.DB
	for i := 0; i < maxRetries; i++ {
		serverDB, err = gorm.Open(mysql.Open(baseDSN), &gorm.Config{})
		if err == nil {
			log.Println("Conexión al servidor MySQL establecida (sin DB seleccionada)")
			break
		}
		log.Printf("Intento %d/%d: No se pudo conectar al servidor MySQL: %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			log.Printf("Reintentando en %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}
	if err != nil {
		log.Fatal("No se pudo conectar al servidor MySQL después de ", maxRetries, " intentos: ", err)
	}

	createDBStmt := "CREATE DATABASE IF NOT EXISTS `" + dbName + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"
	if err := serverDB.Exec(createDBStmt).Error; err != nil {
		log.Fatal("Error creando la base de datos: ", err)
	}
	if sqlDB, cerr := serverDB.DB(); cerr == nil {
		_ = sqlDB.Close()
	}

	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(mysql.Open(fullDSN), &gorm.Config{})
		if err == nil {
			log.Println("Conexión a la base de datos establecida exitosamente: ", dbName)
			break
		}
		log.Printf("Intento %d/%d: No se pudo conectar a la DB %q: %v", i+1, maxRetries, dbName, err)
		if i < maxRetries-1 {
			log.Printf("Reintentando en %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}
	if err != nil {
		log.Fatal("No se pudo conectar a la DB después de ", maxRetries, " intentos: ", err)
	}

	if sqlDB, cerr := DB.DB(); cerr == nil {
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(25)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
	}

	if err := DB.AutoMigrate(&models.Game{}, &models.User{}); err != nil {
		log.Fatal("Error en migración de modelos: ", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

//hola mundo
