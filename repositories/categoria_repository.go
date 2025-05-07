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

type SubCategorias struct {
	Nombre      string
	Presupuesto float64
	Gasto       float64
	Diferencia  float64
}

type JSONResponse struct {
	Presupuesto   float64
	Gasto         float64
	Diferencia    float64
	SubCategorias []SubCategorias
}

func (r *CategoriaRepository) GetCategoriesData(finanzaId uint, categoriaId uint) (*JSONResponse, error) {

	var respuesta JSONResponse
	err := r.DB.Model(models.SubCategoriaEgreso{}).Where("sub_categoria_egresos.finanzas_id = ? AND sub_categoria_egresos.categoria_egreso_id = ?", finanzaId, categoriaId).Select("COALESCE(SUM(sub_categoria_egresos.presupuesto_mensual), 0) AS presupuesto, COALESCE(SUM(transacciones.monto), 0) AS gasto").Joins("LEFT JOIN transacciones on transacciones.sub_categoria_egresos_id = sub_categoria_egresos.id").Scan(&respuesta).Error
	if err != nil {
		return nil, err
	}

	respuesta.Diferencia = respuesta.Presupuesto - respuesta.Gasto

	var subCategorias []SubCategorias
	err = r.DB.Model(models.SubCategoriaEgreso{}).Where("sub_categoria_egresos.finanzas_id = ? AND sub_categoria_egresos.categoria_egreso_id = ?").Select("sub_categoria_egresos.nombre_sub_categoria AS nombre, COALESCE(SUM(sub_categoria_egresos.presupuesto_mensual),0) AS presupuesto , COALESCE(SUM(transacciones.monto),0) AS gasto").Joins("LEFT JOIN transacciones ON transacciones.sub_categoria_egreso_id = sub_categoria_egresos.id").Group("sub_categoria_egresos.id").Scan(&subCategorias).Error
	if err != nil {
		return nil, err
	}

	for index := range subCategorias {
		subCategorias[index].Diferencia = subCategorias[index].Presupuesto - subCategorias[index].Gasto
	}

	respuesta.SubCategorias = subCategorias

	return &respuesta, nil
}
