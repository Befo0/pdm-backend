package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nombre           string             `json:"nombre" gorm:"not null"`
	Correo           string             `json:"correo" gorm:"unique"`
	Contrasena       string             `json:"contrasena" gorm:"size255;not null"`
	Finanzas         []Finanzas         `json:"finanzas" gorm:"foreignKey:UserID"`
	FinanzasConjunto []FinanzasConjunto `json:"finanzas_conjunto" gorm:"foreignKey:UserID"`
}
