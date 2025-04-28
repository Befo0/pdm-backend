package models

import (
	"time"

	"gorm.io/gorm"
)

type Transacciones struct {
	gorm.Model
	FinanzasID           uint                `json:"finanza_id" gorm:"index;not null"`
	Finanzas             Finanzas            `gorm:"foreignKey:FinanzasID"`
	EsCompartida         bool                `json:"es_compartida" gorm:"not null"`
	FinanzasConjuntoID   *uint               `json:"finanza_conjunto_id" gorm:"index"`
	FinanzasConjunto     *FinanzasConjunto   `gorm:"foreignKey:FinanzasConjuntoID"`
	Descripcion          *string             `json:"descripcion" gorm:"size:500"`
	TipoRegistroID       uint                `json:"tipo_registro_id" gorm:"index;not null"`
	TipoRegistro         TipoRegistro        `gorm:"foreignKey:TipoRegistroID"`
	TipoIngresosID       *uint               `json:"tipo_ingreso_id" gorm:"index"`
	TipoIngresos         *TipoIngresos       `gorm:"foreignKey:TipoIngresosID"`
	CategoriaEgresoID    *uint               `json:"categoria_egreso_id" gorm:"index"`
	CategoriaEgreso      *CategoriaEgreso    `gorm:"foreignKey:CategoriaEgresoID"`
	SubCategoriaEgresoID *uint               `json:"sub_categoria_egreso_id" gorm:"index"`
	SubCategoriaEgreso   *SubCategoriaEgreso `gorm:"foreignKey:SubCategoriaEgresoID"`
	TipoPresupuestoID    *uint               `json:"tipo_gasto_id" gorm:"index"`
	TipoPresupuesto      *TipoPresupuesto    `gorm:"foreignKey:TipoPresupuestoID"`
	FechaRegistro        time.Time           `json:"fecha_registro" gorm:"not null"`
	Monto                float64             `json:"monto" gorm:"not null"`
}
