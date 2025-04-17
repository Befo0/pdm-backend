package controllers

import (
	"errors"
	"net/http"
	"pdm-backend/models"
	"pdm-backend/repositories"
	"pdm-backend/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	Repo *repositories.UserRepository
}

func NewUserHandler(repo *repositories.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la petición esta incorrecto"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ocurrio un error al hashaer la contraseña"})
		return
	}

	user.Password = string(hashedPassword)

	err = h.Repo.Create(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al crear el usuario"})
		return
	}

	token, err := services.GenerateJWT(user.ID, user.Name, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "No se pudo cargar el token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

type LoginRequest struct {
	Email    string `json:"correo"`
	Password string `json:"contraseña"`
}

func (h *UserHandler) Login(c *gin.Context) {

	var userRequest LoginRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la peticion esta incorrecto"})
		return
	}

	user, err := h.Repo.GetUserByEmail(userRequest.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "errors": gin.H{"email": "No hay cuenta asociada a este correo electronico"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"sucess": false, "message": err.Error()})
		return
	}

	password := user.Password

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(userRequest.Password)); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "La contraseña proporcionada es incorrecta"})
		return
	}

	token, err := services.GenerateJWT(user.ID, user.Name, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "No se pudo cargar el token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "token": token})
}
