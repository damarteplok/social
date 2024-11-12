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

	KantorNgetesId interface {
		Get(context.Context, int64) (*store.KantorNgetesId, error)
		Set(context.Context, *store.KantorNgetesId) error
		Delete(context.Context, int64)
	}


	SetujuiSomething interface {
		Get(context.Context, int64) (*store.SetujuiSomething, error)
		Set(context.Context, *store.SetujuiSomething) error
		Delete(context.Context, int64)
	}


	ReviewSomething interface {
		Get(context.Context, int64) (*store.ReviewSomething, error)
		Set(context.Context, *store.ReviewSomething) error
		Delete(context.Context, int64)
	}


	BikinSomething interface {
		Get(context.Context, int64) (*store.BikinSomething, error)
		Set(context.Context, *store.BikinSomething) error
		Delete(context.Context, int64)
	}

}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users: &UsersStore{
			rdb: rbd,
		},
		// GENERATED CACHE CODE CONSTRUCTOR

		KantorNgetesId: &KantorNgetesIdStore{
			rdb: rbd,
		},


		SetujuiSomething: &SetujuiSomethingStore{
			rdb: rbd,
		},


		ReviewSomething: &ReviewSomethingStore{
			rdb: rbd,
		},


		BikinSomething: &BikinSomethingStore{
			rdb: rbd,
		},

	}
}