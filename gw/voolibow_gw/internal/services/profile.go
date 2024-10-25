package services

import (
	"context"
	"voolibow_gw/internal/models"
	profile_pb "voolibow_gw/proto/api/profile/profile"
	"voolibow_gw/types"

	"google.golang.org/grpc"
)

type (
	ProfileService interface {
		UpdateUsername(int32, string) *types.Error
		UpdateAvatar(int32, *models.UpdateProfileDTO) *types.Error
		GetProfileBySid(string) (*models.Profile, *types.Error)
		GetProfileByUsername(string) (*models.Profile, *types.Error)
		GetProfileByUserId(int32) (*models.Profile, *types.Error)
	}
	profileService struct {
		profileClient profile_pb.ProfileServiceClient
	}
)

func NewProfileService(profileServiceConnection *grpc.ClientConn) ProfileService {
	return &profileService{
		profileClient: profile_pb.NewProfileServiceClient(profileServiceConnection),
	}
}

func (c *profileService) UpdateUsername(userId int32, username string) *types.Error {
	_, err := c.profileClient.UpdateUsername(context.Background(), &profile_pb.UpdateUsernameRequest{
		UserId:   userId,
		Username: username,
	})

	if err != nil {
		return types.ExtractGrpcError(err)
	}

	return nil
}

func (c *profileService) UpdateAvatar(userId int32, data *models.UpdateProfileDTO) *types.Error {
	_, err := c.profileClient.UpdateAvatar(context.Background(), &profile_pb.UpdateProfileRequest{
		UserId:          userId,
		Avatar:          data.Avatar,
		DisplayUsername: data.DisplayUsername,
	})

	if err != nil {
		return types.ExtractGrpcError(err)
	}

	return nil
}

func (c *profileService) GetProfileBySid(sid string) (*models.Profile, *types.Error) {
	res, err := c.profileClient.GetProfileBySid(context.Background(), &profile_pb.GetProfileBySidRequest{
		Sid: sid,
	})

	if err != nil {
		return nil, types.ExtractGrpcError(err)
	}

	return &models.Profile{
		DisplaySid:      res.DisplaySid,
		DisplayUsername: res.DisplayUsername,
		Username:        res.Username,
		Avatar:          res.Avatar,
	}, nil
}

func (c *profileService) GetProfileByUsername(username string) (*models.Profile, *types.Error) {
	res, err := c.profileClient.GetProfileByUsername(context.Background(), &profile_pb.GetProfileByUsernameRequest{
		Username: username,
	})

	if err != nil {
		return nil, types.ExtractGrpcError(err)
	}

	return &models.Profile{
		DisplaySid:      res.DisplaySid,
		DisplayUsername: res.DisplayUsername,
		Username:        res.Username,
		Avatar:          res.Avatar,
	}, nil
}

func (c *profileService) GetProfileByUserId(userId int32) (*models.Profile, *types.Error) {
	res, err := c.profileClient.GetProfileByUserId(context.Background(), &profile_pb.GetProfileByUserIdRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, types.ExtractGrpcError(err)
	}

	return &models.Profile{
		DisplaySid:      res.DisplaySid,
		DisplayUsername: res.DisplayUsername,
		Username:        res.Username,
		Avatar:          res.Avatar,
	}, nil
}
