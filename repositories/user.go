package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/polar-bear-cu/sgt-user-service/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	FindByID(ctx context.Context, id string) (models.User, error)
	Update(ctx context.Context, u models.User) (models.User, error)
}

type UserPostgres struct {
	db *pgxpool.Pool
}

func NewUserPostgres(db *pgxpool.Pool) UserRepository {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) FindByID(ctx context.Context, id string) (models.User, error) {
	var u models.User
	err := r.db.QueryRow(ctx,
		`SELECT id, email, google_sub, display_name, picture_url FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.GoogleSub, &u.DisplayName, &u.PictureURL)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	}
	return u, err
}

func (r *UserPostgres) Update(ctx context.Context, u models.User) (models.User, error) {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET display_name = $2, picture_url = $3 WHERE id = $1`,
		u.ID, u.DisplayName, u.PictureURL,
	)
	return u, err
}
