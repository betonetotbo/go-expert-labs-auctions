package user_entity

import (
	"context"
	"github.com/betonetotbo/go-expert-labs-auctions/internal/internal_error"
)

type User struct {
	Id   string
	Name string
}

type UserRepositoryInterface interface {
	FindUserById(
		ctx context.Context, userId string) (*User, *internal_error.InternalError)
}
