package security

import (
	"Booking_system/gateway/internal/config"
	"Booking_system/gateway/internal/security/pb"
	"Booking_system/gateway/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gookit/slog"
	"io"
	"net/http"
	"time"
)

type SecurityRepository struct {
	cfg    *config.Config
	client pb.SecurityServiceClient
}

type ChatServiceResponse struct {
	Message string `json:"message"`
}

func NewSecurityRepository(cfg *config.Config, client pb.SecurityServiceClient) *SecurityRepository {
	return &SecurityRepository{cfg: cfg, client: client}
}

func (repo *SecurityRepository) Login(loginUser model.UserCredentials) (model.Tokens, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.cfg.LongTimeout)*time.Second)
	defer cancel()

	res, err := repo.client.Login(ctx, &pb.UserCredentials{
		Email:    loginUser.Email,
		Password: loginUser.Password,
	})
	var tokens model.Tokens
	if err != nil {
		slog.Errorf("Could not login user: %v", err)
		return tokens, err
	}
	tokens.AccessToken = res.AccessToken
	tokens.RefreshToken = res.RefreshToken
	tokens.UserUUID = res.UserId
	return tokens, nil
}

func (repo *SecurityRepository) Create(createUser model.UserCredentials) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.cfg.LongTimeout)*time.Second)
	defer cancel()

	res, err := repo.client.CreateUser(ctx, &pb.UserCredentials{
		Email:    createUser.Email,
		Password: createUser.Password,
	})

	if err != nil {
		slog.Errorf("Error creating user: %v", err)
		return "", err
	}

	if res == nil && res.Message == "" {
		slog.Errorf("CreateUser response is empty")
		return res.Message, errors.New("CreateUser response is empty")
	}

	if err := RegisterUserChatService(createUser, repo); err != nil {
		slog.Warnf("Failed to notify chat service: %v", err)
	}

	return res.Message, nil
}

func (repo *SecurityRepository) Refresh(token string) (model.Tokens, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.cfg.LongTimeout)*time.Second)
	defer cancel()

	res, err := repo.client.RefreshSession(ctx, &pb.RefreshRequest{
		Token: token,
	})
	if err != nil {
		slog.Errorf("Could not refresh user session: %v", err)
		return model.Tokens{}, err
	}
	tokens := model.Tokens{
		AccessToken:  res.Tokens["access_token"],
		RefreshToken: res.Tokens["refresh_token"],
	}
	return tokens, nil
}

func RegisterUserChatService(registerUser model.UserCredentials, repo *SecurityRepository) error {
	registerUserData := model.UserData{
		Username: registerUser.Email,
		Password: registerUser.Password,
	}

	payload, err := json.Marshal(registerUserData)
	if err != nil {
		return fmt.Errorf("failed to marshal user credentials: %v", err)
	}

	resp, err := http.Post(repo.cfg.ChatServiceURL+"register", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send request to chat service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("chat service responded with status: %v", resp.StatusCode)
	}

	return nil
}

func (repo *SecurityRepository) CheckUserExistence(username model.Username) (string, error) {
	usernamePayload, err := json.Marshal(username)
	if err != nil {
		return "", fmt.Errorf("failed to marshal username: %v", err)
	}

	resp, err := http.Post(repo.cfg.ChatServiceURL+"verify-contact", "application/json", bytes.NewBuffer(usernamePayload))

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var chatResponse ChatServiceResponse
	if err := json.Unmarshal(responseBody, &chatResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal response JSON: %v", err)
	}

	if chatResponse.Message != "" {
		slog.Errorf("Chat service returned an error: %s", chatResponse.Message)
		return "", fmt.Errorf("chat service error: %s", chatResponse.Message)
	}

	return chatResponse.Message, nil
}
