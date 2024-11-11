package repository

import (
	"Booking_system/user_service/internal/models"
	"Booking_system/user_service/internal/util"
	"Booking_system/user_service/pkg/postgres"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type UserRepository struct {
	dbClient *gorm.DB
}

func NewUserRepository(dbClient *gorm.DB) *UserRepository {
	return &UserRepository{dbClient: dbClient}
}

func (repo *UserRepository) Delete(userEmail string) error {
	var user models.User
	err := repo.dbClient.Debug().Model(models.User{}).Where("email = ?", userEmail).First(&user).Error
	if err != nil {
		return err
	}

	err = repo.dbClient.Debug().Delete(&user).Error

	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetUserByEmail(userEmail string) (*models.User, error) {
	user := &models.User{}
	err := repo.dbClient.Debug().Where("email = ?", userEmail).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) GetUserEmailById(userId string) (*models.User, error) {
	user := &models.User{}

	err := repo.dbClient.Debug().Where("uuid = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) UpdateUser(user *models.User) error {
	existingUser := &models.User{}
	err := repo.dbClient.Where("email = ?", user.Email).First(&existingUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("User with email %s  and uuid %s not found", user.Email, user.UUID)
		}
		return err
	}
	existingUser.Phone = user.Phone
	existingUser.LastName = user.LastName
	existingUser.FirstName = user.FirstName
	existingUser.DateOfBirth = user.DateOfBirth
	existingUser.Address = user.Address

	err = repo.dbClient.Save(&existingUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetAllUsers(pagination util.Pagination) ([]models.User, string) {
	var users []models.User
	var totalNumber string
	var totalInt int64
	repo.dbClient.Scopes(postgres.Paginate(pagination)).Preload("Status").Order("ID").Find(&users)
	repo.dbClient.Model(&models.User{}).Count(&totalInt)
	totalNumber = strconv.FormatInt(totalInt, 10)
	return users, totalNumber
}

func (repo *UserRepository) AddAdminRole(userEmail string) error {
	user := &models.User{}
	err := repo.dbClient.Debug().Model(models.User{}).Where("email = ?", userEmail).First(user).Error
	if err != nil {
		return err
	}
	if user.Role == "admin" {
		return nil
	}

	err = repo.dbClient.Debug().Model(models.User{}).Where("email = ?", userEmail).Update("role", "admin").Error

	return err
}

func (repo *UserRepository) SetUserRole(userEmail string) error {
	user := &models.User{}
	err := repo.dbClient.Debug().Model(models.User{}).Where("email = ?", userEmail).First(user).Error
	if err != nil {
		return err
	}
	if user.Role == "user" {
		return nil
	}
	err = repo.dbClient.Debug().Model(models.User{}).Where("email = ?", userEmail).Update("role", "user").Error

	return err
}

func (repo *UserRepository) CreateUser(user *models.User) error {
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
