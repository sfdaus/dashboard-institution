package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"prakarsa-app/domain"
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
