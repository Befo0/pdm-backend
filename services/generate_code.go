package services

import (
	"crypto/rand"
	"math/big"
)

const caracteres = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func GenerateInvitacionCode(nCharacters int) (string, error) {
	codigo := make([]byte, nCharacters)

	for i := range codigo {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(caracteres))))
		if err != nil {
			return "", err
		}
		codigo[i] = caracteres[num.Int64()]
	}

	return string(codigo), nil
}
