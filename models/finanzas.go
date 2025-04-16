package models

import "gorm.io/gorm"

type Finanzas struct {
	gorm.Model
	UserID        uint   `json:"user_id" gorm:"index;not null"`
	TipoFinanzaID uint   `json:"tipo_finanza_id" gorm:"index;not null"`
	Titulo        string `json:"titulo" gorm:"size:255;not null`
	Descripcion   string `json:"descripcion" gorm:"size:500"`
}
