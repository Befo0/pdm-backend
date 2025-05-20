package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nombre           string             `json:"nombre" gorm:"not null"`
	Correo           string             `json:"correo" gorm:"uniqueIndex"`
	Contrasena       string             `json:"contrasena" gorm:"size:255;not null"`
	FinanzasConjunto []FinanzasConjunto `json:"finanzas_conjunto" gorm:"foreignKey:UserID"`
}
