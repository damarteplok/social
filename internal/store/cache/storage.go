package cache

import (
	"context"

	"github.com/damarteplok/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
		Delete(context.Context, int64)
	}
	// GENERATED CACHE CODE INTERFACE
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users: &UsersStore{
			rdb: rbd,
		},
		// GENERATED CACHE CODE CONSTRUCTOR
	}
}
