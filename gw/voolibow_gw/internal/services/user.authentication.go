package services

import (
	"context"
	"fmt"
	"voolibow_gw/internal/models"
	profile_pb "voolibow_gw/proto/api/profile/profile"
	pb "voolibow_gw/proto/api/user_authentication/authentication"
	token_pb "voolibow_gw/proto/api/user_authentication/token"
	"voolibow_gw/types"

	"google.golang.org/grpc"
)

type (
	UserAuthenticationService interface {
		Signup(phoneNumber, agent, ip string) *types.Error
		Verify(*models.VerificationDTO) (*models.TokenInfo, *types.Error)
		Signin(phoneNumber, agent, ip string) *types.Error
		RenewToken(data *models.RenewTokenDTO) (*models.TokenInfo, *types.Error)
	}
	userAuthenticationService struct {
		userAuthenticationClient pb.AuthenticationServiceClient
		userTokenClient          token_pb.TokenServiceClient
		profileClient            profile_pb.ProfileServiceClient
	}
)

func NewUserAuthenticationService(userAuthenticationServiceConnection *grpc.ClientConn, profileServiceConnection *grpc.ClientConn) UserAuthenticationService {
	return &userAuthenticationService{
		userAuthenticationClient: pb.NewAuthenticationServiceClient(userAuthenticationServiceConnection),
		userTokenClient:          token_pb.NewTokenServiceClient(userAuthenticationServiceConnection),
		profileClient:            profile_pb.NewProfileServiceClient(profileServiceConnection),
	}
}

func (c *userAuthenticationService) Signup(phoneNumber, agent, ip string) *types.Error {
	_, err := c.userAuthenticationClient.Signup(context.Background(), &pb.SignupRequest{
		Phone: phoneNumber,
		Agent: agent,
		Ip:    ip,
	})
	if err != nil {
		return types.ExtractGrpcError(err)
	}
	return nil
}

func (c *userAuthenticationService) Signin(phoneNumber, agent, ip string) *types.Error {
	_, err := c.userAuthenticationClient.Signin(context.Background(), &pb.SigninRequest{
		Phone: phoneNumber,
		Agent: agent,
		Ip:    ip,
	})
	if err != nil {
		return types.ExtractGrpcError(err)
	}
	return nil
}

func (c *userAuthenticationService) Verify(data *models.VerificationDTO) (*models.TokenInfo, *types.Error) {
	res, err := c.userAuthenticationClient.Verify(context.Background(), &pb.VerificationRequest{
		VerificationMethod: pb.VerificationMethod(data.VerificationMethod),
		Code:               fmt.Sprintf("%d", data.Code),
		Agent:              data.Agent,
		Ip:                 data.Ip,
	})

	if err != nil {
		return nil, types.ExtractGrpcError(err)
	}
	if _, err := c.profileClient.GetProfileByUserId(context.Background(), &profile_pb.GetProfileByUserIdRequest{
		UserId: res.UserId,
	}); err != nil {
		rerr := types.ExtractGrpcError(err)
		if rerr.Code == 404 {
			_, err := c.profileClient.AddProfile(context.Background(), &profile_pb.AddProfileRequest{UserId: res.UserId})
			if err != nil {
				return nil, types.ExtractGrpcError(err)
			}
		} else {
			return nil, types.NewInternalError("failed to fetch profile")
		}
	}
	return &models.TokenInfo{
		AccessToken:  res.TokenInfo.AccessToken,
		RefreshToken: res.TokenInfo.RefreshToken,
		ExpireAt:     int32(res.TokenInfo.Expiry),
	}, nil
}

func (c *userAuthenticationService) RenewToken(data *models.RenewTokenDTO) (*models.TokenInfo, *types.Error) {
	res, err := c.userTokenClient.RenewToken(context.Background(), &token_pb.RenewTokenRequest{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		Agent:        data.Agent,
		Ip:           data.Ip,
	})

	if err != nil {
		return nil, types.ExtractGrpcError(err)
	}

	return &models.TokenInfo{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpireAt:     res.Expiry,
	}, nil
}
