package repositories

import (
	"time"

	"gorm.io/gorm"
)

type FinanzaRepository struct {
	DB *gorm.DB
}

func NewFinanzaRepository(db *gorm.DB) *FinanzaRepository {
	return &FinanzaRepository{DB: db}
}

func (r *FinanzaRepository) GetSummary(inicio, final time.Time) (interface{ any }, error) {

}
