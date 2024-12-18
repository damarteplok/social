package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrMethodNotAllowed  = errors.New("method not allowed")
	ErrNotFound          = errors.New("resource not found")
	ErrBadRequest        = errors.New("bad request")
	ErrConflict          = errors.New("resource already exist")
	ErrDuplicateEmail    = errors.New("a user with that email already exist")
	ErrDuplicateUsername = errors.New("a user with that username already exist")
	ErrTypeNotAllowed    = errors.New("file extension not allowed")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		GetByID(context.Context, int64) (*Post, error)
		Create(context.Context, *Post) error
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		GetUserAll(context.Context, PaginatedFeedQuery) ([]User, error)
		GetByID(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		GetByEmailAndPassword(context.Context, string, string) (*User, error)
		Create(context.Context, *sql.Tx, *User) error
		CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
	}
	Followers interface {
		Follow(ctx context.Context, followerID int64, userID int64) error
		Unfollow(ctx context.Context, followerID int64, userID int64) error
	}
	Roles interface {
		GetByName(context.Context, string) (*Role, error)
	}
	// GENERATED CODE INTERFACE

	PembuatanMediaBeritaTechnology interface {
		Create(context.Context, *PembuatanMediaBeritaTechnology) error
		Delete(context.Context, int64) error
		GetByID(context.Context, int64) (*PembuatanMediaBeritaTechnology, error)
		Update(context.Context, *PembuatanMediaBeritaTechnology) error
		Search(context.Context, PaginatedQuery) (map[string]interface{}, error)
	}

	ApprovingArtikel interface {
		Create(context.Context, *ApprovingArtikel) error
		Delete(context.Context, int64) error
		GetByID(context.Context, int64) (*ApprovingArtikel, error)
	}

	ReviewingArtikel interface {
		Create(context.Context, *ReviewingArtikel) error
		Delete(context.Context, int64) error
		GetByID(context.Context, int64) (*ReviewingArtikel, error)
	}

	PembuatanArtikel interface {
		Create(context.Context, *PembuatanArtikel) error
		Delete(context.Context, int64) error
		GetByID(context.Context, int64) (*PembuatanArtikel, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
		Roles:     &RoleStore{db},
		// GENERATED CODE CONSTRUCTOR

		PembuatanMediaBeritaTechnology: &PembuatanMediaBeritaTechnologyStore{db},

		ApprovingArtikel: &ApprovingArtikelStore{db},

		ReviewingArtikel: &ReviewingArtikelStore{db},

		PembuatanArtikel: &PembuatanArtikelStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
