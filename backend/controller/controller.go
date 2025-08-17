package controller

import (
	"gametracker/models"
	"gametracker/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllGames(c *gin.Context) {
	games, err := service.GetAllGames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obtaining games"})
		return
	}
	c.JSON(http.StatusOK, games)
}

func GetGameByID(c *gin.Context) {
	id := c.Param("id")
	game, err := service.GetGameByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	c.JSON(http.StatusOK, game)
}

func CreateGame(c *gin.Context) {
	var game models.Game
	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.CreateGame(&game); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating game"})
		return
	}
	c.JSON(http.StatusOK, game)
}

func UpdateGame(c *gin.Context) {
	id := c.Param("id")
	game, err := service.GetGameByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game not found"})
		return
	}
	if err := service.UpdateGame(&game); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating game"})
		return
	}
	c.JSON(http.StatusOK, game)
}

func DeleteGame(c *gin.Context) {
	id := c.Param("id")
	if err := service.DeleteGame(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting game"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Game deleted successfully"})
}

func GetByTitle(c *gin.Context) {
	title := c.Query("title") //esto obtiene el query param ?title=...

	games, err := service.GetByTitle(title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching games"})
		return
	}
	c.JSON(http.StatusOK, games)
}

func GetByStatus(c *gin.Context) {
	status := c.Query("status") //esto obtiene el query param ?status=...

	games, err := service.GetByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching games"})
		return
	}
	c.JSON(http.StatusOK, games)
}

func GetByGenre(c *gin.Context) {
	genre := c.Query("genre") //esto obtiene el query param ?genre=...

	games, err := service.GetByGenre(genre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching games"})
		return
	}
	c.JSON(http.StatusOK, games)
}

func GetStats(c *gin.Context) {
	stats, err := service.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener estadísticas"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
