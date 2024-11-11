package chat

import "github.com/gofiber/fiber/v2"

func RegisterChatRoutes(r *fiber.App, chatCtrl *ChatController) {
	r.Get("/contact-list", chatCtrl.GetContactList)
	r.Get("/chat-history", chatCtrl.GetChatHistory)
}
