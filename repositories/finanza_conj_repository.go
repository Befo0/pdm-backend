package repositories

import (
	"errors"
	"pdm-backend/models"
	"time"

	"gorm.io/gorm"
)

type FinanzaConjRepository struct {
	DB *gorm.DB
}

func NewFinanzaConjRepository(db *gorm.DB) *FinanzaConjRepository {
	return &FinanzaConjRepository{DB: db}
}

func (r *FinanzaConjRepository) CreateConjFinance(userId uint, titulo, descripcion string) error {

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		finanza := models.Finanzas{
			UserID:         userId,
			TipoFinanzasID: 2,
			Titulo:         &titulo,
			Descripcion:    &descripcion,
		}
		if err := tx.Create(&finanza).Error; err != nil {
			return err
		}

		categoria := models.CategoriaEgreso{
			FinanzasID:      finanza.ID,
			NombreCategoria: "Ahorro",
			UserID:          userId,
		}
		if err := tx.Create(&categoria).Error; err != nil {
			return err
		}

		subCategoria := models.SubCategoriaEgreso{
			FinanzasID:         finanza.ID,
			NombreSubCategoria: "Ahorro",
			TipoPresupuestoID:  3,
			CategoriaEgresoID:  categoria.ID,
			PresupuestoMensual: 0.00,
			UserID:             userId,
		}
		if err := tx.Create(&subCategoria).Error; err != nil {
			return err
		}

		finanzaConj := models.FinanzasConjunto{
			FinanzasID: finanza.ID,
			UserID:     userId,
			RolesID:    1,
			Activo:     true,
			FechaUnion: time.Now(),
		}
		if err := tx.Create(&finanzaConj).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *FinanzaConjRepository) JoinUser(userId uint, codigo string) error {

	var invitacion models.Invitaciones

	err := r.DB.Where("codigo = ?", codigo).First(&invitacion).Error
	if err != nil {
		return err
	}

	var existente models.FinanzasConjunto
	err = r.DB.Where("finanzas_id = ? AND user_id = ?", invitacion.FinanzaID, userId).First(&existente).Error
	if err == nil {
		return errors.New("Ya perteneces a esta finanza conjunta")
	}

	if time.Now().Before(invitacion.ExpiraEn) {
		finanzaConjunta := models.FinanzasConjunto{
			FinanzasID: invitacion.FinanzaID,
			UserID:     userId,
			RolesID:    2,
			Activo:     true,
			FechaUnion: time.Now(),
		}

		err := r.DB.Create(&finanzaConjunta).Error
		if err != nil {
			return err
		}

	} else {
		return errors.New("El codigo ya ha expirado")
	}

	return nil
}

type FinancesResponse struct {
	FinanzaId     uint
	FinanzaNombre string
	NombreAdmin   string
}

func (r *FinanzaConjRepository) GetConjFinances(userId uint) ([]FinancesResponse, error) {

	var financeResponse []FinancesResponse

	err := r.DB.Model(models.FinanzasConjunto{}).Where("user_id = ? AND activo = ?", userId, true).
		Select("finanzas.id AS finanza_id, finanzas.titulo AS finanza_nombre, users.nombre AS nombre_admin").
		Joins("INNER JOIN finanzas ON finanzas.id = finanzas_conjuntos.finanzas_id").
		Joins("LEFT JOIN finanzas_conjuntos ON finanzas_conjuntos.finanzas_id = finanzas.id AND finanzas_conjuntos.roles_id = 1").
		Joins("LEFT JOIN users ON users.id = finanzas_conjuntos.user_id").
		Scan(&financeResponse).Error
	if err != nil {
		return nil, err
	}

	return financeResponse, nil
}

type MiembrosFinanza struct {
	IdUsuario     uint
	NombreUsuario string
	RolUsuario    uint
}

type ConjFinancesDetails struct {
	FinanzaTitulo      string
	FinanzaDescripcion string
	Miembros           []MiembrosFinanza
}

func (r *FinanzaConjRepository) GetConjFinancesDetails(finanzaId uint) (*ConjFinancesDetails, error) {
	var detallesFinanza ConjFinancesDetails
	var miembros []MiembrosFinanza

	err := r.DB.Model(models.Finanzas{}).
		Select("finanzas.titulo AS finanza_titulo, finanzas.descripcion AS finanza_descripcion").
		Where("finanzas.id = ?", finanzaId).
		Scan(&detallesFinanza).Error
	if err != nil {
		return nil, err
	}

	err = r.DB.Model(models.FinanzasConjunto{}).
		Select("finanzas_conjuntos.user_id AS id_usuario, users.nombre AS nombre_usuario, finanzas_conjuntos.roles_id AS rol_usuario").
		Joins("JOIN users ON users.id = finanzas_conjuntos.user_id").
		Where("finanzas_conjuntos.finanzas_id = ?", finanzaId).
		Scan(&miembros).Error
	if err != nil {
		return nil, err
	}

	detallesFinanza.Miembros = miembros
	return &detallesFinanza, nil
}

func (r *FinanzaConjRepository) LeaveConjFinance(userId, finanzaId uint) error {

	var finanzaConj models.FinanzasConjunto

	err := r.DB.Model(models.FinanzasConjunto{}).Where("finanzas_id = ? AND user_id = ?", finanzaId, userId).
		First(&finanzaConj).Error
	if err != nil {
		return err
	}

	finanzaConj.Activo = false

	return r.DB.Save(&finanzaConj).Error
}
