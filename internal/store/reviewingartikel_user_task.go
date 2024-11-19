package store
import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

const (
	ReviewingArtikelID = "reviewing_artikel"
	ReviewingArtikelName = "Reviewing Artikel"
	ReviewingArtikelFormID = "reviewing_artikel_form"
	ReviewingArtikelAssignee = ""
	ReviewingArtikelCandidateGroup = ""
	ReviewingArtikelCandidateUser = ""
	ReviewingArtikelSchedule = ``

)

type ReviewingArtikel struct {
    ID         int64    `json:"id"`
	Name       string   `json:"name"`
	FormId     string   `json:"form_id"`
	Properties []string `json:"properties"`
	CreatedBy  int64    `json:"created_by"`
	UpdatedBy  *int64   `json:"updated_by"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	DeletedAt  *string  `json:"deleted_at"`
}

type ReviewingArtikelStore struct {
	db *sql.DB
}

func (s *ReviewingArtikelStore) Create(ctx context.Context, model *ReviewingArtikel) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.create(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}

func (s *ReviewingArtikelStore) Delete(ctx context.Context, id int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, id); err != nil {
			return err
		}
		return nil
	})	
}

func (s *ReviewingArtikelStore) Update(ctx context.Context, model *ReviewingArtikel) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.update(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}
	
func (s *ReviewingArtikelStore) create(ctx context.Context, tx *sql.Tx, model *ReviewingArtikel) error {
	if model.Properties == nil {
		model.Properties = []string{}
	}

	propertiesJSON, errProperties := json.Marshal(model.Properties)
	if errProperties != nil {
		return errProperties
	}

	query := `
		INSERT INTO reviewingartikel (name, form_id, properties, created_by)
		VALUES (
			$1,
			$2,
			$3,
			$4
		) RETURNING 
		 	id, name, form_id, properties, created_by, updated_by,
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
	).Scan(
		&model.ID,
		&model.Name,
		&model.FormId,
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

func (s *ReviewingArtikelStore) GetByID(ctx context.Context, id int64) (*ReviewingArtikel, error) {
	query := `
		SELECT id, name, form_id, properties, created_by, 
		updated_by, created_at, updated_at
		FROM reviewingartikel
		WHERE id = $1 AND deleted_at IS NULL
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var model ReviewingArtikel
	var propertiesData []byte
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.Name,
		&model.FormId,
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

func (s *ReviewingArtikelStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `UPDATE reviewingartikel SET deleted_at = NOW() WHERE id = $1;`

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

func (s *ReviewingArtikelStore) update(ctx context.Context, tx *sql.Tx, model *ReviewingArtikel) error {

	if model.Properties == nil {
		model.Properties = []string{}
	}

	propertiesJSON, errProperties := json.Marshal(model.Properties)
	if errProperties != nil {
		return errProperties
	}
	query := `
		UPDATE reviewingartikel
		SET name = $1, form_id = $2, properties = $3, updated_by = $4  updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING id, name, form_id, properties, created_by, updated_by, created_at, updated_at;
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
		model.ID,
	).Scan(&model.ID, 
		&model.Name, 
		&model.FormId, 
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

