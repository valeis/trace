package user

import (
	"Booking_system/gateway/internal/config"
	"Booking_system/gateway/internal/middleware"
	"Booking_system/gateway/internal/middleware/authorization"
	pbSecurity "Booking_system/gateway/internal/security/pb"
	"Booking_system/gateway/internal/user/pb"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(r *fiber.App, rbacCfg *config.PermissionsConfig, userCtrl *UserController, userClient pb.UserServiceClient, securityClient pbSecurity.SecurityServiceClient) {
	authenticateMiddleware := middleware.NewAuthenticationMiddleware(securityClient)
	authorizeMiddleware := authorization.NewAuthorizationMiddleware(userClient, rbacCfg)
	route := r.Group("/user")
	route.Delete("", authenticateMiddleware.Auth(), authorizeMiddleware.Authorize("delete", "user"), userCtrl.DeleteUser)
	route.Get("/:uid", userCtrl.GetUserById)
	route.Post("/update", authenticateMiddleware.Auth(), authorizeMiddleware.Authorize("update", "user"), userCtrl.UpdateUser)
	route.Get("/all/:page/:limit", authenticateMiddleware.Auth(), authorizeMiddleware.Authorize("update", "user"), userCtrl.GetAllUsers)
	route.Post("/admin", authenticateMiddleware.Auth(), authorizeMiddleware.Authorize("write", "user"), userCtrl.AddAdmin)
	route.Post("/user", authenticateMiddleware.Auth(), authorizeMiddleware.Authorize("write", "user"), userCtrl.SetUserRole)
	route.Post("/create", authenticateMiddleware.Auth(), authorizeMiddleware.Authorize("write", "user"), userCtrl.CreateUser)
}
