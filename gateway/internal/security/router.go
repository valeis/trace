package security

import "github.com/gofiber/fiber/v2"

func RegisterSecurityRoutes(r *fiber.App, securityCtrl *SecurityController) {
	r.Post("/login", securityCtrl.Login)
	r.Post("/refresh", securityCtrl.Refresh)
	r.Post("/register", securityCtrl.CreateUser)
	r.Post("/verify-contact", securityCtrl.CheckUserExistence)
}
