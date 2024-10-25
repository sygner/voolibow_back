package helpers

import (
	"time"
	"voolibow_gw/internal/global"
	"voolibow_gw/internal/models"

	"github.com/gofiber/contrib/websocket"
)

func SocketValidation(c *websocket.Conn) (*global.TokenValue, *models.WebSocketResponse) {
	webSocketLogin := new(models.WebSocketLoginDTO)
	err := c.ReadJSON(&webSocketLogin)
	if err != nil {
		return nil, &models.WebSocketResponse{Event: "error", State: "signin", Data: "failed to signin"}
	}
	if global.TOKENS[webSocketLogin.Key] == nil {
		return nil, &models.WebSocketResponse{Event: "error", State: "signin", Data: "key not found"}
	}
	token := global.TOKENS[webSocketLogin.Key]
	if time.Now().After(token.ExpireTime) {
		delete(global.TOKENS, webSocketLogin.Key)
		return nil, &models.WebSocketResponse{Event: "error", State: "signin", Data: "key expired"}
	}
	return token, nil

}
