package models

import (
	"time"

	"gorm.io/gorm"
)

type FinanzasConjunto struct {
	gorm.Model
	FinanzasID uint               `json:"finanza_id" gorm:"index;not null"`
	Finanzas   Finanzas           `gorm:"foreignKey:FinanzasID"`
	UserID     uint               `json:"user_id" gorm:"index;not null"`
	User       User               `gorm:"foreignKey:UserID"`
	RolesID    uint               `json:"rol_id" gorm:"not null"`
	Rol        RolFinanzaConjunto `gorm:"foreignKey:RolesID" json:"rol"`
	FechaUnion time.Time          `json:"fecha_union" gorm:"not null"`
}
