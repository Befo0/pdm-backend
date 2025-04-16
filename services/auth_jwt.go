package services

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret string

type JWTClaims struct {
	UserId    uint   `json:"userId"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
	jwt.RegisteredClaims
}

func init() {
	secret = os.Getenv("SECRET_WORD")
	if secret == "" {
		log.Fatal("No se ha encontrado la clave secreta")
	}
}

func GenerateJWT(userId uint, userName string, userEmail string) (string, error) {

	claims := JWTClaims{
		UserId:    userId,
		UserName:  userName,
		UserEmail: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
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
