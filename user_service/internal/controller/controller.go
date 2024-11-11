package controller

import (
	"Booking_system/user_service/internal/models"
	"Booking_system/user_service/internal/pb"
	"Booking_system/user_service/internal/util"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IUserService interface {
	Delete(userEmail string) error
	GetUserByEmail(userEmail string) (*models.User, error)
	GetUserEmailById(userId string) (*models.User, error)
	UpdateUser(user *models.User) error
	GetAllUsers(pagination util.Pagination) ([]models.User, string)
	AddAdmin(userEmail string) error
	SetUserRole(userEmail string) error
	CreateUser(user *models.User) (string, error)
}

type UserController struct {
	userService IUserService
}

func NewUserController(userService IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (ctrl *UserController) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.UserResponseMessage, error) {
	userEmail := req.GetEmail()

	if userEmail == "" {
		return nil, status.Error(codes.InvalidArgument, "Email field cannot be empty")
	}
	err := ctrl.userService.Delete(userEmail)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Couldn't delete")
	}

	return &pb.UserResponseMessage{Message: "User deleted successfully"}, nil
}

func (ctrl *UserController) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	userEmail := req.GetEmail()
	user, err := ctrl.userService.GetUserByEmail(userEmail)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	userResponse := &pb.GetUserByEmailResponse{
		Email:    req.Email,
		Uuid:     user.UUID,
		Password: user.Password,
		Role:     user.Role,
	}
	return userResponse, nil
}

func (ctrl *UserController) GetUserEmailById(ctx context.Context, req *pb.GetUserEmailByIdRequest) (*pb.GetUserEmailByIdResponse, error) {
	userId := req.Uuid
	user, err := ctrl.userService.GetUserEmailById(userId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User email not found")
	}

	userResponse := &pb.GetUserEmailByIdResponse{
		Email:       user.Email,
		Role:        user.Role,
		Address:     user.Address,
		DateOfBirth: user.DateOfBirth,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Phone:       user.Phone,
	}
	return userResponse, nil
}

func (ctrl *UserController) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponseMessage, error) {
	if req.Email == "" {
		return nil, errors.New("Email cannot be empty")
	}
	user := &models.User{
		Email:       req.Email,
		Address:     req.Address,
		DateOfBirth: req.DateOfBirth,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Phone:       req.Phone,
	}
	err := ctrl.userService.UpdateUser(user)

	if err != nil {
		return &pb.UserResponseMessage{Message: "Error updating user"}, err
	}
	return &pb.UserResponseMessage{Message: "User data updated successfully"}, nil
}

func (ctrl *UserController) GetAllUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	pag := util.Pagination{
		Page:  int(req.Page),
		Limit: int(req.Limit),
	}

	users, totalNumber := ctrl.userService.GetAllUsers(pag)

	getUsersResponse := make([]*pb.GetUserEmailByIdResponse, len(users))

	for i := range getUsersResponse {
		u := users[i]
		getUsersResponse[i] = &pb.GetUserEmailByIdResponse{
			Email:       u.Email,
			Address:     u.Address,
			DateOfBirth: u.DateOfBirth,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			Phone:       u.Phone,
			Role:        u.Role,
		}
	}
	return &pb.GetUsersResponse{
		Users:         getUsersResponse,
		NumberOfUsers: totalNumber,
	}, nil
}

func (ctrl *UserController) AddAdmin(ctx context.Context, req *pb.AddAdminRequest) (*pb.UserResponseMessage, error) {
	userEmail := req.GetEmail()
	err := ctrl.userService.AddAdmin(userEmail)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Couldn't update role")
	}
	return &pb.UserResponseMessage{Message: "User role updated successfully"}, nil
}

func (ctrl *UserController) SetUserRole(ctx context.Context, req *pb.SetUserRoleRequest) (*pb.UserResponseMessage, error) {
	userEmail := req.GetEmail()
	err := ctrl.userService.SetUserRole(userEmail)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Couldn't update role")
	}
	return &pb.UserResponseMessage{Message: "User role updated successfully"}, nil
}

func (ctrl *UserController) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponseMessage, error) {
	if req.Password == "" || req.FirstName == "" || req.LastName == "" || req.Role == "" || req.Email == "" || req.Address == "" || req.DateOfBirth == "" || req.Phone == "" {
		return nil, errors.New("user data isn't completed")
	}

	user := &models.User{
		Password:    req.Password,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Role:        req.Role,
		Email:       req.Email,
		Address:     req.Address,
		DateOfBirth: req.DateOfBirth,
		Phone:       req.Phone,
	}
	message, err := ctrl.userService.CreateUser(user)
	return &pb.UserResponseMessage{Message: message}, err
}
