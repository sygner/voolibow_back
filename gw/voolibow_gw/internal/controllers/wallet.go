package controllers

import (
	"voolibow_gw/internal/models"
	"voolibow_gw/internal/services"

	"github.com/gofiber/fiber/v2"
)

type (
	WalletController interface {
		RegisterWallet(*fiber.Ctx) error
		GetBalance(*fiber.Ctx) error
		TransferETH(*fiber.Ctx) error
		GetWalletsByUserId(*fiber.Ctx) error
	}
	walletController struct {
		walletService services.WalletService
	}
)

func NewWalletController(walletService services.WalletService) WalletController {
	return &walletController{
		walletService: walletService,
	}
}

func (c *walletController) RegisterWallet(ctx *fiber.Ctx) error {
	registerWalletDTO := new(models.RegisterWalletDTO)
	if err := ctx.BodyParser(registerWalletDTO); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error":   "failed to get data",
			"success": false,
		})
	}
	userId := ctx.Locals("user_id").(int32)

	switch registerWalletDTO.Currency {
	case 0:
		registerWalletDTO.CurrencyFullName = "Ethereum"
	default:
		return ctx.Status(400).JSON(fiber.Map{
			"error":   "no currency found",
			"success": false,
		})
	}
	registerWalletDTO.UserId = userId
	res, err := c.walletService.RegitserWallet(registerWalletDTO)
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

func (c *walletController) GetBalance(ctx *fiber.Ctx) error {
	currencyDTO := new(models.CurrencyIdDTO)
	if err := ctx.BodyParser(currencyDTO); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error":   "failed to get data",
			"success": false,
		})
	}
	switch currencyDTO.Currency {
	case 0:
	default:
		return ctx.Status(400).JSON(fiber.Map{
			"error":   "no currency found",
			"success": false,
		})
	}
	userId := ctx.Locals("user_id").(int32)

	res, err := c.walletService.GetBalance(userId, currencyDTO.Currency)
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

func (c *walletController) TransferETH(ctx *fiber.Ctx) error {
	transferDTO := new(models.TransferETHDTO)
	if err := ctx.BodyParser(transferDTO); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error":   "failed to get data",
			"success": false,
		})
	}
	userId := ctx.Locals("user_id").(int32)
	err := c.walletService.TransferETH(&models.TransferETHDTO{
		UserId:    userId,
		ToAddress: transferDTO.ToAddress,
		Amount:    transferDTO.Amount,
		RoomId:    transferDTO.RoomId,
	})
	if err != nil {
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error":   err.Message,
			"success": false,
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"message": "transfered",
		"success": true,
	})
}

func (c *walletController) GetWalletsByUserId(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(int32)
	res, err := c.walletService.GetWalletsByUserId(userId)
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

func (c *walletController) GetTransactionsByAddress(ctx *fiber.Ctx) error {
	addressDTO := new(models.AddressDTO)
	if err := ctx.BodyParser(addressDTO); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error":   "failed to get data",
			"success": false,
		})
	}

	res, err := c.walletService.GetTransactionsByAddress(addressDTO.Address)
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
