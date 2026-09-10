package repositories

import (
	"errors"
	"sync"

	"github.com/polar-bear-cu/sgt-user-service/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	FindByID(id string) (models.User, error)
	Update(u models.User) (models.User, error)
}

type inMemoryUser struct {
	mu    sync.Mutex
	items map[string]models.User
}

func NewInMemoryUser() UserRepository {
	r := &inMemoryUser{items: map[string]models.User{}}
	id := "11111111-1111-1111-1111-111111111111"
	r.items[id] = models.User{ID: id, Email: "mock1@gmail.com", DisplayName: "Mock 1"}
	return r
}

func (r *inMemoryUser) FindByID(id string) (models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.items[id]
	if !ok {
		return models.User{}, ErrUserNotFound
	}
	return u, nil
}

func (r *inMemoryUser) Update(u models.User) (models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[u.ID]; !ok {
		return models.User{}, ErrUserNotFound
	}
	r.items[u.ID] = u
	return u, nil
}
