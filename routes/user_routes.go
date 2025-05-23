package routes

import (
	"pdm-backend/controllers"
	"pdm-backend/middlewares"
	"pdm-backend/repositories"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine) {

	userRepo := repositories.NewUserRepository(repositories.GetDB())
	financeRepo := repositories.NewFinanzaRepository(repositories.GetDB())
	handler := controllers.NewUserHandler(userRepo, financeRepo)

	user := r.Group("/user")
	user.POST("/login", handler.Login)
	user.POST("/register", handler.Register)

	user.Use(middlewares.AuthMiddleware())
	{
		user.PATCH("/change-profile", handler.UpdateProfile)
		user.PATCH("/change-password", handler.UpdatePassword)
	}
}
