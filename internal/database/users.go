package database

import (
	"bca-go-final/internal/types"
	"log"

	"github.com/google/uuid"
)

func (s *service) GetAllUsers(companyId uuid.UUID) ([]types.User, error) {
	users := []types.User{}

	sql := "select id, name, email, company_id, role_id from \"user\" where company_id = $1"
	rows, err := s.db.Query(sql, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := types.User{}
		if err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.CompanyId, &u.RoleId); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	log.Println("GetAllUsers: ", users)

	return users, nil
}

func (s *service) CreateUser(u *types.UserCreate) (types.User, error) {
	sql := "insert into \"user\" (name, email, password, company_id, role_id) values ($1, $2, $3, $4, $5) returning id"
	err := s.db.QueryRow(sql, u.Name, u.Email, u.Password, u.CompanyId, u.RoleId).Scan(&u.Id)
	if err != nil {
		return types.User{}, err
	}

	us := types.User{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		CompanyId: u.CompanyId,
		RoleId:    u.RoleId,
	}

	return us, nil
}
