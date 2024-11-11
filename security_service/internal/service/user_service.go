package service

import (
	"Booking_system/security_service/internal/config"
	"Booking_system/security_service/internal/models"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gookit/slog"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type IUserRepository interface {
	Create(user *models.User) error
	ValidateUserExistence(userEmail string) (*models.User, error)
	CheckIfEmailExists(mail string) bool
	GetUserByEmail(email string) (models.User, error)
}

type IRedisRepository interface {
	ReplaceToken(currentToken, newToken string, expires time.Duration) error
	InsertUserToken(key string, value string, expires time.Duration) error
}

type SecurityService struct {
	userRepo  IUserRepository
	redisRepo IRedisRepository
}

func NewSecurityService(userRepo IUserRepository, redisRepo IRedisRepository) *SecurityService {
	return &SecurityService{
		userRepo:  userRepo,
		redisRepo: redisRepo,
	}
}

func (svc *SecurityService) Create(user *models.User) (string, error) {
	existingUser, err := svc.userRepo.ValidateUserExistence(user.Email)

	if err != nil && existingUser == nil {
		slog.Errorf("ERR validating user existence: %v\n", err.Error())
		return "Error adding user", err
	} else if existingUser == nil && err == nil {
		user.Role = "user"
		hashedPassword, err := svc.generatePasswordHash(user.Password)
		if err != nil {
			slog.Errorf("Can't register new user: %v\n", err.Error())
			return "Error adding user", errors.New("can't register")
		}
		user.Password = hashedPassword
		user.UUID = uuid.New().String()
		err = svc.userRepo.Create(user)
		if err != nil {
			slog.Errorf("user failed to insert in database: %v\n", err.Error())
			return "Error adding user", err
		}
		slog.Info("User added successfully")
		return "User added successfully", nil
	} else {
		slog.Info("User already Exists")
		return "Error adding user", errors.New("user already Exists")
	}
}

func (svc *SecurityService) Login(userLogin *models.UserCredentialsModel) (map[string]string, string, error) {
	user, err := svc.userRepo.GetUserByEmail(userLogin.Email)
	if err != nil {
		slog.Errorf("Invalid credentials: %v", err)
		return nil, "", errors.New("Invalid credentials")
	}
	if err = svc.comparePasswordHash(user.Password, userLogin.Password); err != nil {
		return nil, "", errors.New("invalid credentials")
	}
	token, err := generateTokenPair(user.Email)
	if err != nil {
		slog.Errorf("Could not generate token pair %v", err)
		return nil, "", errors.New("Could not generate token pair")
	}

	if err = svc.redisRepo.InsertUserToken(token["refresh_token"], user.Email, time.Hour*5); err != nil {
		return nil, "", errors.New("could not insert refresh token")
	}
	return token, user.UUID, nil
}

func (svc *SecurityService) RefreshUserToken(token string, email string) (map[string]string, error) {
	tokenMap, err := generateTokenPair(email)
	if err != nil {
		return nil, err
	}
	if err := svc.redisRepo.ReplaceToken(token, tokenMap["refresh_token"], time.Hour*5); err != nil {
		return nil, err
	}
	return tokenMap, nil
}

func (svc *SecurityService) generatePasswordHash(pass string) (string, error) {
	const salt = 14
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), salt)
	if err != nil {
		slog.Errorf("Failed to generate password hash: %v\n", err)
		return "", err
	}
	return string(hash), nil
}

func (svc *SecurityService) comparePasswordHash(hash, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		slog.Errorf("Failed to compare password and hash: %v\n", err)
		return err
	}
	return nil
}

func generateTokenPair(email string) (map[string]string, error) {
	keyConfig := config.LoadConfig()

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userEmail"] = email
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	t, err := token.SignedString([]byte(keyConfig.Token.TKey))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	if err != nil {
		return nil, err
	}

	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["userEmail"] = email
	rtClaims["exp"] = time.Now().Add(time.Hour * 6).Unix()

	rt, err := refreshToken.SignedString([]byte(keyConfig.Token.RTKey))

	if err != nil {
		return nil, err
	}
	return map[string]string{
		"access_token":  t,
		"refresh_token": rt,
	}, nil
}
