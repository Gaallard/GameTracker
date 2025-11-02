package controller

import (
	"bytes"
	"encoding/json"
	"gametracker/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewAuthController(t *testing.T) {
	controller := NewAuthController()
	assert.NotNil(t, controller)
}

func TestAuthController_Register_Success(t *testing.T) {
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	req := models.RegisterRequest{
		Username:  "newuser",
		Email:     "newuser@example.com",
		Password:  "password123",
		FirstName: "New",
		LastName:  "User",
	}

	mock.ExpectQuery("SELECT .* FROM `users` WHERE username = .* OR email = .* ORDER BY `users`.`id` LIMIT .").
		WillReturnError(gorm.ErrRecordNotFound)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	authController := NewAuthController()
	router.POST("/register", authController.Register)

	w := httptest.NewRecorder()
	jsonData, _ := json.Marshal(req)
	req_http, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req_http.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req_http)

	assert.Equal(t, http.StatusCreated, w.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthController_Register_InvalidData(t *testing.T) {
	_, _, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	authController := NewAuthController()
	router.POST("/register", authController.Register)

	invalidJSON := `{"username": "test"`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthController_Login_Success(t *testing.T) {
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	// Hash password
	testUser := models.User{}
	err := testUser.HashPassword("password123")
	require.NoError(t, err)
	hashedPassword := testUser.Password

	req := models.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "first_name", "last_name", "created_at", "updated_at"}).
		AddRow(1, "testuser", "test@example.com", hashedPassword, "Test", "User", time.Now(), time.Now())

	mock.ExpectQuery("SELECT .* FROM `users` WHERE username = .* OR email = .* ORDER BY `users`.`id` LIMIT .").
		WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	authController := NewAuthController()
	router.POST("/login", authController.Login)

	w := httptest.NewRecorder()
	jsonData, _ := json.Marshal(req)
	req_http, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req_http.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req_http)

	assert.Equal(t, http.StatusOK, w.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthController_Login_InvalidData(t *testing.T) {
	_, _, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	authController := NewAuthController()
	router.POST("/login", authController.Login)

	invalidJSON := `{"username": }`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthController_GetProfile_Authenticated(t *testing.T) {
	_, _, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	authController := NewAuthController()
	router.GET("/profile", func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	}, authController.GetProfile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthController_GetProfile_Unauthenticated(t *testing.T) {
	_, _, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	authController := NewAuthController()
	router.GET("/profile", authController.GetProfile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

