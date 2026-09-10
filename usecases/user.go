package usecases

import (
	"context"
	"errors"

	"github.com/polar-bear-cu/sgt-user-service/models"
	"github.com/polar-bear-cu/sgt-user-service/repositories"
)

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

func (u *UserUsecase) FindOrCreate(ctx context.Context, email, googleSub string) (models.User, bool, error) {
	existing, err := u.repo.FindByGoogleSub(ctx, googleSub)
	if err == nil {
		return existing, false, nil
	}
	if !errors.Is(err, repositories.ErrUserNotFound) {
		return models.User{}, false, err
	}

	created, err := u.repo.Create(ctx, models.User{Email: email, GoogleSub: googleSub})
	if err != nil {
		return models.User{}, false, err
	}
	return created, true, nil
}
