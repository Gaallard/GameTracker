package routes

import (
	"gametracker/controller"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine) {
	authController := controller.NewAuthController()
	
	// Rutas públicas de autenticación
	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// Rutas protegidas
	protected := r.Group("/api")
	protected.Use(authController.AuthMiddleware())
	{
		protected.GET("/profile", authController.GetProfile)
	}
}
