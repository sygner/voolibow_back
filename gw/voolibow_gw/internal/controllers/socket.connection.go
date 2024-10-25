package controllers

import (
	"math/rand"
	"time"
	"voolibow_gw/internal/global"

	"github.com/gofiber/fiber/v2"
)

type (
	SocketConnectionController interface{}
	socketConnectionController struct {
	}
)

func NewSocketCOnnectionController() SocketConnectionController {
	return &socketConnectionController{}
}

func (c *socketConnectionController) KeyRequest(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(int32)

	key := RandStringRunes(40)
	global.TOKENS[key] = &global.TokenValue{UserId: userId, Username: "username", ExpireTime: time.Now().Add(time.Second * 10)}
	for id, token := range global.TOKENS {
		if token.ExpireTime.Before(time.Now()) {
			delete(global.TOKENS, id)
		}
	}

	return ctx.Status(200).JSON(fiber.Map{
		"data":    key,
		"success": true,
	})
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
