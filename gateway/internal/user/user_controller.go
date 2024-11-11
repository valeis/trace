package user

import (
	"Booking_system/gateway/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
	"log"
	"net/http"
	"strconv"
)

type IUserService interface {
	Delete(email string) (string, error)
	GetByEmail(email string) (model.User, error)
	GetById(id string) (model.User, error)
	UpdateUser(updateUser model.User) (string, error)
	GetAllUsers(page uint32, limit uint32) ([]model.User, string, error)
	AddAdmin(email string) (string, error)
	SetUserRole(email string) (string, error)
	CreateUser(user model.User) (string, error)
}

type UserController struct {
	service IUserService
}

type UserResponse struct {
	Users       []model.User `json:"users"`
	TotalNumber string       `json:"totalNumber"`
}

func NewUserController(service IUserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	var request struct {
		Email string `json:"email"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	_, err := c.service.Delete(request.Email)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Infof("DeleteUser request successful")
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}

func (c *UserController) GetUserByEmail(ctx *fiber.Ctx) error {
	var request struct {
		Email string `json:"email"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	res, err := c.service.GetByEmail(request.Email)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("GetUserByEmail request successful")
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *UserController) GetUserById(ctx *fiber.Ctx) error {
	pid := ctx.Params("uid")

	res, err := c.service.GetById(pid)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	slog.Info("GetUserById request successful")
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	var user model.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	_, err := c.service.UpdateUser(user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User updated successfully"})
}

func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	pageStr := ctx.Params("page")
	limitStr := ctx.Params("limit")

	page, err := strconv.ParseUint(pageStr, 10, 32)
	if err != nil {
		slog.Errorf("Failed to get the page: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	limit, err := strconv.ParseUint(limitStr, 10, 32)
	if err != nil {
		slog.Errorf("Failed to get the limit: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	users, totalNumber, err := c.service.GetAllUsers(uint32(page), uint32(limit))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	userResponse := UserResponse{
		Users:       users,
		TotalNumber: totalNumber,
	}

	slog.Infof("All users retrieved successfully")
	return ctx.Status(fiber.StatusOK).JSON(userResponse)
}

func (c *UserController) AddAdmin(ctx *fiber.Ctx) error {
	var request struct {
		Email string `json:"email"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	_, err := c.service.AddAdmin(request.Email)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	slog.Infof("AddAdmin request successful")
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Admin added successfully"})
}

func (c *UserController) SetUserRole(ctx *fiber.Ctx) error {
	var request struct {
		Email string `json:"email"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		slog.Errorf("Invalid request format: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	_, err := c.service.SetUserRole(request.Email)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	slog.Infof("SetUserRole request successful")
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User role set successfully"})
}

func (ctrl *UserController) CreateUser(ctx *fiber.Ctx) error {
	var user model.User

	err := ctx.BodyParser(&user)
	if err != nil {
		log.Printf("Invalid request format")
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = ctrl.service.CreateUser(user)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("Create-User request successful")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "User created successfully"})
}
