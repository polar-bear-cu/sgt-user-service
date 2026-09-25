package usecases

import (
	"context"
	"errors"

	"github.com/polar-bear-cu/sgt-user-service/models"
	"github.com/polar-bear-cu/sgt-user-service/repositories"
)

var ErrForbidden = errors.New("admin only")
var ErrUnauthenticated = errors.New("unauthenticated")

type UserUsecase struct {
	repo repositories.UserRepository
}

func NewUser(repo repositories.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetByID(ctx context.Context, id string) (models.User, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *UserUsecase) UpdateProfile(ctx context.Context, id, displayName, pictureURL string) (models.User, error) {
	current, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	current.DisplayName = displayName
	current.PictureURL = pictureURL
	return u.repo.Update(ctx, current)
}

func (u *UserUsecase) FindOrCreate(
	ctx context.Context,
	email, googleSub, displayName, pictureURL string,
) (models.User, bool, error) {
	existing, err := u.repo.FindByGoogleSub(ctx, googleSub)
	if err == nil {
		return existing, false, nil
	}
	if !errors.Is(err, repositories.ErrUserNotFound) {
		return models.User{}, false, err
	}

	created, err := u.repo.Create(ctx, models.User{
		Email:       email,
		GoogleSub:   googleSub,
		DisplayName: displayName,
		PictureURL:  pictureURL,
	})
	if err != nil {
		return models.User{}, false, err
	}
	return created, true, nil
}

func (u *UserUsecase) DeleteSelf(ctx context.Context, callerID string) error {
	return u.repo.Delete(ctx, callerID)
}

func (u *UserUsecase) requireAdmin(ctx context.Context, callerID string) error {
	if callerID == "" {
		return ErrUnauthenticated
	}
	caller, err := u.repo.FindByID(ctx, callerID)
	if err != nil {
		return err
	}
	if caller.Role != models.RoleAdmin {
		return ErrForbidden
	}
	return nil
}

func (u *UserUsecase) GetAll(ctx context.Context, callerID string, limit, offset int) ([]models.User, error) {
	if err := u.requireAdmin(ctx, callerID); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 10
	}
	if offset <= 0 {
		offset = 0
	}
	return u.repo.FindAll(ctx, limit, offset)
}

func (u *UserUsecase) GetByIDAsAdmin(ctx context.Context, callerID, targetID string) (models.User, error) {
	if err := u.requireAdmin(ctx, callerID); err != nil {
		return models.User{}, err
	}
	return u.repo.FindByID(ctx, targetID)
}

func (u *UserUsecase) UpdateRole(ctx context.Context, callerID, targetID, role string) (models.User, error) {
	if err := u.requireAdmin(ctx, callerID); err != nil {
		return models.User{}, err
	}
	return u.repo.UpdateRole(ctx, targetID, role)
}

func (u *UserUsecase) DeleteUser(ctx context.Context, callerID, targetID string) error {
	if err := u.requireAdmin(ctx, callerID); err != nil {
		return err
	}
	return u.repo.Delete(ctx, targetID)
}
