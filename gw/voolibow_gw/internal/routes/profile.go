package routes

import (
	"voolibow_gw/internal/controllers"
	"voolibow_gw/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProfileRoutes(router fiber.Router, controllers controllers.ProfileController, middleware middleware.MiddlewareHandler) {
	profileRouter := router.Group("/profile")
	profileRouter.Use(middleware.TokenAuthentication)
	profileRouter.Post("/update/username", controllers.UpdateUsername)
	profileRouter.Post("/update/profile", controllers.UpdateProfile)
	profileRouter.Get("/get/s/:sid", controllers.GetProfileBySid)
	profileRouter.Get("/get/u/:username", controllers.GetProfileByUsername)
	profileRouter.Get("/get", controllers.GetSelfProfile)
}
