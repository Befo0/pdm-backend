package repositories

import (
	"gorm.io/gorm"
)

type CategoriaRepository struct {
	DB *gorm.DB
}

func NewCategoriaRepository(db *gorm.DB) *CategoriaRepository {
	return &CategoriaRepository{DB: db}
}
