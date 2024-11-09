package store
import (
	"context"
	"database/sql"
	"errors"
)

const (
	SetujuiSomethingID = "setujui_something"
	SetujuiSomethingName = "setujui something"
	SetujuiSomethingFormID = "setujui_form_id"
	SetujuiSomethingAssignee = "atasan_2"
	SetujuiSomethingCandidateGroup = "atasan_2"
	SetujuiSomethingCandidateUser = "atasan_2"
	SetujuiSomethingSchedule = ``

)

type SetujuiSomething struct {
    ID int64 `json:"id"`
	Name string  `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}

type SetujuiSomethingStore struct {
	db *sql.DB
}

func (s *SetujuiSomethingStore) Create(ctx context.Context, model *SetujuiSomething) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.create(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}

func (s *SetujuiSomethingStore) Delete(ctx context.Context, id int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, id); err != nil {
			return err
		}
		return nil
	})	
}

func (s *SetujuiSomethingStore) Update(ctx context.Context, model *SetujuiSomething) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.update(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}
	
func (s *SetujuiSomethingStore) create(ctx context.Context, tx *sql.Tx, model *SetujuiSomething) error {
	query := `
		INSERT INTO setujuisomething (name)
		VALUES (
			$1
		) RETURNING 
		 	id, name, 
			created_at, updated_at
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		model.Name,
	).Scan(
		&model.ID,
		&model.Name,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *SetujuiSomethingStore) GetByID(ctx context.Context, id int64) (*SetujuiSomething, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM setujuisomething
		WHERE id = $1 AND deleted_at IS NULL
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var model SetujuiSomething
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.Name,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &model, nil
}

func (s *SetujuiSomethingStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `UPDATE setujuisomething SET deleted_at = NOW() WHERE id = $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *SetujuiSomethingStore) update(ctx context.Context, tx *sql.Tx, model *SetujuiSomething) error {
	query := `
		UPDATE setujuisomething
		SET name = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING id, name, created_at updated_at;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		model.Name,
		model.ID,
	).Scan(&model.ID, &model.Name, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

