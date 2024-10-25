package controllers

import (
	"net/http"
	"voolibow_gw/internal/models"
	"voolibow_gw/internal/services"

	"github.com/gofiber/fiber/v2"
)

type (
	ProfileController interface {
		UpdateUsername(*fiber.Ctx) error
		UpdateProfile(*fiber.Ctx) error
		GetProfileBySid(*fiber.Ctx) error
		GetProfileByUsername(*fiber.Ctx) error
		GetSelfProfile(*fiber.Ctx) error
	}
	profileControlelr struct {
		profileService services.ProfileService
	}
)

func NewProfileControlelr(profileService services.ProfileService) ProfileController {
	return &profileControlelr{
		profileService: profileService,
	}
}

func (c *profileControlelr) UpdateUsername(ctx *fiber.Ctx) error {
	updateUsernameDTO := new(models.UpdateUsernameDTO)
	if err := ctx.BodyParser(updateUsernameDTO); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error":   "failed to get data",
			"success": false,
		})
	}
	userId := ctx.Locals("user_id").(int32)
	err := c.profileService.UpdateUsername(userId, updateUsernameDTO.Username)
	if err != nil {
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error":   err.Message,
			"success": false,
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"message": "username updated",
		"success": true,
	})
}

func (c *profileControlelr) UpdateProfile(ctx *fiber.Ctx) error {
	updateProfileDTO := new(models.UpdateProfileDTO)
	if err := ctx.BodyParser(updateProfileDTO); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error":   "failed to get data",
			"success": false,
		})
	}
	userId := ctx.Locals("user_id").(int32)
	err := c.profileService.UpdateAvatar(userId, updateProfileDTO)
	if err != nil {
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error":   err.Message,
			"success": false,
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"message": "profile updated",
		"success": true,
	})
}

func (c *profileControlelr) GetProfileBySid(ctx *fiber.Ctx) error {
	sid := ctx.Params("sid", "")
	if sid == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to get data",
			"success": false,
		})
	}

	res, err := c.profileService.GetProfileBySid(sid)
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

func (c *profileControlelr) GetProfileByUsername(ctx *fiber.Ctx) error {
	username := ctx.Params("username", "")
	if username == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to get data",
			"success": false,
		})
	}

	res, err := c.profileService.GetProfileByUsername(username)
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

func (c *profileControlelr) GetSelfProfile(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(int32)

	res, err := c.profileService.GetProfileByUserId(userId)
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
