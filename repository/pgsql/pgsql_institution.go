package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"prakarsa-app/domain"
	"prakarsa-app/transport/request"
	"prakarsa-app/transport/response"
	"prakarsa-app/utils"
	"strings"
)

type pgsqlInstitutionRepository struct {
	db *sql.DB
}

// NewPgsqlInstitutionRepository will create new an todoRepository object representation of InstitutionRepository interface
func NewPgsqlInstitutionRepository(db *sql.DB) *pgsqlInstitutionRepository {
	return &pgsqlInstitutionRepository{
		db: db,
	}
}

func (r *pgsqlInstitutionRepository) Create(ctx context.Context, institution *domain.Institution) (err error) {
	query := `INSERT INTO institutions (id, name, type, alias, is_active, created_by, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err = r.db.ExecContext(ctx, query, institution.ID, institution.Name, institution.Type, institution.Alias, institution.IsActive,
		institution.CreatedBy, institution.CreatedAt); err != nil {
		return err
	}

	return
}

func (r *pgsqlInstitutionRepository) Update(ctx context.Context, institution *domain.Institution) (err error) {
	// Build dynamic SET clauses from Institution struct
	sets := []string{}
	args := []interface{}{}
	idx := 1

	if institution.Name != "" {
		sets = append(sets, fmt.Sprintf("name = $%d", idx))
		args = append(args, institution.Name)
		idx++
	}

	if institution.Alias != "" {
		sets = append(sets, fmt.Sprintf("alias = $%d", idx))
		args = append(args, institution.Alias)
		idx++
	}

	if institution.Type != "" {
		sets = append(sets, fmt.Sprintf("type = $%d", idx))
		args = append(args, institution.Type)
		idx++
	}

	if institution.IsActive != nil {
		sets = append(sets, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, institution.IsActive)
		idx++
	}

	// kalau ada sesuatu untuk di‐update, commit ke SQL
	if len(sets) > 0 {
		// Update stamp
		sets = append(sets, fmt.Sprintf("updated_at = $%d", idx))
		args = append(args, institution.UpdatedAt)
		idx++

		sets = append(sets, fmt.Sprintf("updated_by = $%d", idx))
		args = append(args, institution.UpdatedBy)
		idx++

		// tambahkan WHERE id = $idx
		args = append(args, institution.ID)
		query := fmt.Sprintf(
			"UPDATE institutions SET %s WHERE id = $%d",
			strings.Join(sets, ", "),
			idx,
		)

		if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
			return
		}
	}

	return
}

func (r *pgsqlInstitutionRepository) Delete(ctx context.Context, institution *domain.Institution) (rowsAffected int64, err error) {
	query := "DELETE FROM institutions WHERE id = $1"
	res, err := r.db.ExecContext(ctx, query, institution.ID)
	if err != nil {
		return
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return
	}

	return
}

func (r *pgsqlInstitutionRepository) GetList(ctx context.Context, request *request.GetListInstitutionReq) (res []response.GetListInstitutionRes, meta response.MetaRes, err error) {
	// 1. Build WHERE clauses
	wheres := []string{}
	args := []interface{}{}
	idx := 1

	if request.Name != "" {
		wheres = append(wheres, fmt.Sprintf("name ILIKE $%d", idx))
		args = append(args, "%"+request.Name+"%")
		idx++
	}

	if request.Type != "" {
		wheres = append(wheres, fmt.Sprintf("type ILIKE $%d", idx))
		args = append(args, "%"+request.Type+"%")
		idx++
	}

	if request.Alias != "" {
		wheres = append(wheres, fmt.Sprintf("alias ILIKE $%d", idx))
		args = append(args, "%"+request.Alias+"%")
		idx++
	}

	if request.IsActive != nil {
		wheres = append(wheres, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, request.IsActive)
		idx++
	}

	whereSQL := ""
	if len(wheres) > 0 {
		whereSQL = "WHERE " + strings.Join(wheres, " AND ")
	}

	// --- 2. Hitung totalCount dulu (tanpa LIMIT/OFFSET) ---
	countQuery := fmt.Sprintf(
		"SELECT COUNT(*) FROM institutions %s",
		whereSQL,
	)
	if err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&meta.TotalData); err != nil {
		return nil, meta, err
	}

	// 2. Calculate LIMIT & OFFSET
	perPage := request.PerPage
	if perPage <= 0 {
		perPage = 10
	}
	page := request.Page
	if page <= 0 {
		page = 1
	}

	// total pages = ceil(total / perPage)
	meta.Page = page
	meta.PerPage = perPage
	meta.TotalPages = (meta.TotalData + perPage - 1) / perPage

	offset := (page - 1) * perPage

	// add LIMIT & OFFSET to args
	args = append(args, perPage, offset)
	limitPos, offsetPos := idx, idx+1

	// 3. Final query
	query := fmt.Sprintf(`
        SELECT
            id, name, type, alias, is_active
        FROM institutions
        %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereSQL, limitPos, offsetPos)

	// 4. Execute
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, meta, err
	}
	defer rows.Close()

	// 5. Scan results
	for rows.Next() {
		var item response.GetListInstitutionRes

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Type,
			&item.Alias,
			&item.IsActive,
		); err != nil {
			return nil, meta, err
		}

		res = append(res, item)
	}
	if errRow := rows.Err(); errRow != nil {
		return nil, meta, errRow
	}

	return
}

func (r *pgsqlInstitutionRepository) GetDetail(ctx context.Context, request *request.GetDetailInstitutionReq) (res domain.Institution, err error) {

	const query = `
					SELECT
					  id,
					  name,
					  alias,
					  type,
					  is_active,
					  created_at,
					  created_by,
					  updated_at,
					  updated_by,
					  deleted_at
					FROM institutions
					WHERE id = $1
					LIMIT 1
					`

	// 1. QueryRowContext untuk ambil satu baris
	row := r.db.QueryRowContext(ctx, query, request.ID)

	// 2. Scan kolom ke field di domain.Partner
	// since created_at is NOT NULL int8:
	var createdAt int64
	// updated_at/deleted_at can be NULL, so use NullInt64:
	var updatedAt, deletedAt sql.NullInt64
	var updatedBy sql.NullString

	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.Alias,
		&res.Type,
		&res.IsActive,
		&createdAt,
		&res.CreatedBy,
		&updatedAt,
		&updatedBy,
		&deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, utils.NewNotFoundError("institution not found")
		}
		return res, err
	}

	// assign into your domain fields
	res.CreatedAt = createdAt
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.Int64
	}
	if deletedAt.Valid {
		res.DeletedAt = deletedAt.Int64
	}
	if updatedBy.Valid {
		res.UpdatedBy = updatedBy.String
	}

	return
}
