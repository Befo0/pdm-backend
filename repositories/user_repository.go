package repositories

import (
	"pdm-backend/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(&user).Error
}

func (r *UserRepository) CreateUserAndFinance(user *models.User) error {

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		finanza := models.Finanzas{
			UserID:         user.ID,
			TipoFinanzasID: 1,
		}
		if err := tx.Create(&finanza).Error; err != nil {
			return err
		}

		categoria := models.CategoriaEgreso{
			FinanzasID:      finanza.ID,
			NombreCategoria: "Ahorro",
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
			EsCompartida:       false,
		}
		if err := tx.Create(&subCategoria).Error; err != nil {
			return err
		}

		ahorro := models.Ahorro{
			FinanzasID: finanza.ID,
			Monto:      0,
		}
		if err := tx.Create(&ahorro).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.DB.Where("correo = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserById(id uint) (*models.User, error) {
	var user models.User

	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetFinanceByUserId(userId uint) (uint, error) {
	var finanza models.Finanzas

	if err := r.DB.Where("user_id = ? AND tipo_finanzas_id", userId, 1).First(&finanza).Error; err != nil {
		return 0, err
	}

	return finanza.ID, nil
}
