package grpc

import (
	"context"

	userv1 "github.com/polar-bear-cu/sgt-proto/gen/go/user/v1"
	"github.com/polar-bear-cu/sgt-user-service/usecases"
)

type UserServer struct {
	userv1.UnimplementedUserServiceServer
	uc *usecases.UserUsecase
}

func NewUserServer(uc *usecases.UserUsecase) *UserServer {
	return &UserServer{uc: uc}
}

func (s *UserServer) FindOrCreateUser(
	ctx context.Context,
	req *userv1.FindOrCreateUserRequest,
) (*userv1.FindOrCreateUserResponse, error) {
	user, created, err := s.uc.FindOrCreate(ctx, req.GetEmail(), req.GetGoogleSub())
	if err != nil {
		return nil, err
	}
	return &userv1.FindOrCreateUserResponse{
		User:    &userv1.User{Id: user.ID, Email: user.Email},
		Created: created,
	}, nil
}

func (s *UserServer) UpdateProfile(
	ctx context.Context,
	req *userv1.UpdateProfileRequest,
) (*userv1.UpdateProfileResponse, error) {
	user, err := s.uc.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &userv1.UpdateProfileResponse{
		User: &userv1.User{Id: user.ID, Email: user.Email},
	}, nil
}
