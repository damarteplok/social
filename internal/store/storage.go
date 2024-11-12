package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exist")
	ErrDuplicateEmail    = errors.New("a user with that email already exist")
	ErrDuplicateUsername = errors.New("a user with that username already exist")
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

	PesenKeRestorant interface {
		Create(context.Context, *PesenKeRestorant) error
		Delete(context.Context, int64) error
		GetByID(context.Context, int64) (*PesenKeRestorant, error)
	}


	KantorNgetesId interface {
		Create(context.Context, *KantorNgetesId) error
		Delete(context.Context, int64) error
		GetByID(context.Context, int64) (*KantorNgetesId, error)
	}


	Process1hti3q2 interface {
		Create(context.Context, *Process1hti3q2) error
		Delete(context.Context, int64) error
		GetByID(context.Context, int64) (*Process1hti3q2, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:          &PostStore{db},
		Users:          &UserStore{db},
		Comments:       &CommentStore{db},
		Followers:      &FollowerStore{db},
		Roles:          &RoleStore{db},
		Process1hti3q2: &Process1hti3q2Store{db},
		// GENERATED CODE CONSTRUCTOR

		PesenKeRestorant:   &PesenKeRestorantStore{db},


		KantorNgetesId:   &KantorNgetesIdStore{db},


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