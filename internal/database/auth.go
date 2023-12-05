package database

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
)

func (s *service) CreateCompany(company *types.CompanyCreate) error {
	type favContextKey string

	var id uuid.UUID
	var role string
	k := favContextKey("company")
	ctx := context.WithValue(context.Background(), k, &company)
	tx, _ := s.db.Begin()
	defer tx.Rollback()

	sql := "insert into company (ruc, name, employees) values ($1, $2, $3) returning id"
	if err := tx.QueryRowContext(ctx, sql, company.Ruc, company.Name, company.Employees).Scan(&id); err != nil {
		return err
	}

	sql = "select id from role where name = 'admin'"
	if err := tx.QueryRowContext(ctx, sql).Scan(&role); err != nil {
		return err
	}
	log.Println("Role id: ", role)

	sql = "insert into \"user\" (name, email, password, company_id, role_id) values ($1, $2, $3, $4, $5)"
	pass, err := utils.EncryptPasssword(company.Password)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, sql, company.Name, company.Email, pass, id, role); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (s *service) Login(l *types.Login) (string, error) {
	sql := "select password, id, name, company_id, role_id from \"user\" where email = $1"
	u := &types.User{}
	var password string
	if err := s.db.QueryRowContext(context.Background(), sql, l.Email).Scan(&password, &u.Id, &u.Name, &u.CompanyId, &u.RoleId); err != nil {
		return "", errors.New("invalid credentials")
	}

	if _, err := utils.ComparePassword(password, l.Password); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(*u)
	if err != nil {
		return "", errors.New("server error")
	}

	sql = "insert into logged_in (user_id, token) values ($1, $2) on conflict (user_id) do update set token = $2"
	if _, err := s.db.ExecContext(context.Background(), sql, u.Id, token); err != nil {
		return "", errors.New("server error")
	}

	return token, nil
}
