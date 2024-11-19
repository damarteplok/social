package store
import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

const (
	ApprovingArtikelID = "approving_artikel"
	ApprovingArtikelName = "Approving Artikel"
	ApprovingArtikelFormID = "approving_artikel_form"
	ApprovingArtikelAssignee = ""
	ApprovingArtikelCandidateGroup = ""
	ApprovingArtikelCandidateUser = ""
	ApprovingArtikelSchedule = ``

)

type ApprovingArtikel struct {
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

type ApprovingArtikelStore struct {
	db *sql.DB
}

func (s *ApprovingArtikelStore) Create(ctx context.Context, model *ApprovingArtikel) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.create(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}

func (s *ApprovingArtikelStore) Delete(ctx context.Context, id int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, id); err != nil {
			return err
		}
		return nil
	})	
}

func (s *ApprovingArtikelStore) Update(ctx context.Context, model *ApprovingArtikel) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.update(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}
	
func (s *ApprovingArtikelStore) create(ctx context.Context, tx *sql.Tx, model *ApprovingArtikel) error {
	if model.Properties == nil {
		model.Properties = []string{}
	}

	propertiesJSON, errProperties := json.Marshal(model.Properties)
	if errProperties != nil {
		return errProperties
	}

	query := `
		INSERT INTO approvingartikel (name, form_id, properties, created_by)
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

func (s *ApprovingArtikelStore) GetByID(ctx context.Context, id int64) (*ApprovingArtikel, error) {
	query := `
		SELECT id, name, form_id, properties, created_by, 
		updated_by, created_at, updated_at
		FROM approvingartikel
		WHERE id = $1 AND deleted_at IS NULL
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var model ApprovingArtikel
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

func (s *ApprovingArtikelStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `UPDATE approvingartikel SET deleted_at = NOW() WHERE id = $1;`

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

func (s *ApprovingArtikelStore) update(ctx context.Context, tx *sql.Tx, model *ApprovingArtikel) error {

	if model.Properties == nil {
		model.Properties = []string{}
	}

	propertiesJSON, errProperties := json.Marshal(model.Properties)
	if errProperties != nil {
		return errProperties
	}
	query := `
		UPDATE approvingartikel
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

