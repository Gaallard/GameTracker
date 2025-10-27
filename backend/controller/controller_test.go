package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"gametracker/db"
	"gametracker/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ------------------------------------------------------------
// Setup helpers
// ------------------------------------------------------------

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Inyectar en la variable global usada por el service
	prev := db.DB
	db.DB = gormDB
	t.Cleanup(func() {
		db.DB = prev
		_ = sqlDB.Close()
	})

	return gormDB, mock, sqlDB
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/games", GetAllGames)
	router.GET("/games/:id", GetGameByID)
	router.POST("/games", CreateGame)
	router.PUT("/games/:id", UpdateGame)
	router.DELETE("/games/:id", DeleteGame)
	router.GET("/games/search/title", GetByTitle)
	router.GET("/games/search/status", GetByStatus)
	router.GET("/games/search/genre", GetByGenre)
	router.GET("/games/stats", GetStats)

	return router
}

// ------------------------------------------------------------
// Tests
// ------------------------------------------------------------

func TestGetAllGames_Success(t *testing.T) {
	// Arrange
	_, mock, _ := setupTestDB(t)

	now := time.Now()
	game := models.Game{
		ID:           1,
		Title:        "Test Game 1",
		Platform:     "PC",
		Genre:        "RPG",
		Status:       "Completed",
		Progress:     100,
		HoursPlayed:  25.5,
		PersonalNote: "Great game",
		Score:        8,
		StartedAt:    &now,
		FinishedAt:   &now,
		CoverURL:     "http://example.com/cover1.jpg",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	rows := sqlmock.NewRows([]string{
		"id", "title", "platform", "genre", "status", "progress", "hours_played",
		"personal_note", "score", "started_at", "finished_at", "cover_url",
		"created_at", "updated_at",
	}).AddRow(
		game.ID, game.Title, game.Platform, game.Genre, game.Status, game.Progress,
		game.HoursPlayed, game.PersonalNote, game.Score, game.StartedAt, game.FinishedAt,
		game.CoverURL, game.CreatedAt, game.UpdatedAt,
	)

	mock.ExpectQuery(`SELECT \* FROM \` + "`games`" + ``).WillReturnRows(rows)

	router := setupRouter()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/games", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Game
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Len(t, response, 1)
	assert.Equal(t, game.Title, response[0].Title)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateGame_InvalidJSON(t *testing.T) {
	// Arrange
	_, mock, _ := setupTestDB(t)
	router := setupRouter()
	invalidJSON := `{"title": "Test", "invalid": }`

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/games", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "invalid")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetGameByID_NotFound(t *testing.T) {
	// Arrange
	_, _, _ = setupTestDB(t) // no seteamos expectativas SQL: el handler puede resolver 404 sin query exacta
	router := setupRouter()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/games/999", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Game not found", response["error"])
}

func TestUpdateGame_NotFound(t *testing.T) {
	// Arrange
	_, _, _ = setupTestDB(t) // idem: no imponemos expectativa SQL
	router := setupRouter()

	body := models.Game{
		Title:        "Updated Game",
		Platform:     "PC",
		Genre:        "RPG",
		Status:       "Completed",
		Progress:     100,
		HoursPlayed:  30.0,
		PersonalNote: "Updated note",
		Score:        9,
		CoverURL:     "http://example.com/cover.jpg",
	}
	jsonData, _ := json.Marshal(body)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/games/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Game not found", response["error"])
}
