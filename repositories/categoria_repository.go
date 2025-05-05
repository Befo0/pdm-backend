package repositories

import (
	"pdm-backend/models"

	"gorm.io/gorm"
)

type CategoriaRepository struct {
	DB *gorm.DB
}

func NewCategoriaRepository(db *gorm.DB) *CategoriaRepository {
	return &CategoriaRepository{DB: db}
}

type CategoriasFinanzas struct {
	CategoriaId     uint
	CategoriaNombre string
}

func (r *CategoriaRepository) GetCategories(finanzaId uint) (*[]CategoriasFinanzas, error) {

	var categorias []CategoriasFinanzas

	err := r.DB.Model(models.CategoriaEgreso{}).Where("finanzas_id = ?", finanzaId).Select("categoria_egresos.id AS categoria_id, categoria_egresos.nombre_categoria AS categoria_nombre").Scan(&categorias).Error
	if err != nil {
		return nil, err
	}

	return &categorias, err

}
