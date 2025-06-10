package repositories

import (
	"pdm-backend/models"

	"gorm.io/gorm"
)

type IngresosRepository struct {
	DB *gorm.DB
}

func NewIngresosRepository(db *gorm.DB) *IngresosRepository {
	return &IngresosRepository{DB: db}
}

type IngresosOpciones struct {
	IdIngreso          uint    `json:"id_ingreso"`
	NombreIngreso      string  `json:"nombre_ingreso"`
	IngresoPresupuesto float64 `json:"ingreso_presupuesto"`
}

func (r *IngresosRepository) GetIncomes(finanzaId uint) ([]IngresosOpciones, error) {

	opciones := []IngresosOpciones{}

	err := r.DB.Model(models.TipoIngresos{}).Where("finanzas_id = ?", finanzaId).
		Select("tipo_ingresos.id AS id_ingreso, tipo_ingresos.nombre_ingresos AS nombre_ingreso, tipo_ingresos.monto_ingreso AS ingreso_presupuesto").
		Scan(&opciones).Error

	if err != nil {
		return nil, err
	}

	return opciones, nil
}

type IngresosLista struct {
	IdIngreso     uint    `json:"id_ingreso"`
	NombreIngreso string  `json:"nombre_ingreso"`
	MontoIngreso  float64 `json:"monto_ingreso"`
	NombreUsuario string  `json:"nombre_usuario"`
}

func (r *IngresosRepository) GetIncomesList(finanzaId uint) ([]IngresosLista, error) {

	listaIngresos := []IngresosLista{}

	err := r.DB.Model(models.TipoIngresos{}).Where("finanzas_id = ?", finanzaId).
		Select("tipo_ingresos.id AS id_ingreso, tipo_ingresos.nombre_ingresos AS nombre_ingreso, tipo_ingresos.monto_ingreso AS monto_ingreso, users.nombre AS nombre_usuario").
		Joins("LEFT JOIN users ON users.id = tipo_ingresos.user_id").
		Scan(&listaIngresos).Error

	if err != nil {
		return nil, err
	}

	return listaIngresos, nil
}

func (r *IngresosRepository) CreateIncome(ingreso *models.TipoIngresos) error {
	return r.DB.Create(&ingreso).Error
}

func (r *IngresosRepository) GetIncomeById(id *uint) (*models.TipoIngresos, error) {
	var ingreso models.TipoIngresos

	if err := r.DB.First(&ingreso, id).Error; err != nil {
		return nil, err
	}

	return &ingreso, nil
}

type IncomeResponse struct {
	IdIngreso          uint    `json:"id_ingreso"`
	NombreIngreso      string  `json:"nombre_ingreso"`
	MontoIngreso       float64 `json:"monto_ingreso"`
	DescripcionIngreso string  `json:"descripcion_ingreso"`
}

func (r *IngresosRepository) GetIncome(id *uint) (*IncomeResponse, error) {
	var response IncomeResponse

	tx := r.DB.Model(models.TipoIngresos{}).Where("tipo_ingresos.id = ?", id).
		Select("tipo_ingresos.id AS id_ingreso, tipo_ingresos.nombre_ingresos AS nombre_ingreso, tipo_ingresos.monto_ingreso AS monto_ingreso, tipo_ingresos.descripcion AS descripcion_ingreso").
		Scan(&response)

	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if err := tx.Error; err != nil {
		return nil, err
	}

	return &response, nil
}

func (r *IngresosRepository) UpdateIncome(ingreso *models.TipoIngresos) error {
	return r.DB.Save(&ingreso).Error
}
