package services

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var secret string

type JWTClaims struct {
	UserId    uint   `json:"userId"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
	FinanzaId uint   `json:"financeId"`
	jwt.RegisteredClaims
}

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("La variable de entorno no ha sido cargada")
	}

	secret = os.Getenv("SECRET_WORD")
	if secret == "" {
		log.Fatal("No se ha encontrado la clave secreta")
	}
}

func GenerateJWT(userId uint, userName string, userEmail string, finanzaId uint) (string, error) {

	claims := JWTClaims{
		UserId:    userId,
		UserName:  userName,
		UserEmail: userEmail,
		FinanzaId: finanzaId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateJWT(cookieToken string) (*jwt.Token, *JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(cookieToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Metodo de firma invalido: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, nil, err
	}

	return token, claims, nil
}
