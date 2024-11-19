package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

const (
	PembuatanArtikelID             = "creating_artikel"
	PembuatanArtikelName           = "Pembuatan Artikel"
	PembuatanArtikelFormID         = "creating_artikel_form"
	PembuatanArtikelAssignee       = ""
	PembuatanArtikelCandidateGroup = ""
	PembuatanArtikelCandidateUser  = ""
	PembuatanArtikelSchedule       = ``
)

type PembuatanArtikel struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	TaskId     string   `json:"task_id"`
	FormId     string   `json:"form_id"`
	Properties []string `json:"properties"`
	CreatedBy  int64    `json:"created_by"`
	UpdatedBy  *int64   `json:"updated_by"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	DeletedAt  *string  `json:"deleted_at"`
}

type PembuatanArtikelStore struct {
	db *sql.DB
}

func (s *PembuatanArtikelStore) Create(ctx context.Context, model *PembuatanArtikel) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.create(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}

func (s *PembuatanArtikelStore) Delete(ctx context.Context, id int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, id); err != nil {
			return err
		}
		return nil
	})
}

func (s *PembuatanArtikelStore) Update(ctx context.Context, model *PembuatanArtikel) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.update(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}

func (s *PembuatanArtikelStore) create(ctx context.Context, tx *sql.Tx, model *PembuatanArtikel) error {
	if model.Properties == nil {
		model.Properties = []string{}
	}

	propertiesJSON, errProperties := json.Marshal(model.Properties)
	if errProperties != nil {
		return errProperties
	}

	query := `
		INSERT INTO pembuatanartikel (name, form_id, properties, created_by, task_id)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		) RETURNING 
		 	id, name, form_id, task_id, properties, created_by, updated_by,
			created_at, updated_at
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var propertiesData []byte
	err := tx.QueryRowContext(
		ctx,
		query,
		model.Name,
		model.FormId,
		propertiesJSON,
		model.CreatedBy,
		model.UpdatedAt,
		model.TaskId,
	).Scan(
		&model.ID,
		&model.Name,
		&model.FormId,
		&model.TaskId,
		&propertiesData,
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

func (s *PembuatanArtikelStore) GetByID(ctx context.Context, id int64) (*PembuatanArtikel, error) {
	query := `
		SELECT id, name, form_id, task_id, properties, created_by,
		updated_by, created_at, updated_at
		FROM pembuatanartikel
		WHERE id = $1 AND deleted_at IS NULL
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var model PembuatanArtikel
	var propertiesData []byte
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.Name,
		&model.FormId,
		&model.TaskId,
		&propertiesData,
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

	if len(propertiesData) > 0 {
		if err := json.Unmarshal(propertiesData, &model.Properties); err != nil {
			return nil, err
		}
	} else {
		model.Properties = []string{}
	}

	return &model, nil
}

func (s *PembuatanArtikelStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `UPDATE pembuatanartikel SET deleted_at = NOW() WHERE id = $1;`

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

func (s *PembuatanArtikelStore) update(ctx context.Context, tx *sql.Tx, model *PembuatanArtikel) error {
	if model.Properties == nil {
		model.Properties = []string{}
	}

	propertiesJSON, errProperties := json.Marshal(model.Properties)
	if errProperties != nil {
		return errProperties
	}
	query := `
		UPDATE pembuatanartikel
		SET name = $1, form_id = $2, properties = $3, updated_by = $4, task_id = $5,  updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING id, name, form_id, task_id, properties, created_by, updated_by, created_at, updated_at;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var propertiesData []byte
	err := tx.QueryRowContext(
		ctx,
		query,
		model.Name,
		model.FormId,
		propertiesJSON,
		model.UpdatedBy,
		model.TaskId,
		model.ID,
	).Scan(&model.ID,
		&model.Name,
		&model.FormId,
		&model.TaskId,
		propertiesData,
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
