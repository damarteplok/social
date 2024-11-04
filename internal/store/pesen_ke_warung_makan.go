package store
import (
	"context"
	"database/sql"
	"errors"
)

// TODO: UPDATE THIS STRUCT AND CODE BELOW
type PesenKeWarungMakan struct {
    ID int64 `json:"id"`
	ProcessDefinitionKey int64  `json:"process_definition_key"`
	Version int32 `json:"version"`
	ResourceName string `json:"resource_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PesenKeWarungMakanStore struct {
	db *sql.DB
}
	
func (s *PesenKeWarungMakanStore) Create(ctx context.Context, model *PesenKeWarungMakan) error {
	model.Version = 1
	model.ProcessDefinitionKey = 2251799815414068
	model.ResourceName = "golang_ngetest.bpmn"

	query := `
		INSERT INTO pesen_ke_warung_makan (process_definition_key, version, resource_name)
		VALUES (
			$1, 
			$2, 
			$3
		) RETURNING 
		 	id, process_definition_key, version, resource_name, 
			created_at, updated_at
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		model.ProcessDefinitionKey,
		model.Version,
		model.ResourceName,
	).Scan(
		&model.ID,
		&model.ProcessDefinitionKey,
		&model.Version,
		&model.ResourceName,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PesenKeWarungMakanStore) GetByID(ctx context.Context, id int64) (*PesenKeWarungMakan, error) {
	query := `
		SELECT id, process_definition_key, version, resource_name, created_at, updated_at
		FROM pesen_ke_warung_makan
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var model PesenKeWarungMakan
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.ProcessDefinitionKey,
		&model.Version,
		&model.ResourceName,
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

func (s *PesenKeWarungMakanStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM pesen_ke_warung_makan WHERE id = $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
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

func (s *PesenKeWarungMakanStore) Update(ctx context.Context, model *PesenKeWarungMakan) error {
	query := `
		UPDATE pesen_ke_warung_makan
		SET process_definition_key = $1, version = $2, resource_name = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id, process_definition_key, version, resource_name, created_at updated_at;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		model.ProcessDefinitionKey,
		model.Version,
		model.ResourceName,
		model.ID,
	).Scan(&model.ID, &model.ProcessDefinitionKey, &model.Version, &model.ResourceName, &model.CreatedAt, &model.UpdatedAt)
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

