package routes

import (
	"pdm-backend/controllers"
	"pdm-backend/middlewares"
	"pdm-backend/repositories"

	"github.com/gin-gonic/gin"
)

func TransaccionRouter(r *gin.Engine) {
	transaccionRepo := repositories.NewTransaccionRepository(repositories.GetDB())
	handler := controllers.NewTransaccionHanlder(transaccionRepo)

	transaccion := r.Group("/transaccion")
	transaccion.Use(middlewares.AuthMiddleware())
	{
	}
}
