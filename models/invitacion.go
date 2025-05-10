package models

import (
	"time"

	"gorm.io/gorm"
)

type Invitaciones struct {
	gorm.Model
	FinanzaConjuntoID uint             `gorm:"index;not null" json:"finanza_conjunto_id"`
	Codigo            string           `gorm:"size:500;not null" json:"codigo"`
	ExpiraEn          time.Time        `gorm:"not null" json:"expira_en"`
	FinanzasConjunto  FinanzasConjunto `gorm:"foreignKey:FinanzaConjuntoID" json:"finanza_conjunto"`
}
