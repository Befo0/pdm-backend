package routes

import (
	"pdm-backend/controllers"
	"pdm-backend/repositories"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {

	userRepo := repositories.NewUserRepository(repositories.GetDB())
	userHandler := controllers.NewUserHandler(userRepo)

	user := r.Group("/user")
	user.POST("/login")
	user.POST("/register", userHandler.Register)
}
