package routes

import (
	"voolibow_gw/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserAuthenticationRoutes(router fiber.Router, controllers controllers.UserAuthenticationController) {
	authRouter := router.Group("/auth")

	authRouter.Post("/signup", controllers.Signup)
	authRouter.Post("/verify", controllers.Verify)
	authRouter.Post("/signin", controllers.Signin)
	authRouter.Post("/token/renew", controllers.RenewToken)
}
