package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name             string             `json:"nombre" gorm:"not null"`
	Email            string             `json:"correo" gorm:"unique"`
	Password         string             `json:"-" gorm:"size255;not null"`
	Finanzas         []Finanzas         `json:"finanzas" gorm:"foreignKey:UserID"`
	FinanzasConjunto []FinanzasConjunto `json:"finanzas_conjunto" gorm:"foreignKey:UserID"`
}
