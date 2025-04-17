package models

import "gorm.io/gorm"

type Presupuesto struct {
	gorm.Model
	FinanzasID           uint               `json:"finanza_id" gorm:"index;not null"`
	Finanzas             Finanzas           `gorm:"foreignKey:FinanzasID"`
	EsCompartida         bool               `json:"es_conjunta" gorm:"not null"`
	FinanzasConjuntoID   *uint              `json:"finanza_conjunto_id" gorm:"index"`
	FinanzasConjunto     *FinanzasConjunto  `gorm:"foreignKey:FinanzasConjuntoID"`
	SubCategoriaEgresoID uint               `json:"sub_categoria_egreso_id" gorm:"index;not null"`
	SubCategoriaEgreso   SubCategoriaEgreso `gorm:"foreignKey:SubCategoriaEgresoID"`
	CategoriaEgresoID    uint               `json:"categoria_egreso_id" gorm:"index;not null"`
	CategoriaEgreso      CategoriaEgreso    `gorm:"foreignKey:CategoriaEgresoID"`
	TipoPresupuestoID    uint               `json:"tipo_presupuesto_id" gorm:"index;not null"`
	TipoPresupuesto      TipoPresupuesto    `gorm:"foreignKey:TipoPresupuestoID"`
	PresupuestoMensual   float64            `json:"presupuesto_mensual" gorm:"not null"`
	UserID               *uint              `json:"id_usuario_registro" gorm:"index"`
	User                 *User              `gorm:"foreignKey:UserID"`
}
