package database

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"

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

func (s *service) GetUser(id, companyId uuid.UUID) (types.User, error) {
	u := types.User{}
	sql := "select id, name, email, company_id, role_id from \"user\" where id = $1 and company_id = $2"
	err := s.db.QueryRow(sql, id, companyId).Scan(&u.Id, &u.Name, &u.Email, &u.CompanyId, &u.RoleId)
	if err != nil {
		return types.User{}, err
	}

	return u, nil
}

func (s *service) UpdateUser(u types.User, id, companyId uuid.UUID) (types.User, error) {
	sql := "update \"user\" set name = $1, email = $2, company_id = $3, role_id = $4 where id = $5 and company_id = $6"
	_, err := s.db.Exec(sql, u.Name, u.Email, u.CompanyId, u.RoleId, id, companyId)
	if err != nil {
		return types.User{}, err
	}

	return types.User{
		Id:        id,
		Name:      u.Name,
		Email:     u.Email,
		CompanyId: companyId,
		RoleId:    u.RoleId,
	}, nil
}

func (s *service) UpdatePassword(pass string, id, companyId uuid.UUID) (types.User, error) {
	u := types.User{}
	hash, err := utils.EncryptPasssword(pass)
	if err != nil {
		return types.User{}, err
	}

	sql := "update \"user\" set password = $1 where id = $2 and company_id = $3 returning id, email, name, role_id"
	if err := s.db.QueryRow(sql, hash, id, companyId).Scan(&u.Id, &u.Email, &u.Name, &u.RoleId); err != nil {
		return types.User{}, err
	}
	u.CompanyId = companyId

	return u, nil
}

func (s *service) DeleteUser(id, companyId uuid.UUID) error {
	sql := "delete from \"user\" where id = $1 and company_id = $2"
	_, err := s.db.Exec(sql, id, companyId)
	return err
}
