package models

import (
	"time"

	"gorm.io/gorm"
)

type Invitaciones struct {
	gorm.Model
	FinanzasID uint      `json:"finanza_id" gorm:"index;not null"`
	Codigo     string    `gorm:"size:10;not null;uniqueIndex" json:"codigo"`
	ExpiraEn   time.Time `gorm:"not null" json:"expira_en"`
	Finanzas   Finanzas  `gorm:"foreignKey:FinanzasID" json:"finanza"`
}
