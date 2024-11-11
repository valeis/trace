package service

import (
	"Booking_system/user_service/internal/models"
	"Booking_system/user_service/internal/util"
	"errors"
	"github.com/google/uuid"
	"github.com/gookit/slog"
	"golang.org/x/crypto/bcrypt"
)

type IUserRepository interface {
	Delete(userEmail string) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserEmailById(userID string) (*models.User, error)
	UpdateUser(user *models.User) error
	GetAllUsers(paginaton util.Pagination) ([]models.User, string)
	AddAdminRole(userEmail string) error
	SetUserRole(userEmail string) error
	CreateUser(user *models.User) error
	ValidateUserExistence(userEmail string) (*models.User, error)
}

type UserService struct {
	userRepo IUserRepository
}

func NewUserService(userRepo IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (svc *UserService) Delete(userEmail string) error {
	err := svc.userRepo.Delete(userEmail)
	if err != nil {
		slog.Errorf("failed to delete user from database: %v", err.Error())
		return err
	}
	return nil
}

func (svc *UserService) GetUserByEmail(userEmail string) (*models.User, error) {
	user, err := svc.userRepo.GetUserByEmail(userEmail)
	if err != nil {
		slog.Errorf("failed to get user from database: %v\n", err)
		return nil, err
	}
	return user, nil
}

func (svc *UserService) GetUserEmailById(userID string) (*models.User, error) {
	user, err := svc.userRepo.GetUserEmailById(userID)
	if err != nil {
		slog.Errorf("Failed to fetch user from database: %v", err.Error())
		return nil, err
	}
	return user, nil
}

func (svc *UserService) UpdateUser(user *models.User) error {
	err := svc.userRepo.UpdateUser(user)
	if err != nil {
		slog.Errorf("Error updating user data: %v\n", err.Error())
	}
	return nil
}

func (svc *UserService) GetAllUsers(pagination util.Pagination) ([]models.User, string) {
	return svc.userRepo.GetAllUsers(pagination)
}

func (svc *UserService) AddAdmin(userEmail string) error {
	err := svc.userRepo.AddAdminRole(userEmail)
	if err != nil {
		slog.Error("failed to add admin role: %v\n", err)
		return err
	}
	return nil
}

func (svc *UserService) SetUserRole(userEmail string) error {
	err := svc.userRepo.SetUserRole(userEmail)
	if err != nil {
		slog.Errorf("failed to set user role to user: %v\n", err)
		return err
	}
	return nil
}

func (svc *UserService) CreateUser(user *models.User) (string, error) {
	existingUser, err := svc.userRepo.ValidateUserExistence(user.Email)

	if err != nil && existingUser == nil {
		slog.Errorf("Err validating user existence: %v\n", err.Error())
		return "Error adding user", err
	} else if existingUser == nil && err == nil {
		hashedPassword, err := svc.generatePasswordHash(user.Password)
		if err != nil {
			slog.Errorf("Can't register new user: %v\n", err.Error())
			return "Error adding user", errors.New("can't register")
		}
		user.Password = hashedPassword
		user.UUID = uuid.New().String()
		err = svc.userRepo.CreateUser(user)
		if err != nil {
			slog.Errorf("Failed to insert user in database: %v\n", err.Error())
			return "Error adding user", err
		}
		slog.Info("User added successfully")
		return "User added successfully", nil
	} else {
		slog.Info("User already Exists")
		return "Error adding user", errors.New("user already Exists")
	}
}

func (svc *UserService) generatePasswordHash(pass string) (string, error) {
	const salt = 14
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), salt)
	if err != nil {
		slog.Errorf("Failed to generate password hash: %v\n", err)
		return "", err
	}
	return string(hash), nil
}
