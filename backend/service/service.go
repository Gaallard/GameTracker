package service

import (
	"errors"
	"gametracker/db"
	"gametracker/models"
)

func GetAllGames() ([]models.Game, error) {
	var games []models.Game
	result := db.DB.Find(&games)
	return games, result.Error
}

func GetGameByID(id string) (models.Game, error) {
	var game models.Game
	result := db.DB.First(&game, id)
	if result.Error != nil {
		return game, errors.New("game not found")
	}
	return game, nil
}

func CreateGame(game *models.Game) error {
	return db.DB.Create(game).Error
}

func UpdateGame(game *models.Game) error {
	return db.DB.Save(game).Error
}

func DeleteGame(id string) error {
	return db.DB.Delete(&models.Game{}, id).Error
}

func GetByTitle(title string) ([]models.Game, error) {
	var games []models.Game
	query := "%" + title + "%"
	result := db.DB.Where("title LIKE ?", query).Find(&games)
	return games, result.Error
}

func GetByStatus(status string) ([]models.Game, error) {
	var games []models.Game
	query := "%" + status + "%"
	result := db.DB.Where("status LIKE ?", query).Find(&games)
	return games, result.Error
}

func GetByGenre(genre string) ([]models.Game, error) {
	var games []models.Game
	query := "%" + genre + "%"
	result := db.DB.Where("genre LIKE ?", query).Find(&games)
	return games, result.Error
}

func GetStats() (models.GameStats, error) {
	var stats models.GameStats
	var games []models.Game

	result := db.DB.Find(&games)
	if result.Error != nil {
		return stats, result.Error
	}

	statusCount := make(map[string]int)
	genreCount := make(map[string]int)
	var pendingCount int
	var totalHours float64
	for _, game := range games {
		statusCount[game.Status]++
		genreCount[game.Genre]++
		totalHours += game.HoursPlayed

		if game.Status != "Completed" && game.Progress < 100 {
			pendingCount++
		}
	}

	mostPlayedGenre := ""
	maxGenreCount := 0
	for genre, count := range genreCount {
		if count > maxGenreCount {
			maxGenreCount = count
			mostPlayedGenre = genre
		}
	}

	totalGames := len(games)
	averageHours := 0.0
	if totalGames > 0 {
		averageHours = totalHours / float64(totalGames)
	}

	stats = models.GameStats{
		TotalGames:      totalGames,
		ByStatus:        statusCount,
		AverageHours:    averageHours,
		MostPlayedGenre: mostPlayedGenre,
		PendingGames:    pendingCount,
	}

	return stats, nil
}
