package main

import (
	"Booking_system/gateway/internal/chat"
	"Booking_system/gateway/internal/config"
	"Booking_system/gateway/internal/middleware"
	"Booking_system/gateway/internal/security"
	"Booking_system/gateway/internal/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

func main() {
	r := fiber.New()
	r.Use(middleware.RateLimiterMiddleware())
	r.Use(middleware.CorsMiddleware())
	r.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*") // Allow requests from any origin
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		return c.Next()
	})
	registerRoutes(r)
	err := r.Listen(":1337")
	if err != nil {
		slog.Fatalf("Failed to start server: %v", err)
	}
}

func registerRoutes(r *fiber.App) {
	cfg := config.LoadConfig()
	rbacCfg := config.LoadConfigRBAC()
	securityClient, err := security.InitAuthServiceClient(cfg)
	userClient := user.InitUserServiceClient(cfg)

	if err != nil {
		slog.Fatalf("Failed to connect to security service grpc: %v", err)
	}

	securityRepo := security.NewSecurityRepository(cfg, securityClient)
	securitySvc := security.NewSecurityService(securityRepo)
	securityCtrl := security.NewSecurityController(securitySvc)

	security.RegisterSecurityRoutes(r, securityCtrl)

	chatRepo := chat.NewChatRepository(cfg)
	chatSvc := chat.NewChatService(chatRepo)
	chatCtrl := chat.NewChatController(chatSvc)

	chat.RegisterChatRoutes(r, chatCtrl)

	userRepo := user.NewUserRepository(cfg, userClient)
	userSvc := user.NewUserService(userRepo)
	userCtrl := user.NewUserController(userSvc)

	user.RegisterUserRoutes(r, rbacCfg, userCtrl, userClient, securityClient)
}
