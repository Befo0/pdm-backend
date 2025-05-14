package repositories

import (
	"pdm-backend/models"

	"gorm.io/gorm"
)

type AhorroRepository struct {
	DB *gorm.DB
}

func NewAhorroRepository(db *gorm.DB) *AhorroRepository {
	return &AhorroRepository{DB: db}
}

func (r *AhorroRepository) GetSaving(finanzaId uint) (*models.Ahorro, error) {

	var ahorro models.Ahorro

	err := r.DB.Where("finanzas_id = ?", finanzaId).First(&ahorro).Error
	if err != nil {
		return nil, err
	}

	return &ahorro, nil
}

func (r *AhorroRepository) UpdateSaving(ahorro *models.Ahorro) error {
	return r.DB.Save(&ahorro).Error
}

type AhorroData struct {
	NombreMes    string
	AhorroMes    float64
	AhorroAcum   float64
	MetaMensual  float64
	Cumplimiento float64
}

func (r *AhorroRepository) GetSavingsData(mensualidad float64, finanzaId uint) (*AhorroData, error) {

	var subCategoriaId uint
	errors := make(chan error, 3)
	var ahorroData []AhorroData
	meses := []string{"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}

	tx := r.DB.Model(models.SubCategoriaEgreso{}).Where("nombre_sub_categoria = ?", "Ahorro").Select("sub_categoria_egresos.id").Scan(&subCategoriaId)
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if err := tx.Error; err != nil {
		return nil, err
	}

	for index, meses := range meses {

		var ahorroMes float64
		var ahorroAcum float64
		metaMensual := float64((index + 1) * 25)
		cumplimiento := 0.0

		go func() {
			err := r.DB.Model(models.Transacciones{}).Where("finanzas_id = ? AND sub_categoria_egreso_id = ? AND EXTRACT(MONTH FROM fecha_registro) = ?", finanzaId, subCasubCategoriaId, index+1).Error

			errors <- err
		}()
	}
}
