package routes

import (
	"pdm-backend/controllers"
	"pdm-backend/middlewares"
	"pdm-backend/repositories"

	"github.com/gin-gonic/gin"
)

func AhorroRouter(r *gin.Engine) {

	ahorroRepo := repositories.NewAhorroRepository(repositories.GetDB())
	handler := controllers.NewAhorroHandler(ahorroRepo)

	ahorro := r.Group("/ahorro")
	ahorro.Use(middlewares.AuthMiddleware())
	{
		ahorro.GET("/lista", handler.GetSavingsData)
		ahorro.PUT("/actualizar", handler.UpdateSaving)
	}
}
