package store

import (
	"context"
	"database/sql"
	"errors"
)

const (
	PesenKeRestorantVersion              = 3
	PesenKeRestorantProcessDefinitionKey = 2251799813928651
	PesenKeRestorantResourceName         = "decide-dinner.bpmn"
)

// TODO: UPDATE THIS STRUCT AND CODE BELOW
type PesenKeRestorant struct {
	ID                   int64   `json:"id"`
	ProcessDefinitionKey int64   `json:"process_definition_key"`
	Version              int32   `json:"version"`
	ResourceName         string  `json:"resource_name"`
	ProcessInstanceKey   int64   `json:"process_instance_key"`
	CreatedBy            int64   `json:"created_by"`
	UpdatedBy            *int64  `json:"updated_by"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
	DeletedAt            *string `json:"deleted_at"`
}

type PesenKeRestorantStore struct {
	db *sql.DB
}

func (s *PesenKeRestorantStore) Create(ctx context.Context, model *PesenKeRestorant) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.create(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}

func (s *PesenKeRestorantStore) Delete(ctx context.Context, id int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, id); err != nil {
			return err
		}
		return nil
	})
}

func (s *PesenKeRestorantStore) Update(ctx context.Context, model *PesenKeRestorant) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.update(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}

func (s *PesenKeRestorantStore) create(ctx context.Context, tx *sql.Tx, model *PesenKeRestorant) error {
	// model.Version = 3
	// model.ProcessDefinitionKey = 2251799813928651
	model.ResourceName = "decide-dinner.bpmn"

	query := `
		INSERT INTO pesen_ke_restorant (process_definition_key, version, resource_name, process_instance_key, created_by)
		VALUES (
			$1, 
			$2, 
			$3,
			$4,
			$5
		) RETURNING 
		 	id, process_definition_key, version, resource_name, process_instance_key, created_by, updated_by,
			created_at, updated_at
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		model.ProcessDefinitionKey,
		model.Version,
		model.ResourceName,
		model.ProcessInstanceKey,
		model.CreatedBy,
	).Scan(
		&model.ID,
		&model.ProcessDefinitionKey,
		&model.Version,
		&model.ResourceName,
		&model.ProcessInstanceKey,
		&model.CreatedBy,
		&model.UpdatedBy,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PesenKeRestorantStore) GetByID(ctx context.Context, id int64) (*PesenKeRestorant, error) {
	query := `
		SELECT id, process_definition_key, version, resource_name, process_instance_key, created_by, updated_by, created_at, updated_at
		FROM pesen_ke_restorant
		WHERE id = $1 AND deleted_at IS NULL
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var model PesenKeRestorant
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.ProcessDefinitionKey,
		&model.Version,
		&model.ResourceName,
		&model.ProcessInstanceKey,
		&model.CreatedBy,
		&model.UpdatedBy,
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

func (s *PesenKeRestorantStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `UPDATE pesen_ke_restorant SET deleted_at = NOW() WHERE id = $1;`

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

func (s *PesenKeRestorantStore) update(ctx context.Context, tx *sql.Tx, model *PesenKeRestorant) error {
	query := `
		UPDATE pesen_ke_restorant
		SET process_definition_key = $1, version = $2, resource_name = $3, process_instance_key = $4, updated_by = $5, updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING id, process_definition_key, version, resource_name, process_instance_key, created_by, updated_by, created_at updated_at;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		model.ProcessDefinitionKey,
		model.Version,
		model.ResourceName,
		model.ProcessInstanceKey,
		model.UpdatedBy,
		model.ID,
	).Scan(&model.ID,
		&model.ProcessDefinitionKey,
		&model.Version,
		&model.ResourceName,
		&model.ProcessInstanceKey,
		&model.CreatedBy,
		&model.UpdatedBy,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
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
