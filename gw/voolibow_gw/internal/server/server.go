package server

import (
	"log"
	"safir/libs/appconfigs"
	"safir/libs/appstates"
	"voolibow_gw/internal/client"
	"voolibow_gw/internal/controllers"
	"voolibow_gw/internal/middleware"
	"voolibow_gw/internal/routes"
	"voolibow_gw/internal/services"
	"voolibow_gw/internal/socket"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func RunServer() {
	var (
		listenAddress                   = appconfigs.String("listen-address", "Server listen address")
		userAuthenticationServerAddress = appconfigs.String("user-authentication-server-address", "User Authentication Server address")
		profileServerAddress            = appconfigs.String("profile-server-address", "Profile Server address")
	)
	// Handle configuration errors.
	if err := appconfigs.Parse(); err != nil {
		appstates.PanicMissingEnvParams(err.Error()) // Log an error if there are missing environment parameters.
	}

	var (
		userAuthenticationConnection = client.GrpcServerConnection(*userAuthenticationServerAddress)
		profileConnection            = client.GrpcServerConnection(*profileServerAddress)
	)
	var (
		middlewareHandler = middleware.NewMiddlewareHandler(userAuthenticationConnection)
	)
	var (
		userAuthenticationService = services.NewUserAuthenticationService(userAuthenticationConnection, profileConnection)
		userAccountService        = services.NewUserAccountService(userAuthenticationConnection)
		profileService            = services.NewProfileService(profileConnection)
	)

	var (
		userAuthenticationController = controllers.NewUserAuthenticationController(userAuthenticationService)
		userAccountController        = controllers.NewUserAccountService(userAccountService)
		profileController            = controllers.NewProfileControlelr(profileService)
	)

	defer userAuthenticationConnection.Close()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	apiRoutes := app.Group("/api")
	routes.UserAuthenticationRoutes(apiRoutes, userAuthenticationController)
	routes.UserAccountRoutes(apiRoutes, userAccountController, middlewareHandler)
	routes.ProfileRoutes(apiRoutes, profileController, middlewareHandler)

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws", websocket.New(socket.SocketConnection))
	log.Fatal(app.Listen(*listenAddress))
}
