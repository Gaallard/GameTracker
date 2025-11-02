package service

import (
	"gametracker/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewAuthService(t *testing.T) {
	service := NewAuthService()
	assert.NotNil(t, service)
}

func TestAuthService_Register_Success(t *testing.T) {
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	req := models.RegisterRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	mock.ExpectQuery("SELECT .* FROM `users` WHERE username = .* OR email = .* ORDER BY `users`.`id` LIMIT .").
		WillReturnError(gorm.ErrRecordNotFound)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	service := NewAuthService()
	user, err := service.Register(req)

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Empty(t, user.Password)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthService_Register_UserExists(t *testing.T) {
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	req := models.RegisterRequest{
		Username: "existinguser",
		Email:    "existing@example.com",
		Password: "password123",
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email"}).
		AddRow(1, "existinguser", "existing@example.com")

	mock.ExpectQuery("SELECT .* FROM `users` WHERE username = .* OR email = .* ORDER BY `users`.`id` LIMIT .").
		WillReturnRows(rows)

	service := NewAuthService()
	user, err := service.Register(req)

	require.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "ya existe")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	req := models.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	mock.ExpectQuery("SELECT .* FROM `users` WHERE username = .* OR email = .* ORDER BY `users`.`id` LIMIT .").
		WillReturnError(gorm.ErrRecordNotFound)

	service := NewAuthService()
	authResponse, err := service.Login(req)

	require.Error(t, err)
	assert.Nil(t, authResponse)
	assert.Contains(t, err.Error(), "no encontrado")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	_, _, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	service := NewAuthService()
	token, err := service.ValidateToken("invalid.token.here")

	// ValidateToken returns a token even if invalid, so err can be nil or not
	// The token.Valid will be false
	if err == nil {
		assert.NotNil(t, token)
		assert.False(t, token.Valid)
	} else {
		assert.Error(t, err)
	}
}

func TestAuthService_GetUserFromToken_InvalidToken(t *testing.T) {
	_, _, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	service := NewAuthService()
	
	// Create an invalid token
	token, _ := service.ValidateToken("invalid.token.here")
	userID, username, err := service.GetUserFromToken(token)

	require.Error(t, err)
	assert.Equal(t, uint(0), userID)
	assert.Empty(t, username)
}

func TestAuthService_Login_Success(t *testing.T) {
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	service := NewAuthService()

	// Hash a real password for testing
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

	authResponse, err := service.Login(req)

	require.NoError(t, err)
	assert.NotNil(t, authResponse)
	assert.NotEmpty(t, authResponse.Token)
	assert.Equal(t, "testuser", authResponse.User.Username)
	assert.Empty(t, authResponse.User.Password)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	service := NewAuthService()

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

	// Login to get valid token
	authResponse, err := service.Login(req)
	require.NoError(t, err)

	// Validate the token
	token, err := service.ValidateToken(authResponse.Token)

	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.True(t, token.Valid)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthService_GetUserFromToken_Success(t *testing.T) {
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	service := NewAuthService()

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

	// Login to get valid token
	authResponse, err := service.Login(req)
	require.NoError(t, err)

	// Validate token
	token, err := service.ValidateToken(authResponse.Token)
	require.NoError(t, err)

	// Extract user info from token
	userID, username, err := service.GetUserFromToken(token)

	require.NoError(t, err)
	assert.Equal(t, uint(1), userID)
	assert.Equal(t, "testuser", username)

	require.NoError(t, mock.ExpectationsWereMet())
}
