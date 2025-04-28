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

type Handler struct {
	UserRepo    *repositories.UserRepository
	FinanceRepo *repositories.FinanzaRepository
}

func NewUserHandler(userRepo *repositories.UserRepository, financeRepo *repositories.FinanzaRepository) *Handler {
	return &Handler{
		UserRepo:    userRepo,
		FinanceRepo: financeRepo,
	}
}

func (h *Handler) Register(c *gin.Context) {

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

	err = h.UserRepo.CreateUserAndFinance(&user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al crear el usuario y la finanza"})
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

func (h *Handler) Login(c *gin.Context) {

	var userRequest LoginRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la peticion esta incorrecto"})
		return
	}

	user, err := h.UserRepo.GetUserByEmail(userRequest.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"success": false, "errors": gin.H{"email": "No hay cuenta asociada a este correo electronico"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"sucess": false, "message": "Error en el servidor"})
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

type UpdateProfileRequest struct {
	Name  string `json:"nombre"`
	Email string `json:"correo"`
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var updateRequest UpdateProfileRequest

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la petición es incorrecto"})
		return
	}

	claimsInterface, exists := c.Get("claims")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "No se encontraron claims"})
		return
	}

	userClaims, ok := claimsInterface.(*services.JWTClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Tipo de claim invalido"})
		return
	}

	user, err := h.UserRepo.GetUserById(userClaims.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "El usuario no existe"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"sucess": false, "message": "Error en el servidor"})
		return
	}

	user.Name = updateRequest.Name
	user.Email = updateRequest.Email

	if err := h.UserRepo.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al modificar datos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "El perfil ha sido actualizado correctamente"})
}

type UpdatePasswordRequest struct {
	ActualPassword  string `json:"actual_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (h *Handler) UpdatePassword(c *gin.Context) {
	var passwordRequest UpdatePasswordRequest

	if err := c.ShouldBindJSON(&passwordRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "El formato de la petición es incorrecto"})
		return
	}

	claimsInterface, exists := c.Get("claims")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "No se encontraron claims"})
		return
	}

	userClaims, ok := claimsInterface.(*services.JWTClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Tipo de claim invalido"})
		return
	}

	user, err := h.UserRepo.GetUserById(userClaims.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "El usuario no existe"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"sucess": false, "message": "Error en el servidor"})
		return
	}

	password := user.Password

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(passwordRequest.ActualPassword)); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "La contraseña no coincide con la actual"})
		return
	}

	if passwordRequest.NewPassword != passwordRequest.ConfirmPassword {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "La confirmación de tu contraseña no coincide con la nueva"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordRequest.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Hubo un error al hashear la contraseña"})
		return
	}

	user.Password = string(hashedPassword)

	if err := h.UserRepo.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al modificar la contraseña"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "La contraseña ha sido actualizado correctamente"})
}
