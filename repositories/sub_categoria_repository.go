package repositories

import (
	"pdm-backend/models"

	"gorm.io/gorm"
)

type SubCategoriaRepository struct {
	DB *gorm.DB
}

func NewSubCategoriaRepository(db *gorm.DB) *SubCategoriaRepository {
	return &SubCategoriaRepository{DB: db}
}

type SubCategoriasFinanzas struct {
	SubCategoriaId     uint
	SubCategoriaNombre string
}

func (r *SubCategoriaRepository) GetSubCategories(finanzaId uint) ([]SubCategoriasFinanzas, error) {

	var subCategorias []SubCategoriasFinanzas

	err := r.DB.Model(models.SubCategoriaEgreso{}).Where("finanzas_id = ?", finanzaId).
		Select("sub_categoria_egresos.id AS sub_categoria_id, sub_categoria_egresos.nombre_sub_categoria AS sub_categoria_nombre").
		Scan(&subCategorias).Error

	if err != nil {
		return nil, err
	}

	return subCategorias, err
}

type GastosOpciones struct {
	TipoId     uint
	TipoNombre string
}

func (r *SubCategoriaRepository) GetSubCategoriesExpensesType() ([]GastosOpciones, error) {

	var opciones []GastosOpciones

	err := r.DB.Model(models.TipoPresupuesto{}).Select("tipo_presupuestos.id AS tipo_id, tipo_presupuestos.nombre_tipo_presupuesto AS tipo_nombre").
		Scan(&opciones).Error
	if err != nil {
		return nil, err
	}

	return opciones, nil
}

type SubCategoriasLista struct {
	SubCategoriaId     uint
	CategoriaNombre    string
	SubCategoriaNombre string
	TipoGasto          string
	Presupuesto        float64
	NombreUsuario      string
}

func (r *SubCategoriaRepository) GetSubCategoriesList(finanzaId uint) ([]SubCategoriasLista, error) {
	var subCategoria []SubCategoriasLista

	err := r.DB.Model(models.SubCategoriaEgreso{}).Where("sub_categoria_egresos.finanzas_id = ?", finanzaId).
		Select("sub_categoria_egresos.id AS sub_categoria_id, categoria_egresos.nombre_categoria AS categoria_nombre, sub_categoria_egresos.nombre_sub_categoria AS sub_categoria_nombre, tipo_presupuestos.nombre_tipo_presupuesto AS tipo_gasto, sub_categoria_egresos.presupuesto_mensual AS presupuesto, users.nombre AS nombre_usuario").
		Joins("LEFT JOIN categoria_egresos ON sub_categoria_egresos.categoria_egreso_id = categoria_egresos.id").
		Joins("LEFT JOIN tipo_presupuestos ON sub_categoria_egresos.tipo_presupuesto_id = tipo_presupuestos.id").
		Joins("LEFT JOIN users ON sub_categoria_egresos.user_id = users.id").
		Scan(&subCategoria).Error

	if err != nil {
		return nil, err
	}

	return subCategoria, nil
}

func (r *SubCategoriaRepository) CreateSubCategory(subCategoria *models.SubCategoriaEgreso) error {
	return r.DB.Create(&subCategoria).Error
}

func (r *SubCategoriaRepository) GetSubCategoryById(id *uint) (*models.SubCategoriaEgreso, error) {

	var subCategoria models.SubCategoriaEgreso

	if err := r.DB.First(&subCategoria, id).Error; err != nil {
		return nil, err
	}

	return &subCategoria, nil
}

type SubCategoriaResponse struct {
	CategoriaId        uint    `json:"categoria_id"`
	NombreSubCategoria string  `json:"nombre_sub_categoria"`
	TipoGastoId        uint    `json:"tipo_gasto_id"`
	Presupuesto        float64 `json:"presupuesto"`
}

func (r *SubCategoriaRepository) GetSubCategory(id *uint) (*SubCategoriaResponse, error) {

	var subCategoria SubCategoriaResponse

	tx := r.DB.Model(models.SubCategoriaEgreso{}).Where("sub_categoria_egresos.id = ?", id).
		Select("sub_categoria_egresos.categoria_egreso_id AS categoria_id, sub_categoria_egresos.nombre_sub_categoria AS nombre_sub_categoria, sub_categoria_egresos.tipo_presupuesto_id AS tipo_gasto_id, sub_categoria_egresos.presupuesto_mensual AS presupuesto").Scan(&subCategoria)

	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if err := tx.Error; err != nil {
		return nil, err
	}

	return &subCategoria, nil
}
