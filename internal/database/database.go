package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"

	"bca-go-final/internal/types"
)

type Service interface {
	Health() map[string]string
	CreateCompany(company *types.CompanyCreate) error
	Login(l *types.Login) (string, *types.User, error)
	RegenerateToken(token string, user uuid.UUID) error
	IsLoggedIn(token string, user uuid.UUID) bool

	Levels(companyId uuid.UUID) []types.Select

	// database/dummy.go
	LoadDummyData(companyId uuid.UUID) error

	// database/users.go
	GetAllUsers(companyId uuid.UUID) ([]types.User, error)
	CreateUser(u *types.UserCreate) (types.User, error)
	GetUser(id, companyId uuid.UUID) (types.User, error)
	UpdateUser(u types.User, id, companyId uuid.UUID) (types.User, error)
	UpdatePassword(pass string, id, companyId uuid.UUID) (types.User, error)
	DeleteUser(id, companyId uuid.UUID) error

	// database/projects.go
	GetAllProjects(companyId uuid.UUID) ([]types.Project, error)
	CreateProject(p types.Project) (types.Project, error)
	GetProject(id, companyId uuid.UUID) (types.Project, error)
	UpdateProject(p types.Project, id, companyId uuid.UUID) error
	GetActiveProjects(companyId uuid.UUID, active bool) []types.Project

	// database/suppliers.go
	GetAllSuppliers(companyId uuid.UUID, search string) ([]types.Supplier, error)
	CreateSupplier(supplier *types.Supplier) error
	GetOneSupplier(id, companyId uuid.UUID) (types.Supplier, error)
	UpdateSupplier(supplier *types.Supplier) error

	// database/budget-items.go
	GetBudgetItems(companyId uuid.UUID, search string) ([]types.BudgetItemResponse, error)
	CreateBudgetItem(bi *types.BudgetItem) error
	GetOneBudgetItem(id uuid.UUID, companyId uuid.UUID) (*types.BudgetItem, error)
	UpdateBudgetItem(bi *types.BudgetItem) error
	GetBudgetItemsByAccumulate(companyId uuid.UUID, accum bool) []types.BudgetItem
	GetBudgetItemsByLevel(companyId uuid.UUID, level uint8) []types.BudgetItem
	GetNonAccumulateChildren(companyId, id *uuid.UUID, budgetItems []types.BudgetItem, results []uuid.UUID) []uuid.UUID

	// database/budget.go
	GetBudgets(companyId, project_id uuid.UUID, search string) ([]types.GetBudget, error)
	CreateBudget(b *types.CreateBudget) (types.Budget, error)
	GetBudgetsByProjectId(companyId, projectId uuid.UUID, level *uint8) ([]types.GetBudget, error)
	GetOneBudget(companyId, projectId, budgetItemId uuid.UUID) (*types.GetBudget, error)
	UpdateBudget(b types.CreateBudget, budget types.Budget) error

	// database/invoice.go
	GetInvoices(companyId uuid.UUID) ([]types.InvoiceResponse, error)
	CreateInvoice(invoice *types.InvoiceCreate) error
	GetOneInvoice(invoiceId, companyId uuid.UUID) (types.InvoiceResponse, error)
	UpdateInvoice(invoice types.InvoiceCreate) error
	DeleteInvoice(invoiceId, companyId uuid.UUID) error
	BalanceInvoice(invoice types.InvoiceResponse) error

	// database/invoice-details.GetOneBudget
	GetAllDetails(invoiceId, companyId uuid.UUID) ([]types.InvoiceDetailsResponse, error)
	AddDetail(detail types.InvoiceDetailCreate) error
	DeleteDetail(invoiceId, budgetItemId, companyId uuid.UUID) error

	// database/reports.go
	GetBalance(companyId, projectId uuid.UUID, date time.Time) types.BalanceResponse
	GetHistoricByProject(companyId, projectId uuid.UUID, date time.Time, level uint8) []types.GetBudget
	GetSpentByBudgetItem(companyId, projectId, budgetItemId uuid.UUID, date time.Time, ids []uuid.UUID) float64
	GetDetailsByBudgetItem(companyId, projectId, budgetItemId uuid.UUID, date time.Time, ids []uuid.UUID) []types.InvoiceDetails

	// database/closure.go
	CreateClosure(companyId, projectId uuid.UUID, date time.Time) error

	// database/categories.go
	GetAllCategories(companyId uuid.UUID) ([]types.Category, error)
	CreateCategory(category types.Category) error
	GetCategory(id, companyId uuid.UUID) (types.Category, error)
	UpdateCategory(category types.Category) error

	// database/materials.go
	GetAllMaterials(companyId uuid.UUID) []types.Material
	CreateMaterial(material types.Material) error
	GetMaterial(id, companyId uuid.UUID) (types.Material, error)
	UpdateMaterial(material types.Material) error

	// database/rubros.go
	GetAllRubros(companyId uuid.UUID) ([]types.Rubro, error)
	CreateRubro(rubro types.Rubro) (uuid.UUID, error)
	GetOneRubro(id, companyId uuid.UUID) (types.Rubro, error)
	UpdateRubro(rubro types.Rubro) error

	GetMaterialsByItem(id, companyId uuid.UUID) []types.ACU
	AddMaterialsByItem(itemId, materialId uuid.UUID, quantity float64, companyId uuid.UUID) error
	DeleteMaterialsByItem(itemId, materialId, companyId uuid.UUID) error
	UpdateMaterialByItem(itemId, materialId uuid.UUID, quantity float64, companyId uuid.UUID) error

	// database/acu.go
	CreateCantidades(projectId, rubroId uuid.UUID, quantity float64, companyId uuid.UUID) error
	DeleteCantidades(id, companyId uuid.UUID) error
	CantidadesTable(companyId uuid.UUID) []types.Quantity
	AnalysisReport(project_id, company_id uuid.UUID) map[string][]types.AnalysisReport
	GetQuantityByMaterialAndItem(itemId, materialId, companyId uuid.UUID) types.ItemMaterialType
	GetOneQuantityById(id, companyId uuid.UUID) types.Quantity
	UpdateQuantity(q types.Quantity, companyId uuid.UUID) error
}

type service struct {
	db *sql.DB
}

func New(connStr string) Service {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	s := &service{db: db}

	if err := createTables(db); err != nil {
		log.Fatalf(fmt.Sprintf("error creating tables. Err: %v", err))
	}

	if err := loadRoles(db); err != nil {
		log.Fatalf(fmt.Sprintf("error loading roles. Err: %v", err))
	}

	return s
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) Levels(companyId uuid.UUID) []types.Select {
	levels := []types.Select{}
	query := "select level from vw_levels where company_id = $1"
	rows, err := s.db.Query(query, companyId)
	if err != nil {
		log.Fatal(err)
		return levels
	}
	defer rows.Close()

	for rows.Next() {
		var level string
		if err := rows.Scan(&level); err != nil {
			log.Fatal(err)
			return levels
		}
		levels = append(levels, types.Select{Key: level, Value: level})
	}

	return levels
}
