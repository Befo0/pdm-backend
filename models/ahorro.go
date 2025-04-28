package models

import "gorm.io/gorm"

type Ahorro struct {
	gorm.Model
	FinanzasID uint     `gorm:"index;not null" json:"finanza_id"`
	Monto      float64  `gorm:"not null" json:"monto"`
	Finanzas   Finanzas `gorm:"foreignKey:FinanzasID" json:"finanza"`
}
