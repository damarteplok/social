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

	PembuatanMediaBeritaTechnology interface {
		Get(context.Context, int64) (*store.PembuatanMediaBeritaTechnology, error)
		Set(context.Context, *store.PembuatanMediaBeritaTechnology) error
		Delete(context.Context, int64)
	}

	ApprovingArtikel interface {
		Get(context.Context, int64) (*store.ApprovingArtikel, error)
		Set(context.Context, *store.ApprovingArtikel) error
		Delete(context.Context, int64)
	}

	ReviewingArtikel interface {
		Get(context.Context, int64) (*store.ReviewingArtikel, error)
		Set(context.Context, *store.ReviewingArtikel) error
		Delete(context.Context, int64)
	}

	PembuatanArtikel interface {
		Get(context.Context, int64) (*store.PembuatanArtikel, error)
		Set(context.Context, *store.PembuatanArtikel) error
		Delete(context.Context, int64)
	}
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users: &UsersStore{
			rdb: rbd,
		},
		// GENERATED CACHE CODE CONSTRUCTOR

		PembuatanMediaBeritaTechnology: &PembuatanMediaBeritaTechnologyStore{
			rdb: rbd,
		},

		ApprovingArtikel: &ApprovingArtikelStore{
			rdb: rbd,
		},

		ReviewingArtikel: &ReviewingArtikelStore{
			rdb: rbd,
		},

		PembuatanArtikel: &PembuatanArtikelStore{
			rdb: rbd,
		},
	}
}
