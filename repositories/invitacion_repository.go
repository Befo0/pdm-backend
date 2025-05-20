package repositories

import (
	"pdm-backend/models"
	"pdm-backend/services"
	"strings"
	"time"

	"gorm.io/gorm"
)

type InvitacionRepository struct {
	DB *gorm.DB
}

func NewInvitacionRepository(db *gorm.DB) *InvitacionRepository {
	return &InvitacionRepository{DB: db}
}

func (r *InvitacionRepository) CreateInvite(finanzaId *uint) error {

	intentosMax := 5

	for i := 0; i < intentosMax; i++ {

		codigo, err := services.GenerateInvitacionCode(10)
		if err != nil {
			return err
		}

		invitacion := models.Invitaciones{
			FinanzaID: *finanzaId,
			Codigo:    codigo,
			ExpiraEn:  time.Now().Add(time.Minute * 15),
		}

		err = r.DB.Create(&invitacion).Error
		if err == nil {
			return nil
		}

		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate key") {
			return err
		}

		return err
	}
	return nil
}
