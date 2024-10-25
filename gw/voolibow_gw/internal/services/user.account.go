package services

import (
	"context"
	"voolibow_gw/internal/models"
	pb "voolibow_gw/proto/api/user_authentication/account"
	"voolibow_gw/types"

	"google.golang.org/grpc"
)

type (
	UserAccountService interface {
		Logout(string, int32) *types.Error
		KillSession(int32) *types.Error
		GetSessions(int32) ([]models.Session, *types.Error)
		GetLoginHistories(*models.GetLoginHistoriesDTO) ([]models.LoginHistory, *types.Error)
	}
	userAccountService struct {
		userAccountClient pb.AccountServiceClient
	}
)

func NewUserAccountService(userAccountServiceConnection *grpc.ClientConn) UserAccountService {
	return &userAccountService{
		userAccountClient: pb.NewAccountServiceClient(userAccountServiceConnection),
	}
}

func (c *userAccountService) Logout(accessToken string, userId int32) *types.Error {
	_, err := c.userAccountClient.Logout(context.Background(), &pb.LogoutRequest{
		AccessToken: accessToken,
		UserId:      userId,
	})

	if err != nil {
		return types.ExtractGrpcError(err)
	}
	return nil
}

func (c *userAccountService) KillSession(sessionId int32) *types.Error {
	_, err := c.userAccountClient.KillSession(context.Background(), &pb.KillSessionRequest{
		SessionId: sessionId,
	})

	if err != nil {
		return types.ExtractGrpcError(err)
	}
	return nil
}

func (c *userAccountService) GetSessions(userId int32) ([]models.Session, *types.Error) {
	res, err := c.userAccountClient.GetSessions(context.Background(), &pb.GetSessionsRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, types.ExtractGrpcError(err)
	}

	sessions := make([]models.Session, 0)
	for _, session := range res.Sessions {
		sessions = append(sessions, models.Session{
			Agent:     session.Agent,
			Ip:        session.Ip,
			Status:    session.Status,
			SessionId: session.SessionId,
			CreatedAt: session.CreatedAt,
		})
	}

	return sessions, nil
}

func (c *userAccountService) GetLoginHistories(data *models.GetLoginHistoriesDTO) ([]models.LoginHistory, *types.Error) {
	res, err := c.userAccountClient.GetLoginHistory(context.Background(), &pb.GetLoginHistories{
		UserId: data.UserId,
		Pagination: &pb.Pagination{
			Offset:   data.Pagination.Offset,
			Limit:    data.Pagination.Limit,
			GetTotal: data.Pagination.Total,
		},
	})

	if err != nil {
		return nil, types.ExtractGrpcError(err)
	}

	loginHistories := make([]models.LoginHistory, 0)
	for _, loginHistory := range res.LoginHistories {
		loginHistories = append(loginHistories, models.LoginHistory{
			Id:           loginHistory.Id,
			UserId:       loginHistory.UserId,
			UserRole:     loginHistory.UserRole,
			Section:      loginHistory.Section,
			Ip:           loginHistory.Ip,
			Agent:        loginHistory.Agent,
			Logged_in_at: loginHistory.LoggedInAt,
		})
	}
	return loginHistories, nil
}
