package store

import (
	"context"
	"database/sql"
	"errors"
)

const (
	PembuatanMediaBeritaTechnologyVersion              = 1
	PembuatanMediaBeritaTechnologyProcessDefinitionKey = 2251799814076217
	PembuatanMediaBeritaTechnologyResourceName         = "pembuatan_media_berita_technology.bpmn"
)

// TODO: UPDATE THIS STRUCT AND CODE BELOW
type PembuatanMediaBeritaTechnology struct {
	ID                   int64   `json:"id"`
	ProcessDefinitionKey int64   `json:"process_definition_key"`
	Version              int32   `json:"version"`
	ResourceName         string  `json:"resource_name"`
	ProcessInstanceKey   int64   `json:"process_instance_key"`
	TaskDefinitionId     *string `json:"task_definition_id"`
	TaskState            *string `json:"task_state"`
	CreatedBy            int64   `json:"created_by"`
	UpdatedBy            *int64  `json:"updated_by"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
	DeletedAt            *string `json:"deleted_at"`
}

type PembuatanMediaBeritaTechnologyStore struct {
	db *sql.DB
}

func (s *PembuatanMediaBeritaTechnologyStore) Create(ctx context.Context, model *PembuatanMediaBeritaTechnology) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.create(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}

func (s *PembuatanMediaBeritaTechnologyStore) Delete(ctx context.Context, id int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, id); err != nil {
			return err
		}
		return nil
	})
}

func (s *PembuatanMediaBeritaTechnologyStore) Update(ctx context.Context, model *PembuatanMediaBeritaTechnology) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.update(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
}

func (s *PembuatanMediaBeritaTechnologyStore) create(ctx context.Context, tx *sql.Tx, model *PembuatanMediaBeritaTechnology) error {
	// model.Version = 1
	// model.ProcessDefinitionKey = 2251799814076217
	model.ResourceName = "pembuatan_media_berita_technology.bpmn"

	query := `
		INSERT INTO pembuatan_media_berita_technology (
			process_definition_key, version, 
			resource_name, process_instance_key, created_by
		) VALUES (
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

func (s *PembuatanMediaBeritaTechnologyStore) GetByID(ctx context.Context, id int64) (*PembuatanMediaBeritaTechnology, error) {
	query := `
		SELECT id, process_definition_key, version, 
			resource_name, process_instance_key, 
			task_definition_id, task_state,
			created_by, updated_by, created_at, updated_at
		FROM pembuatan_media_berita_technology
		WHERE id = $1 AND deleted_at IS NULL
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var model PembuatanMediaBeritaTechnology
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.ProcessDefinitionKey,
		&model.Version,
		&model.ResourceName,
		&model.ProcessInstanceKey,
		&model.TaskDefinitionId,
		&model.TaskState,
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

func (s *PembuatanMediaBeritaTechnologyStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `UPDATE pembuatan_media_berita_technology SET deleted_at = NOW() WHERE id = $1;`

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

func (s *PembuatanMediaBeritaTechnologyStore) update(ctx context.Context, tx *sql.Tx, model *PembuatanMediaBeritaTechnology) error {
	query := `
		UPDATE pembuatan_media_berita_technology
		SET process_definition_key = $1, 
			version = $2, 
			resource_name = $3, 
			process_instance_key = $4, 
			updated_by = $5, 
			updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING id, process_definition_key, 
			version, 
			resource_name, 
			process_instance_key, 
			created_by, 
			updated_by, 
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

func (s *PembuatanMediaBeritaTechnologyStore) Search(ctx context.Context, pq PaginatedQuery) (map[string]interface{}, error) {
	sortOrder := "DESC"
	if pq.Sort == "desc" || pq.Sort == "DESC" {
		sortOrder = "DESC"
	}
	if pq.Sort == "asc" || pq.Sort == "ASC" {
		sortOrder = "ASC"
	}

	// Base Query
	query := `
        SELECT p.id, p.process_definition_key, p.version,
            p.resource_name, p.process_instance_key,
            p.task_definition_id, p.task_state,
            p.created_by, p.updated_by, p.created_at, p.updated_at
        FROM pembuatan_media_berita_technology p
        LEFT JOIN users u ON p.created_by = u.id
        WHERE p.deleted_at IS NULL
    `

	var params []interface{}
	params = append(params, pq.Limit, pq.Offset)

	if pq.Search != "" {
		query += `
			AND (
				p.process_definition_key::text ILIKE '%' || $3 || '%' OR 
				p.resource_name ILIKE '%' || $3 || '%' OR
				p.process_instance_key::text ILIKE '%' || $3 || '%' OR
				p.task_definition_id::text ILIKE '%' || $3 || '%' OR
				p.task_state ILIKE '%' || $3 || '%' OR
				u.email ILIKE '%' || $3 || '%' OR
				u.username ILIKE '%' || $3 || '%'
			)
		`
		params = append(params, pq.Search)
	}

	if pq.Since != "" && pq.Until != "" {
		query += `
            AND p.created_at BETWEEN $4 AND $5
        `
		params = append(params, pq.Since, pq.Until)
	}

	query += `
        ORDER BY p.created_at ` + sortOrder + `
        LIMIT $1 OFFSET $2
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var results []*PembuatanMediaBeritaTechnology
	rows, err := s.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item PembuatanMediaBeritaTechnology
		if err := rows.Scan(
			&item.ID,
			&item.ProcessDefinitionKey,
			&item.Version,
			&item.ResourceName,
			&item.ProcessInstanceKey,
			&item.TaskDefinitionId,
			&item.TaskState,
			&item.CreatedBy,
			&item.UpdatedBy,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	countQuery := `
        SELECT COUNT(*)
        FROM pembuatan_media_berita_technology p
        LEFT JOIN users u ON p.created_by = u.id
        WHERE p.deleted_at IS NULL
    `
	var args []interface{}
	if pq.Search != "" {
		countQuery += `
            AND (
                p.process_definition_key::text ILIKE '%' || $1 || '%' OR 
                p.resource_name ILIKE '%' || $1 || '%' OR
                p.process_instance_key::text ILIKE '%' || $1 || '%' OR
                p.task_definition_id::text ILIKE '%' || $1 || '%' OR
                p.task_state ILIKE '%' || $1 || '%' OR
                u.email ILIKE '%' || $1 || '%' OR
                u.username ILIKE '%' || $1 || '%'
            )
        `
		args = append(args, pq.Search)
	}

	var totalCount int
	err = s.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	totalPages := (totalCount + pq.Limit - 1) / pq.Limit

	response := map[string]interface{}{
		"content":      results,
		"totalElement": totalCount,
		"totalPages":   totalPages,
		"limit":        pq.Limit,
		"offset":       pq.Offset,
		"sort":         pq.Sort,
		"search":       pq.Search,
		"since":        pq.Since,
		"until":        pq.Until,
	}

	return response, nil
}
