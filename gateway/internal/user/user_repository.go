package user

import (
	"Booking_system/gateway/internal/config"
	"Booking_system/gateway/internal/user/pb"
	"Booking_system/gateway/model"
	"context"
	"errors"
	"github.com/gookit/slog"
	"time"
)

type UserRepository struct {
	cfg    *config.Config
	client pb.UserServiceClient
}

func NewUserRepository(cfg *config.Config, client pb.UserServiceClient) *UserRepository {
	return &UserRepository{cfg: cfg, client: client}
}

func mapUser(pbUser *pb.GetUserEmailByIdResponse) model.User {
	var user model.User

	if pbUser.FirstName == "" {
	} else {
		user.FirstName = pbUser.FirstName
	}

	if pbUser.LastName == "" {
	} else {
		user.LastName = pbUser.LastName
	}

	if pbUser.DateOfBirth == "" {
	} else {
		user.DateOfBirth = pbUser.DateOfBirth
	}

	if pbUser.Email == "" {
	} else {
		user.Email = pbUser.Email
	}

	if pbUser.Address == "" {
	} else {
		user.Address = pbUser.Address
	}

	if pbUser.Phone == "" {
	} else {
		user.Phone = pbUser.Phone
	}

	if pbUser.Role == "" {
	} else {
		user.Role = pbUser.Role
	}
	return user
}

func (repo *UserRepository) Delete(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.cfg.LongTimeout)*time.Second)
	defer cancel()

	res, err := repo.client.DeleteUser(ctx, &pb.DeleteUserRequest{
		Email: email,
	})
	if err != nil {
		slog.Errorf("Error deleting user: %v", err)
		return "", err
	}

	if res == nil && res.Message == "" {
		slog.Errorf("DeleteUser response is empty")
		return res.Message, errors.New("DeleteUser response is empty")
	}
	return res.Message, nil
}

func (repo *UserRepository) GetByEmail(email string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.cfg.LongTimeout)*time.Second)
	defer cancel()

	res, err := repo.client.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{
		Email: email,
	})
	if err != nil {
		slog.Errorf("Error getting user by email: %v", err)
		return model.User{}, err
	}
	var user model.User
	user.UserCredentials.Email = res.Email
	user.UserCredentials.Password = res.Password
	user.UUID = res.Uuid
	user.Role = res.Role

	return user, nil
}

func (repo *UserRepository) GetById(id string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.cfg.LongTimeout)*time.Second)
	defer cancel()

	res, err := repo.client.GetUserEmailById(ctx, &pb.GetUserEmailByIdRequest{
		Uuid: id,
	})

	if err != nil {
		slog.Errorf("Error getting user by id: %v", err)
		return model.User{}, err
	}

	var user model.User
	user.UserCredentials.Email = res.Email
	user.Role = res.Role
	user.FirstName = res.FirstName
	user.LastName = res.LastName
	user.Phone = res.Phone
	user.DateOfBirth = res.DateOfBirth
	user.Address = res.Address

	return user, nil
}

func (repo *UserRepository) UpdateUser(updateUser model.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.cfg.LongTimeout)*time.Second)
	defer cancel()

	res, err := repo.client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Email:       updateUser.Email,
		FirstName:   updateUser.FirstName,
		LastName:    updateUser.LastName,
		Phone:       updateUser.Phone,
		DateOfBirth: updateUser.DateOfBirth,
		Address:     updateUser.Address,
	})

	if err != nil {
		slog.Errorf("Error updating user data: %v", err)
		return "", err
	}

	if res == nil && res.Message == "" {
		slog.Errorf("UpdateUser response is empty")
		return res.Message, nil
	}
	return res.Message, nil
}

func (repo *UserRepository) GetAllUsers(page uint32, limit uint32) ([]model.User, string, error) {
	var users []model.User
	var totalNumber string

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(50)*time.Second)
	defer cancel()

	resp, err := repo.client.GetAllUsers(ctx, &pb.GetUsersRequest{
		Page:  page,
		Limit: limit,
	})

	if err != nil {
		slog.Errorf("Failed to get all users: ", err)
		return users, totalNumber, err
	}

	totalNumber = resp.GetNumberOfUsers()

	for _, pbUser := range resp.Users {
		users = append(users, mapUser(pbUser))
	}

	return users, totalNumber, nil
}

func (repo *UserRepository) AddAdmin(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.cfg.LongTimeout)*time.Second)
	defer cancel()

	res, err := repo.client.AddAdmin(ctx, &pb.AddAdminRequest{
		Email: email,
	})

	if err != nil {
		slog.Errorf("Error adding admin: %v", err)
		return "", err
	}

	if res == nil && res.Message == "" {
		slog.Errorf("AddAdmin response is empty")
		return res.Message, errors.New("AddAdmin response is empty")
	}

	return res.Message, nil
}

func (repo *UserRepository) SetUserRole(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.cfg.LongTimeout)*time.Second)
	defer cancel()

	res, err := repo.client.SetUserRole(ctx, &pb.SetUserRoleRequest{
		Email: email,
	})

	if err != nil {
		slog.Errorf("Error setting user role to user: %v", err)
		return "", err
	}

	if res == nil && res.Message == "" {
		slog.Errorf("SetUserRole response is empty")
		return res.Message, errors.New("SetUserRole response is empty")
	}
	return res.Message, nil
}

func (repo *UserRepository) CreateUser(user model.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.cfg.LongTimeout)*time.Second)
	defer cancel()

	res, err := repo.client.CreateUser(ctx, &pb.CreateUserRequest{
		Email:       user.Email,
		Password:    user.Password,
		Role:        user.Role,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DateOfBirth: user.DateOfBirth,
		Phone:       user.Phone,
		Address:     user.Address,
	})

	if err != nil {
		slog.Errorf("Error creating user: %v", err)
		return "", err
	}

	if res == nil && res.Message == "" {
		slog.Errorf("CreateUser response is empty")
		return res.Message, errors.New("CreateUser response is empty")
	}

	return res.Message, nil
}
