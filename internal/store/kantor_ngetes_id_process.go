package store
import (
	"context"
	"database/sql"
	"errors"
)

const (
	KantorNgetesIdVersion = 1
	KantorNgetesIdProcessDefinitionKey = 2251799813765404
	KantorNgetesIdResourceName = "bikin_something.bpmn"
)

// TODO: UPDATE THIS STRUCT AND CODE BELOW
type KantorNgetesId struct {
    ID int64 `json:"id"`
	ProcessDefinitionKey int64  `json:"process_definition_key"`
	Version int32 `json:"version"`
	ResourceName string `json:"resource_name"`
	ProcessInstanceKey int64 `json:"process_instance_key"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}

type KantorNgetesIdStore struct {
	db *sql.DB
}

func (s *KantorNgetesIdStore) Create(ctx context.Context, model *KantorNgetesId) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.create(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}

func (s *KantorNgetesIdStore) Delete(ctx context.Context, id int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, id); err != nil {
			return err
		}
		return nil
	})	
}

func (s *KantorNgetesIdStore) Update(ctx context.Context, model *KantorNgetesId) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.update(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}
	
func (s *KantorNgetesIdStore) create(ctx context.Context, tx *sql.Tx, model *KantorNgetesId) error {
	//model.Version = 1
	//model.ProcessDefinitionKey = 2251799813765404
	model.ResourceName = "bikin_something.bpmn"

	query := `
		INSERT INTO kantor_ngetes_id (process_definition_key, version, resource_name, process_instance_key)
		VALUES (
			$1, 
			$2, 
			$3,
			$4
		) RETURNING 
		 	id, process_definition_key, version, resource_name, process_instance_key,
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
		mode.ProcessInstanceKey,
	).Scan(
		&model.ID,
		&model.ProcessDefinitionKey,
		&model.Version,
		&model.ResourceName,
		&model.ProcessInstanceKey,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *KantorNgetesIdStore) GetByID(ctx context.Context, id int64) (*KantorNgetesId, error) {
	query := `
		SELECT id, process_definition_key, version, resource_name, process_instance_key, created_at, updated_at
		FROM kantor_ngetes_id
		WHERE id = $1 AND deleted_at IS NULL
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var model KantorNgetesId
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.ProcessDefinitionKey,
		&model.Version,
		&model.ResourceName,
		&model.ProcessInstanceKey,
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

func (s *KantorNgetesIdStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `UPDATE kantor_ngetes_id SET deleted_at = NOW() WHERE id = $1;`

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

func (s *KantorNgetesIdStore) update(ctx context.Context, tx *sql.Tx, model *KantorNgetesId) error {
	query := `
		UPDATE kantor_ngetes_id
		SET process_definition_key = $1, version = $2, resource_name = $3, process_instance_key = $4, updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING id, process_definition_key, version, resource_name, process_instance_key, created_at updated_at;
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
		model.ID,
	).Scan(&model.ID, &model.ProcessDefinitionKey, &model.Version, &model.ResourceName, &model.ProcessInstanceKey, &model.CreatedAt, &model.UpdatedAt)
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

