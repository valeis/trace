package security

import (
	"Booking_system/gateway/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
	"log"
	"net/http"
)

type ISecurityService interface {
	Login(loginUser model.UserCredentials) (model.Tokens, error)
	Refresh(token string) (model.Tokens, error)
	Create(createUser model.UserCredentials) (string, error)
	CheckUserExistence(username model.Username) (string, error)
}

type SecurityController struct {
	service ISecurityService
}

func NewSecurityController(service ISecurityService) *SecurityController {
	return &SecurityController{service: service}
}

func (ctrl *SecurityController) Login(ctx *fiber.Ctx) error {
	var user model.UserCredentials
	err := ctx.BodyParser(&user)
	if err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tokens, err := ctrl.service.Login(user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("Login request successful")
	return ctx.JSON(fiber.Map{
		"access-token":  tokens.AccessToken,
		"refresh-token": tokens.RefreshToken,
		"userId":        tokens.UserUUID,
		"status":        "true",
	})
}

func (ctrl *SecurityController) Refresh(ctx *fiber.Ctx) error {
	type refreshToken struct {
		Token string `json:"refreshToken"`
	}
	var rt refreshToken
	err := ctx.BodyParser(&rt)
	if err != nil {
		slog.Errorf("Invalid request format: %v", err)

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	tokens, err := ctrl.service.Refresh(rt.Token)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	slog.Info("Refresh request successful")
	return ctx.JSON(fiber.Map{
		"access-token":  tokens.AccessToken,
		"refresh-token": tokens.RefreshToken,
	})
}

func (ctrl *SecurityController) CreateUser(ctx *fiber.Ctx) error {
	var user model.UserCredentials

	err := ctx.BodyParser(&user)
	if err != nil {
		log.Printf("Invalid request format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = ctrl.service.Create(user)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("Create-User request successful")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "User created successfully", "status": "true"})
}

func (ctrl *SecurityController) CheckUserExistence(ctx *fiber.Ctx) error {
	var username model.Username

	err := ctx.BodyParser(&username)
	if err != nil {
		log.Printf("Invalid request format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	_, err = ctrl.service.CheckUserExistence(username)
	if err != nil {
		return ctx.Status(http.StatusNoContent).JSON(fiber.Map{"error": err.Error()})
	}
	log.Printf("Check-User existence request successful")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Check user existence successfully", "status": "true"})
}
