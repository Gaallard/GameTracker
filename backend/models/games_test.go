package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGame_Struct(t *testing.T) {
	t.Run("Game struct should have all required fields", func(t *testing.T) {
		now := time.Now()
		game := Game{
			ID:           1,
			Title:        "Test Game",
			Platform:     "PC",
			Genre:        "Action",
			Status:       "In Progress",
			Progress:     50,
			HoursPlayed:  10.5,
			PersonalNote: "Great game!",
			Score:        8,
			StartedAt:    &now,
			FinishedAt:   nil,
			CoverURL:     "https://example.com/cover.jpg",
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		assert.Equal(t, uint(1), game.ID)
		assert.Equal(t, "Test Game", game.Title)
		assert.Equal(t, "PC", game.Platform)
		assert.Equal(t, "Action", game.Genre)
		assert.Equal(t, "In Progress", game.Status)
		assert.Equal(t, 50, game.Progress)
		assert.Equal(t, 10.5, game.HoursPlayed)
		assert.Equal(t, "Great game!", game.PersonalNote)
		assert.Equal(t, 8, game.Score)
		assert.NotNil(t, game.StartedAt)
		assert.Nil(t, game.FinishedAt)
		assert.Equal(t, "https://example.com/cover.jpg", game.CoverURL)
		assert.Equal(t, now, game.CreatedAt)
		assert.Equal(t, now, game.UpdatedAt)
	})

	t.Run("Game with minimal fields", func(t *testing.T) {
		game := Game{
			Title:    "Minimal Game",
			Platform: "Switch",
		}

		assert.Equal(t, "Minimal Game", game.Title)
		assert.Equal(t, "Switch", game.Platform)
		assert.Equal(t, 0, game.Progress)
		assert.Equal(t, 0.0, game.HoursPlayed)
		assert.Equal(t, 0, game.Score)
	})
}

func TestGameStats_Struct(t *testing.T) {
	t.Run("GameStats with all fields populated", func(t *testing.T) {
		stats := GameStats{
			TotalGames:      5,
			ByStatus:        map[string]int{"Completed": 2, "In Progress": 3},
			AverageHours:    25.5,
			MostPlayedGenre: "RPG",
			PendingGames:    3,
		}

		assert.Equal(t, 5, stats.TotalGames)
		assert.Equal(t, 2, stats.ByStatus["Completed"])
		assert.Equal(t, 3, stats.ByStatus["In Progress"])
		assert.Equal(t, 25.5, stats.AverageHours)
		assert.Equal(t, "RPG", stats.MostPlayedGenre)
		assert.Equal(t, 3, stats.PendingGames)
	})

	t.Run("GameStats with empty status map", func(t *testing.T) {
		stats := GameStats{
			TotalGames:      0,
			ByStatus:        make(map[string]int),
			AverageHours:    0,
			MostPlayedGenre: "",
			PendingGames:    0,
		}

		assert.Equal(t, 0, stats.TotalGames)
		assert.Empty(t, stats.ByStatus)
		assert.Equal(t, 0.0, stats.AverageHours)
		assert.Empty(t, stats.MostPlayedGenre)
		assert.Equal(t, 0, stats.PendingGames)
	})
}

func TestGame_Validation(t *testing.T) {
	t.Run("Game with valid progress range", func(t *testing.T) {
		game := Game{
			Title:    "Test",
			Platform: "PC",
			Progress: 0,
		}
		assert.Equal(t, 0, game.Progress)

		game.Progress = 100
		assert.Equal(t, 100, game.Progress)

		game.Progress = 50
		assert.Equal(t, 50, game.Progress)
	})

	t.Run("Game with valid score range", func(t *testing.T) {
		game := Game{
			Title:    "Test",
			Platform: "PC",
			Score:    0,
		}
		assert.Equal(t, 0, game.Score)

		game.Score = 10
		assert.Equal(t, 10, game.Score)

		game.Score = 5
		assert.Equal(t, 5, game.Score)
	})
}
