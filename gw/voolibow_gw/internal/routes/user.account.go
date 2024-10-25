package routes

import (
	"voolibow_gw/internal/controllers"
	"voolibow_gw/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserAccountRoutes(router fiber.Router, controllers controllers.UserAccountController, middleware middleware.MiddlewareHandler) {
	accountRouter := router.Group("/account")
	accountRouter.Use(middleware.TokenAuthentication)
	accountRouter.Get("/logout", controllers.Logout)
	accountRouter.Delete("/session/kill/:session_id", controllers.KillSession)
	accountRouter.Get("/session/get", controllers.GetSessions)
	accountRouter.Post("/auth/login/history", controllers.GetLoginHistories)
}
