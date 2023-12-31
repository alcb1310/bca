package tests

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

type DBMock struct{}

func (s *DBMock) LoadDummyData(companyId uuid.UUID) error {
	return nil
}

func (s *DBMock) Login(l *types.Login) (string, error) {
	return "", nil
}

func (s *DBMock) CreateCompany(company *types.CompanyCreate) error {
	return nil
}

func (s *DBMock) Health() map[string]string {
	health := make(map[string]string)
	health["message"] = "It's healthy"
	return health
}

func (s *DBMock) IsLoggedIn(token string, user uuid.UUID) bool {
	return true
}

func (s *DBMock) GetAllUsers(companyId uuid.UUID) ([]types.User, error) {
	users := []types.User{}

	users = append(users, types.User{
		Id:        uuid.New(),
		Email:     "test@test.com",
		CompanyId: companyId,
		Name:      "test",
		RoleId:    "a",
	})

	users = append(users, types.User{
		Id:        uuid.New(),
		Email:     "test2@test.com",
		CompanyId: companyId,
		Name:      "test2",
		RoleId:    "a",
	})

	return users, nil
}

func (s *DBMock) CreateUser(u *types.UserCreate) (types.User, error) {
	return types.User{}, nil
}

func (s *DBMock) GetUser(id, companyId uuid.UUID) (types.User, error) {
	return types.User{}, nil
}

func (s *DBMock) UpdateUser(u types.User, id, companyId uuid.UUID) (types.User, error) {
	return types.User{}, nil
}

func (s *DBMock) UpdatePassword(pass string, id, companyId uuid.UUID) (types.User, error) {
	return types.User{}, nil
}

func (s *DBMock) DeleteUser(id, companyId uuid.UUID) error {
	return nil
}

func (s *DBMock) GetAllProjects(companyId uuid.UUID) ([]types.Project, error) {
	return []types.Project{}, nil
}

func (s *DBMock) CreateProject(p types.Project) (types.Project, error) {
	return types.Project{}, nil
}

func (s *DBMock) GetProject(id, companyId uuid.UUID) (types.Project, error) {
	return types.Project{}, nil
}

func (s *DBMock) UpdateProject(p types.Project, id, companyId uuid.UUID) error {
	return nil
}

func (s *DBMock) GetAllSuppliers(companyId uuid.UUID) ([]types.Supplier, error) {
	return []types.Supplier{}, nil
}

func (s *DBMock) CreateSupplier(supplier *types.Supplier) error {
	supplier.ID = uuid.New()
	return nil
}

func (s *DBMock) GetOneSupplier(id, companyId uuid.UUID) (types.Supplier, error) {
	return types.Supplier{}, nil
}

func (s *DBMock) UpdateSupplier(supplier *types.Supplier) error {
	return nil
}

func (s *DBMock) GetBudgetItems(companyId uuid.UUID) ([]types.BudgetItem, error) {
	return []types.BudgetItem{}, nil
}

func (s *DBMock) CreateBudgetItem(bi *types.BudgetItem) error {
	bi.ID = uuid.New()
	return nil
}

func (s *DBMock) GetOneBudgetItem(id uuid.UUID, companyId uuid.UUID) (*types.BudgetItem, error) {
	return &types.BudgetItem{}, nil
}

func (s *DBMock) UpdateBudgetItem(bi *types.BudgetItem) error {
	return nil
}

func (s *DBMock) GetBudgets(companyId uuid.UUID) ([]types.GetBudget, error) {
	return []types.GetBudget{}, nil
}

func (s *DBMock) CreateBudget(b *types.CreateBudget) (types.Budget, error) {
	return types.Budget{}, nil
}

func (s *DBMock) GetBudgetsByProjectId(companyId, projectId uuid.UUID) ([]types.GetBudget, error) {
	return []types.GetBudget{}, nil
}

func (s *DBMock) GetOneBudget(companyId, projectId, budgetItemId uuid.UUID) (*types.GetBudget, error) {
	return &types.GetBudget{}, nil
}

func (s *DBMock) UpdateBudget(b *types.CreateBudget, budget *types.Budget) error {
	return nil
}
