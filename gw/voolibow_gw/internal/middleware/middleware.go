package middleware

import (
	"context"
	"net/http"
	"time"
	pb "voolibow_gw/proto/api/user_authentication/token"
	"voolibow_gw/types"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

type (
	MiddlewareHandler interface {
		TokenAuthentication(*fiber.Ctx) error
	}
	middlewareHandler struct {
		tokenAuthenticationClient pb.TokenServiceClient
	}
)

func NewMiddlewareHandler(tokenAuthenticationServerConnection *grpc.ClientConn) MiddlewareHandler {
	return &middlewareHandler{
		tokenAuthenticationClient: pb.NewTokenServiceClient(tokenAuthenticationServerConnection),
	}
}

// TokenAuthentication is a middleware for token validation.
func (c *middlewareHandler) TokenAuthentication(ctx *fiber.Ctx) error {
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
	ctxTimeOut, cancel := context.WithTimeout(ctx.Context(), time.Second*5)
	defer cancel()
	res, err := c.tokenAuthenticationClient.VerifyToken(ctxTimeOut, &pb.VerificationRequest{
		AccessToken: token,
		Agent:       string(ctx.Context().UserAgent()),
		Role:        "User",
	})

	if err != nil {
		rerr := types.ExtractGrpcError(err)
		return ctx.Status(rerr.Code).JSON(fiber.Map{
			"message": rerr.Message,
			"success": false,
		})
	}
	ctx.Locals("role", res.Role)
	ctx.Locals("session_id", res.SessionId)
	ctx.Locals("user_id", res.UserId)

	return ctx.Next()
}
