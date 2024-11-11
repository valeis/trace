package chat

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
)

type IChatService interface {
	GetContactList(username string) (ChatServiceResponse, error)
	GetChatHistory(participants ChatParticipants) (ChatHistory, error)
}

type ChatController struct {
	service IChatService
}

func NewChatController(service IChatService) *ChatController {
	return &ChatController{service: service}
}

func (ctrl *ChatController) GetContactList(ctx *fiber.Ctx) error {
	username := ctx.Query("username")

	data, err := ctrl.service.GetContactList(username)

	if err != nil {
		return ctx.Status(http.StatusNoContent).JSON(fiber.Map{"error": err.Error()})
	}
	log.Printf("Get-ContactList request successful")
	return ctx.JSON(fiber.Map{
		"data": data.Data,
	})
}

func (ctrl *ChatController) GetChatHistory(ctx *fiber.Ctx) error {
	u1 := ctx.Query("u1")
	u2 := ctx.Query("u2")

	participants := ChatParticipants{
		u1, u2,
	}

	data, err := ctrl.service.GetChatHistory(participants)
	log.Printf("Error %v", err)
	if err != nil {
		return ctx.Status(http.StatusNoContent).JSON(fiber.Map{"error": err.Error()})
	}
	log.Printf("Get-ChatHistory request successful")
	return ctx.JSON(fiber.Map{
		"data": data.Data,
	})
}
