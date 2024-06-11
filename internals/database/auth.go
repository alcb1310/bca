package database

import (
	"errors"

	"github.com/alcb1310/bca/internals/types"
	"github.com/alcb1310/bca/internals/utils"
)

func (s *service) Login(email, password string) (types.User, error) {
	var user = types.User{
		Email: email,
	}
	var pass string

	query := "select id, password, name, company_id from \"user\" where email = $1"
	if err := s.DB.QueryRow(query, email).Scan(&user.ID, &pass, &user.Name, &user.CompanyID); err != nil {
		return user, err
	}

	if utils.IsValidPassword(pass, password) {
		return user, errors.New("Invalid credentials")
	}

	return user, nil
}
