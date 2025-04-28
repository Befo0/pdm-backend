package routes

import (
	"pdm-backend/controllers"
	"pdm-backend/repositories"

	"github.com/gin-gonic/gin"
)

func FinanzaRouter(r *gin.Engine) {

	financeRepo := repositories.NewFinanzaRepository(repositories.GetDB())
	handler := controllers.NewFinanzaHandler(financeRepo)

	user := r.Group("/finanza")
	user.GET("/principal", handler.GetDashboard)
}
