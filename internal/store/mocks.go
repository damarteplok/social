package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{
			users: []User{},
		},
	}
}

type MockUserStore struct {
	users []User
}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, u *User) error {
	return nil
}

func (m *MockUserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	return &User{}, nil
}

func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	return &User{}, nil
}

func (m *MockUserStore) GetByEmailAndPassword(ctx context.Context, email, password string) (*User, error) {
	return &User{}, nil
}

func (m *MockUserStore) GetUserAll(ctx context.Context, fq PaginatedFeedQuery) ([]User, error) {
	return []User{}, nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return nil
}

func (m *MockUserStore) Activate(ctx context.Context, token string) error {
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, userID int64) error {
	return nil
}
