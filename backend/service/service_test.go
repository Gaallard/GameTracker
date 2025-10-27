package service

import (
	"database/sql"
	"gametracker/db"
	"gametracker/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// setupTestDB crea una base de datos mock para testing
func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	// Configurar la variable global db.DB para que las funciones de servicio la usen
	db.DB = gormDB

	return gormDB, mock, sqlDB
}

func TestGetAllGames(t *testing.T) {
	// Arrange
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	// Mock de datos de prueba
	games := []models.Game{
		{
			ID:           1,
			Title:        "Test Game 1",
			Platform:     "PC",
			Genre:        "RPG",
			Status:       "Completed",
			Progress:     100,
			HoursPlayed:  25.5,
			PersonalNote: "Great game",
			Score:        8,
			CoverURL:     "http://example.com/cover1.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           2,
			Title:        "Test Game 2",
			Platform:     "PS5",
			Genre:        "Action",
			Status:       "Playing",
			Progress:     50,
			HoursPlayed:  15.0,
			PersonalNote: "Fun game",
			Score:        7,
			CoverURL:     "http://example.com/cover2.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	// Configurar mock para retornar los juegos
	rows := sqlmock.NewRows([]string{"id", "title", "platform", "genre", "status", "progress", "hours_played", "personal_note", "score", "started_at", "finished_at", "cover_url", "created_at", "updated_at"}).
		AddRow(games[0].ID, games[0].Title, games[0].Platform, games[0].Genre, games[0].Status, games[0].Progress, games[0].HoursPlayed, games[0].PersonalNote, games[0].Score, games[0].StartedAt, games[0].FinishedAt, games[0].CoverURL, games[0].CreatedAt, games[0].UpdatedAt).
		AddRow(games[1].ID, games[1].Title, games[1].Platform, games[1].Genre, games[1].Status, games[1].Progress, games[1].HoursPlayed, games[1].PersonalNote, games[1].Score, games[1].StartedAt, games[1].FinishedAt, games[1].CoverURL, games[1].CreatedAt, games[1].UpdatedAt)

	mock.ExpectQuery("SELECT \\* FROM `games`").
		WillReturnRows(rows)

	// Act
	result, err := GetAllGames()

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, games[0].Title, result[0].Title)
	assert.Equal(t, games[1].Title, result[1].Title)
}

func TestGetGameByID_Success(t *testing.T) {
	// Arrange
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	game := models.Game{
		ID:           1,
		Title:        "Test Game",
		Platform:     "PC",
		Genre:        "RPG",
		Status:       "Completed",
		Progress:     100,
		HoursPlayed:  25.5,
		PersonalNote: "Great game",
		Score:        8,
		CoverURL:     "http://example.com/cover.jpg",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "title", "platform", "genre", "status", "progress", "hours_played", "personal_note", "score", "started_at", "finished_at", "cover_url", "created_at", "updated_at"}).
		AddRow(game.ID, game.Title, game.Platform, game.Genre, game.Status, game.Progress, game.HoursPlayed, game.PersonalNote, game.Score, game.StartedAt, game.FinishedAt, game.CoverURL, game.CreatedAt, game.UpdatedAt)

	mock.ExpectQuery("SELECT \\* FROM `games` WHERE `games`.`id` = \\? ORDER BY `games`.`id` LIMIT \\?").
		WithArgs("1", 1).
		WillReturnRows(rows)

	// Act
	result, err := GetGameByID("1")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, game.Title, result.Title)
	assert.Equal(t, game.Platform, result.Platform)
}

func TestGetGameByID_NotFound(t *testing.T) {
	// Arrange
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	mock.ExpectQuery("SELECT \\* FROM `games` WHERE `games`.`id` = \\? ORDER BY `games`.`id` LIMIT \\?").
		WithArgs("999", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Act
	result, err := GetGameByID("999")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "game not found", err.Error())
	assert.Empty(t, result.Title)
}

func TestCreateGame_Success(t *testing.T) {
	// Arrange
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	game := &models.Game{
		Title:        "New Game",
		Platform:     "PC",
		Genre:        "RPG",
		Status:       "Not Started",
		Progress:     0,
		HoursPlayed:  0,
		PersonalNote: "New game to play",
		Score:        0,
		CoverURL:     "http://example.com/cover.jpg",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `games`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act
	err := CreateGame(game)

	// Assert
	require.NoError(t, err)
}

func TestGetStats_Success(t *testing.T) {
	// Arrange
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	// Mock de datos para estadísticas
	games := []models.Game{
		{
			ID:          1,
			Title:       "Game 1",
			Platform:    "PC",
			Genre:       "RPG",
			Status:      "Completed",
			Progress:    100,
			HoursPlayed: 25.5,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Title:       "Game 2",
			Platform:    "PS5",
			Genre:       "RPG",
			Status:      "Playing",
			Progress:    50,
			HoursPlayed: 15.0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          3,
			Title:       "Game 3",
			Platform:    "Xbox",
			Genre:       "Action",
			Status:      "Not Started",
			Progress:    0,
			HoursPlayed: 0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "platform", "genre", "status", "progress", "hours_played", "personal_note", "score", "started_at", "finished_at", "cover_url", "created_at", "updated_at"})
	for _, game := range games {
		rows.AddRow(game.ID, game.Title, game.Platform, game.Genre, game.Status, game.Progress, game.HoursPlayed, game.PersonalNote, game.Score, game.StartedAt, game.FinishedAt, game.CoverURL, game.CreatedAt, game.UpdatedAt)
	}

	mock.ExpectQuery("SELECT \\* FROM `games`").
		WillReturnRows(rows)

	// Act
	stats, err := GetStats()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 3, stats.TotalGames)
	assert.Equal(t, 2, stats.PendingGames) // Game 2 (Playing, Progress 50) y Game 3 (Not Started, Progress 0) están pendientes
	assert.Equal(t, "RPG", stats.MostPlayedGenre)
	assert.Equal(t, 13.5, stats.AverageHours) // (25.5 + 15.0 + 0) / 3
	assert.Equal(t, 1, stats.ByStatus["Completed"])
	assert.Equal(t, 1, stats.ByStatus["Playing"])
	assert.Equal(t, 1, stats.ByStatus["Not Started"])
}

func TestGetByTitle_Success(t *testing.T) {
	// Arrange
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	game := models.Game{
		ID:           1,
		Title:        "Test Game",
		Platform:     "PC",
		Genre:        "RPG",
		Status:       "Completed",
		Progress:     100,
		HoursPlayed:  25.5,
		PersonalNote: "Great game",
		Score:        8,
		CoverURL:     "http://example.com/cover.jpg",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "title", "platform", "genre", "status", "progress", "hours_played", "personal_note", "score", "started_at", "finished_at", "cover_url", "created_at", "updated_at"}).
		AddRow(game.ID, game.Title, game.Platform, game.Genre, game.Status, game.Progress, game.HoursPlayed, game.PersonalNote, game.Score, game.StartedAt, game.FinishedAt, game.CoverURL, game.CreatedAt, game.UpdatedAt)

	mock.ExpectQuery("SELECT \\* FROM `games` WHERE title LIKE \\?").
		WithArgs("%Test%").
		WillReturnRows(rows)

	// Act
	result, err := GetByTitle("Test")

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, game.Title, result[0].Title)
}

func TestGetByStatus_Success(t *testing.T) {
	// Arrange
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	game := models.Game{
		ID:           1,
		Title:        "Test Game",
		Platform:     "PC",
		Genre:        "RPG",
		Status:       "Completed",
		Progress:     100,
		HoursPlayed:  25.5,
		PersonalNote: "Great game",
		Score:        8,
		CoverURL:     "http://example.com/cover.jpg",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "title", "platform", "genre", "status", "progress", "hours_played", "personal_note", "score", "started_at", "finished_at", "cover_url", "created_at", "updated_at"}).
		AddRow(game.ID, game.Title, game.Platform, game.Genre, game.Status, game.Progress, game.HoursPlayed, game.PersonalNote, game.Score, game.StartedAt, game.FinishedAt, game.CoverURL, game.CreatedAt, game.UpdatedAt)

	mock.ExpectQuery("SELECT \\* FROM `games` WHERE status LIKE \\?").
		WithArgs("%Completed%").
		WillReturnRows(rows)

	// Act
	result, err := GetByStatus("Completed")

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, game.Status, result[0].Status)
}

func TestGetByGenre_Success(t *testing.T) {
	// Arrange
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	game := models.Game{
		ID:           1,
		Title:        "Test Game",
		Platform:     "PC",
		Genre:        "RPG",
		Status:       "Completed",
		Progress:     100,
		HoursPlayed:  25.5,
		PersonalNote: "Great game",
		Score:        8,
		CoverURL:     "http://example.com/cover.jpg",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "title", "platform", "genre", "status", "progress", "hours_played", "personal_note", "score", "started_at", "finished_at", "cover_url", "created_at", "updated_at"}).
		AddRow(game.ID, game.Title, game.Platform, game.Genre, game.Status, game.Progress, game.HoursPlayed, game.PersonalNote, game.Score, game.StartedAt, game.FinishedAt, game.CoverURL, game.CreatedAt, game.UpdatedAt)

	mock.ExpectQuery("SELECT \\* FROM `games` WHERE genre LIKE \\?").
		WithArgs("%RPG%").
		WillReturnRows(rows)

	// Act
	result, err := GetByGenre("RPG")

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, game.Genre, result[0].Genre)
}

func TestUpdateGame_Success(t *testing.T) {
	// Arrange
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	game := &models.Game{
		ID:           1,
		Title:        "Updated Game",
		Platform:     "PC",
		Genre:        "RPG",
		Status:       "Completed",
		Progress:     100,
		HoursPlayed:  30.0,
		PersonalNote: "Updated note",
		Score:        9,
		CoverURL:     "http://example.com/cover.jpg",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `games` SET").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act
	err := UpdateGame(game)

	// Assert
	require.NoError(t, err)
}

func TestDeleteGame_Success(t *testing.T) {
	// Arrange
	_, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `games` WHERE `games`.`id` = \\?").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act
	err := DeleteGame("1")

	// Assert
	require.NoError(t, err)
}
