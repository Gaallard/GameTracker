package routes

import (
	"gametracker/controller"
	"github.com/gin-gonic/gin"
)

func SetupGameRoutes(r *gin.Engine) {
	games := r.Group("/games")
	{
		games.GET("/", controller.GetAllGames)
		games.POST("/", controller.CreateGame)
		games.GET("/:id", controller.GetGameByID)
		games.PUT("/:id", controller.UpdateGame)
		games.DELETE("/:id", controller.DeleteGame)
		games.GET("/title", controller.GetByTitle)
		games.GET("/status", controller.GetByStatus)
		games.GET("/genre", controller.GetByGenre)
		games.GET("/stats", controller.GetStats)
	}
}
