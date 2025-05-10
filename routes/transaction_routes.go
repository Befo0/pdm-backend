package routes

import (
	"pdm-backend/controllers"
	"pdm-backend/middlewares"
	"pdm-backend/repositories"

	"github.com/gin-gonic/gin"
)

func TransaccionRouter(r *gin.Engine) {
	transaccionRepo := repositories.NewTransaccionRepository(repositories.GetDB())
	handler := controllers.NewTransaccionHandler(transaccionRepo)

	transaccion := r.Group("/transacciones")
	transaccion.Use(middlewares.AuthMiddleware())
	{
		transaccion.GET("/lista-transaccion", handler.GetTransactions)
		transaccion.GET("/transaccion/:id", handler.GetTransactionById)
		transaccion.POST("/crear-transaccion", handler.CreateTransaction)
	}
}
