package controllers

import (
	"voolibow_gw/internal/models"
	"voolibow_gw/internal/services"

	"github.com/gofiber/fiber/v2"
)

type (
	UserAuthenticationController interface {
		Signup(*fiber.Ctx) error
		Verify(*fiber.Ctx) error
		Signin(*fiber.Ctx) error
		RenewToken(*fiber.Ctx) error
	}
	userAuthenticationController struct {
		userAuthenticationService services.UserAuthenticationService
	}
)

func NewUserAuthenticationController(userAuthenticationService services.UserAuthenticationService) UserAuthenticationController {
	return &userAuthenticationController{
		userAuthenticationService: userAuthenticationService,
	}
}

func (c *userAuthenticationController) Signup(ctx *fiber.Ctx) error {
	signupDTO := new(models.PhoneNumberDTO)
	if err := ctx.BodyParser(signupDTO); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error":   "failed to get data",
			"success": false,
		})
	}
	err := c.userAuthenticationService.Signup(signupDTO.PhoneNumber, string(ctx.Context().UserAgent()), ctx.IP())
	if err != nil {
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error":   err.Message,
			"success": false,
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "code has been sent",
		"success": true,
	})
}

func (c *userAuthenticationController) Verify(ctx *fiber.Ctx) error {
	verifyDTO := new(models.VerificationDTO)
	if err := ctx.BodyParser(verifyDTO); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error":   "failed to get data",
			"success": false,
		})
	}
	verifyDTO.Agent = string(ctx.Context().UserAgent())
	verifyDTO.Ip = ctx.IP()
	res, err := c.userAuthenticationService.Verify(verifyDTO)
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

func (c *userAuthenticationController) Signin(ctx *fiber.Ctx) error {
	signinDTO := new(models.PhoneNumberDTO)
	if err := ctx.BodyParser(signinDTO); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error":   "failed to get data",
			"success": false,
		})
	}
	err := c.userAuthenticationService.Signin(signinDTO.PhoneNumber, string(ctx.Context().UserAgent()), ctx.IP())
	if err != nil {
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error":   err.Message,
			"success": false,
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "code has been sent",
		"success": true,
	})
}

func (c *userAuthenticationController) RenewToken(ctx *fiber.Ctx) error {
	renewTokenDTO := new(models.RenewTokenDTO)
	if err := ctx.BodyParser(renewTokenDTO); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error":   "failed to get data",
			"success": false,
		})
	}
	renewTokenDTO.Agent = string(ctx.Context().UserAgent())
	renewTokenDTO.Ip = ctx.IP()

	res, err := c.userAuthenticationService.RenewToken(renewTokenDTO)
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
