package models

import "gorm.io/gorm"

type FinanzasConjunto struct {
	gorm.Model
	FinanzasID uint `json:"finanza_id" gorm:"index;not null"`
	UserID     uint `json:"user_id" gorm:"index;not null"`
	RolesID    uint `json:"rol_id" gorm:"not null"`
}
