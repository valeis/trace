package controller

import (
	"Booking_system/security_service/internal/config"
	"Booking_system/security_service/internal/models"
	"Booking_system/security_service/internal/pb"
	"Booking_system/security_service/pkg/jwtoken"
	"context"
	"errors"
	"github.com/gookit/slog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ISecurityService interface {
	Create(user *models.User) (string, error)
	Login(user *models.UserCredentialsModel) (map[string]string, string, error)
	RefreshUserToken(token string, email string) (map[string]string, error)
}

type SecurityController struct {
	securityService ISecurityService
	cfg             *config.Config
}

func NewSecurityController(securityService ISecurityService, cfg *config.Config) *SecurityController {
	return &SecurityController{
		securityService: securityService,
		cfg:             cfg,
	}
}

func (ctrl *SecurityController) CreateUser(ctx context.Context, req *pb.UserCredentials) (*pb.ResponseMessage, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password cannot be empty")
	}
	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}
	message, err := ctrl.securityService.Create(user)
	return &pb.ResponseMessage{Message: message}, err
}

func (s *SecurityController) Login(ctx context.Context, req *pb.UserCredentials) (*pb.Tokens, error) {
	userLogin := models.UserCredentialsModel{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	token, userId, err := s.securityService.Login(&userLogin)
	if err != nil {
		slog.Errorf("Failed to login user: %v", err)
		return nil, err
	}

	return &pb.Tokens{
		AccessToken:  token["access_token"],
		RefreshToken: token["refresh_token"],
		UserId:       userId,
	}, nil
}

func (s *SecurityController) RefreshSession(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	refToken := req.Token
	userEmail, err := jwtoken.IsTokenValid(refToken)
	if err != nil {
		slog.Errorf("Failed to validate token: %v", err)
		return nil, status.Error(codes.Unauthenticated, errors.New("failed to refresh user session").Error())
	}
	tokenMap, err := s.securityService.RefreshUserToken(refToken, userEmail)
	if err != nil {
		slog.Errorf("Failed to refresh user session: %v", err)
		return nil, status.Error(codes.Internal, errors.New("failed to refresh user session").Error())
	}
	return &pb.RefreshResponse{Tokens: tokenMap}, nil
}

func (s *SecurityController) ValidateToken(ctx context.Context, req *pb.Token) (*pb.ValidateTokenResponse, error) {
	email, err := jwtoken.IsTokenValid(req.Token)
	if err != nil {
		slog.Errorf("Failed to validate token: %v", err)
		return nil, status.Error(codes.Unauthenticated, errors.New("failed to validate token").Error())
	}
	result := &pb.ValidateTokenResponse{Token: req.Token, Email: email}
	return result, nil
}
