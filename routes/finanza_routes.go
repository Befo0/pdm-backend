package routes

import (
	"pdm-backend/controllers"
	"pdm-backend/middlewares"
	"pdm-backend/repositories"

	"github.com/gin-gonic/gin"
)

func FinanzaRouter(r *gin.Engine) {

	financeRepo := repositories.NewFinanzaRepository(repositories.GetDB())
	handler := controllers.NewFinanzaHandler(financeRepo)

	finanza := r.Group("/finanza")
	finanza.Use(middlewares.AuthMiddleware())
	{
		finanza.GET("/resumen", handler.GetDashboardSummary)
		finanza.GET("/datos", handler.GetDashboardData)
	}
}
