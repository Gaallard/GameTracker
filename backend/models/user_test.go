package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_HashPassword(t *testing.T) {
	t.Run("HashPassword should successfully hash password", func(t *testing.T) {
		user := User{
			Username: "testuser",
			Email:    "test@example.com",
		}

		err := user.HashPassword("password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, user.Password)
		assert.NotEqual(t, "password123", user.Password) // Should be hashed
	})

	t.Run("HashPassword should handle empty password", func(t *testing.T) {
		user := User{
			Username: "testuser",
			Email:    "test@example.com",
		}

		err := user.HashPassword("")
		assert.NoError(t, err)
		assert.NotEmpty(t, user.Password)
	})
}

func TestUser_CheckPassword(t *testing.T) {
	t.Run("CheckPassword should return true for correct password", func(t *testing.T) {
		user := User{
			Username: "testuser",
			Email:    "test@example.com",
		}

		plainPassword := "password123"
		err := user.HashPassword(plainPassword)
		assert.NoError(t, err)

		isValid := user.CheckPassword(plainPassword)
		assert.True(t, isValid)
	})

	t.Run("CheckPassword should return false for incorrect password", func(t *testing.T) {
		user := User{
			Username: "testuser",
			Email:    "test@example.com",
		}

		err := user.HashPassword("password123")
		assert.NoError(t, err)

		isValid := user.CheckPassword("wrongpassword")
		assert.False(t, isValid)
	})

	t.Run("CheckPassword should return false for empty password", func(t *testing.T) {
		user := User{
			Username: "testuser",
			Email:    "test@example.com",
		}

		err := user.HashPassword("password123")
		assert.NoError(t, err)

		isValid := user.CheckPassword("")
		assert.False(t, isValid)
	})
}

func TestUser_Struct(t *testing.T) {
	t.Run("User struct should have all required fields", func(t *testing.T) {
		user := User{
			ID:        1,
			Username:  "testuser",
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
		}

		assert.Equal(t, uint(1), user.ID)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "Test", user.FirstName)
		assert.Equal(t, "User", user.LastName)
	})
}
