package models

import (
	"time"

	"gorm.io/gorm"
)

type Invitaciones struct {
	gorm.Model
	FinanzaID uint      `gorm:"index;not null" json:"finanza_conjunto_id"`
	Codigo    string    `gorm:"size:10;not null;uniqueIndex" json:"codigo"`
	ExpiraEn  time.Time `gorm:"not null" json:"expira_en"`
	Finanzas  Finanzas  `gorm:"foreignKey:FinanzasID" json:"finanza"`
}
