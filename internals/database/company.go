package database

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/alcb1310/bca/internals/types"
)

func (s *service) CreateCompany(c *types.Company, u *types.CreateUser) error {
	tx, _ := s.DB.Begin()
	defer tx.Rollback()

	query := "insert into company (ruc, name, employees, is_active) values($1, $2, $3, $4) returning id"
	if err := tx.QueryRow(query, c.Ruc, c.Name, c.Employees, c.IsActive).Scan(&c.ID); err != nil {
		if e, ok := err.(*pgconn.PgError); ok {
			slog.Error("Error saving the company", "code", e.Code, "message", e.Message)
			return e
		}

		slog.Error("Unknown error, saving the company", "message", err.Error())
		return err
	}

	u.CompanyID = c.ID

	query = "insert into \"user\" (email, name, password, company_id) values($1, $2, $3, $4) returning id"
	if err := tx.QueryRow(query, u.Email, u.Name, u.Password, u.CompanyID).Scan(&u.ID); err != nil {
		if e, ok := err.(*pgconn.PgError); ok {
			slog.Error("Error saving the user", "code", e.Code, "message", e.Message)
			return e
		}

		slog.Error("Unknown error saving the user", "message", err.Error())
		return err
	}

	tx.Commit()
	return nil
}
