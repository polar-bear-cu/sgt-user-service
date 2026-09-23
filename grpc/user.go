package grpc

import (
	"context"

	userv1 "github.com/polar-bear-cu/sgt-proto/gen/go/user/v1"
	"github.com/polar-bear-cu/sgt-user-service/middlewares"
	"github.com/polar-bear-cu/sgt-user-service/models"
	"github.com/polar-bear-cu/sgt-user-service/usecases"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	user, created, err := s.uc.FindOrCreate(ctx, req.GetEmail(), req.GetGoogleSub(), req.GetName(), req.GetPictureUrl())
	if err != nil {
		return nil, err
	}
	return &userv1.FindOrCreateUserResponse{
		User:    toProto(user),
		Created: created,
	}, nil
}

func (s *UserServer) UpdateProfile(
	ctx context.Context,
	req *userv1.UpdateProfileRequest,
) (*userv1.UpdateProfileResponse, error) {
	if callerID, ok := middlewares.UserIDFromContext(ctx); ok && callerID != req.GetId() {
		return nil, status.Error(codes.PermissionDenied, "cannot update another user's profile")
	}

	user, err := s.uc.UpdateProfile(ctx, req.GetId(), req.GetDisplayName(), req.GetPictureUrl())
	if err != nil {
		return nil, err
	}
	return &userv1.UpdateProfileResponse{
		User: toProto(user),
	}, nil
}

func (s *UserServer) GetUser(
	ctx context.Context,
	req *userv1.GetUserRequest,
) (*userv1.GetUserResponse, error) {
	user, err := s.userByIDForCaller(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &userv1.GetUserResponse{
		User: toProto(user),
	}, nil
}

func (s *UserServer) userByIDForCaller(ctx context.Context, targetID string) (models.User, error) {
	callerID, ok := middlewares.UserIDFromContext(ctx)
	if !ok || callerID == targetID {
		return s.uc.GetByID(ctx, targetID)
	}
	return s.uc.GetByIDAsAdmin(ctx, callerID, targetID)
}

func (s *UserServer) DeleteUser(
	ctx context.Context,
	req *userv1.DeleteUserRequest,
) (*userv1.DeleteUserResponse, error) {
	callerID, ok := middlewares.UserIDFromContext(ctx)

	var err error
	if !ok || callerID == req.GetId() {
		err = s.uc.DeleteSelf(ctx, req.GetId())
	} else {
		err = s.uc.DeleteUser(ctx, callerID, req.GetId())
	}
	if err != nil {
		return nil, err
	}
	return &userv1.DeleteUserResponse{Success: true}, nil
}

func toProto(u models.User) *userv1.User {
	return &userv1.User{
		Id:         u.ID,
		Email:      u.Email,
		Name:       u.DisplayName,
		PictureUrl: u.PictureURL,
	}
}
