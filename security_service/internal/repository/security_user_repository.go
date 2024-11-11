package repository

import (
	"Booking_system/security_service/internal/models"
	"errors"
	"gorm.io/gorm"
)

type UserRepository struct {
	dbClient *gorm.DB
}

func NewUserRepository(dbClient *gorm.DB) *UserRepository {
	return &UserRepository{
		dbClient: dbClient,
	}
}

func (repo *UserRepository) Create(user *models.User) error {
	err := repo.dbClient.Debug().Model(models.User{}).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) ValidateUserExistence(userEmail string) (*models.User, error) {
	user := &models.User{}
	err := repo.dbClient.Debug().Where("email = ?", userEmail).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err == gorm.ErrRecordNotFound || user.UUID == "0" {
		return nil, nil
	}
	return user, nil
}

func (repo *UserRepository) CheckIfEmailExists(mail string) bool {
	var user models.User

	err := repo.dbClient.Debug().Model(models.User{}).Find(&user).Where("email= ?", mail).Error

	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (repo *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	err := repo.dbClient.Debug().Model(models.User{}).Where("email = ?", email).First(&user).Error

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
