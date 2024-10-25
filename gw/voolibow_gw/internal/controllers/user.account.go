package controllers

import (
	"net/http"
	"strconv"
	"voolibow_gw/internal/models"
	"voolibow_gw/internal/services"

	"github.com/gofiber/fiber/v2"
)

type (
	UserAccountController interface {
		Logout(*fiber.Ctx) error
		KillSession(*fiber.Ctx) error
		GetSessions(*fiber.Ctx) error
		GetLoginHistories(*fiber.Ctx) error
	}
	userAccountController struct {
		userAccountService services.UserAccountService
	}
)

func NewUserAccountService(userAccountService services.UserAccountService) UserAccountController {
	return &userAccountController{
		userAccountService: userAccountService,
	}
}

func (c *userAccountController) Logout(ctx *fiber.Ctx) error {
	// Extract the token from the "Authorization" header
	token := ctx.Get("Authorization")

	// Check if the token is empty or does not start with "Bearer "
	if token == "" || token[:7] != "Bearer " {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "failed to authorize",
			"success": false,
		})
	}

	// Extract the token string (remove "Bearer " prefix)
	token = token[7:]
	userId := ctx.Locals("user_id").(int32)
	err := c.userAccountService.Logout(token, userId)
	if err != nil {
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error":   err.Message,
			"success": false,
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "loged out",
		"success": true,
	})
}

func (c *userAccountController) KillSession(ctx *fiber.Ctx) error {
	sessionIdS := ctx.Params("session_id", "")
	if sessionIdS == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to get data",
			"success": false,
		})
	}

	sessionId, rerr := strconv.ParseInt(sessionIdS, 10, 32)
	if rerr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse",
			"success": false,
		})
	}

	err := c.userAccountService.KillSession(int32(sessionId))
	if err != nil {
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error":   err.Message,
			"success": false,
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"message": "session killed",
		"success": true,
	})
}

func (c *userAccountController) GetSessions(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(int32)

	res, err := c.userAccountService.GetSessions(userId)
	if err != nil {
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error":   err.Message,
			"success": false,
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"data":    res,
		"success": true,
	})
}

func (c *userAccountController) GetLoginHistories(ctx *fiber.Ctx) error {
	paginationDTO := new(models.Pagination)
	if err := ctx.BodyParser(paginationDTO); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error":   "failed to get data",
			"success": false,
		})
	}
	userId := ctx.Locals("user_id").(int32)
	loginHistory := &models.GetLoginHistoriesDTO{
		UserId:     userId,
		Pagination: *paginationDTO,
	}
	res, err := c.userAccountService.GetLoginHistories(loginHistory)
	if err != nil {
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error":   err.Message,
			"success": false,
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"data":    res,
		"success": true,
	})
}
