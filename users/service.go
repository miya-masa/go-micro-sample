package users

import (
	"context"
	"errors"
)

var (
	users = map[string]User{
		"user1": User{ID: 1, Name: "user1"},
		"user2": User{ID: 2, Name: "user2"},
		"user3": User{ID: 3, Name: "user3"},
	}

	ErrNotFound = errors.New("user not found.")
)

type UserService interface {
	UserByName(ctx context.Context, name string) (*User, error)
}

func NewService() UserService {
	return &impl{}
}

type impl struct {
}

func (u *impl) UserByName(ctx context.Context, name string) (*User, error) {

	user, ok := users[name]
	if ok {
		return &user, nil
	}

	return nil, ErrNotFound
}
