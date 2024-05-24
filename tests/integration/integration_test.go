package integration_test

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"testing"
// 	"time"
//
// 	"github.com/jackc/pgx/v5/pgconn"
// 	_ "github.com/onsi/ginkgo/v2"
// 	_ "github.com/onsi/gomega"
// 	"github.com/testcontainers/testcontainers-go"
// 	"github.com/testcontainers/testcontainers-go/modules/postgres"
// 	"github.com/testcontainers/testcontainers-go/wait"
//
// 	databasem "bca-go-final/internal/database"
// 	"bca-go-final/internal/types"
// 	"bca-go-final/internal/utils"
// )
//
// var (
// 	dbName     = "bca-test"
// 	dbUser     = "bca-test"
// 	dbPassword = "bca-test"
// 	secret     = "super-secret-key"
// )
//
// func TestClient(t *testing.T) {
// 	RegisterFailHandler(Fail)
// 	t.Setenv("SECRET", secret)
// 	RunSpecs(t, "Integration Test Suite")
// }
//
// var _ = Describe("Integration Test", Ordered, func() {
// 	var ctx context.Context
// 	var container *postgres.PostgresContainer
// 	var connString string
// 	loggedIn := false
// 	var db databasem.Service
// 	trueVal := true
// 	var payload utils.Payload
// 	var token string
//
// 	BeforeAll(func() {
// 		ctx = context.Background()
//
// 		postgresContainer, err := postgres.RunContainer(ctx,
// 			testcontainers.WithImage("docker.io/postgres:16-alpine"),
// 			postgres.WithDatabase(dbName),
// 			postgres.WithUsername(dbUser),
// 			postgres.WithPassword(dbPassword),
// 			testcontainers.WithWaitStrategy(
// 				wait.ForLog("database system is ready to accept connections").
// 					WithOccurrence(2).
// 					WithStartupTimeout(5*time.Second),
// 			),
// 		)
// 		Expect(err).NotTo(HaveOccurred())
//
// 		connString, err = postgresContainer.ConnectionString(ctx)
// 		Expect(err).NotTo(HaveOccurred())
//
// 		db = databasem.New(fmt.Sprintf("%ssslmode=disable", connString), "../../internal/database/tables.sql")
//
// 		h := db.Health()
// 		Expect(h["message"]).To(Equal("It's healthy"))
//
// 		container = postgresContainer
// 	})
//
// 	AfterAll(func() {
// 		err := container.Terminate(ctx)
// 		Expect(err).NotTo(HaveOccurred())
// 	})
//
// 	BeforeEach(func() {
// 		if loggedIn {
// 			maker, err := utils.NewJWTMaker(secret)
// 			Expect(err).NotTo(HaveOccurred())
// 			tokenData, err := maker.VerifyToken(token)
// 			payload = *tokenData
// 		}
// 	})
//
// 	When("User", func() {
// 		When("Register a new company", func() {
// 			It("should create a new company", func() {
// 				err := db.CreateCompany(&types.CompanyCreate{
// 					Ruc:       "1234567890",
// 					Name:      "Test Company",
// 					Email:     "test@test.com",
// 					Password:  "test",
// 					Employees: 100,
// 					User:      "Test",
// 				})
//
// 				Expect(err).NotTo(HaveOccurred())
// 			})
//
// 			It("should fail if company already exists", func() {
// 				err := db.CreateCompany(&types.CompanyCreate{
// 					Ruc:       "1234567890",
// 					Name:      "Test Company",
// 					Email:     "test@test.com",
// 					Password:  "test",
// 					Employees: 100,
// 					User:      "Test",
// 				})
//
// 				Expect(err).To(HaveOccurred())
// 			})
//
// 			It("should fail if user's email already exists", func() {
// 				err := db.CreateCompany(&types.CompanyCreate{
// 					Ruc:       "0123456789",
// 					Name:      "Test Company Again",
// 					Email:     "test@test.com",
// 					Password:  "test",
// 					Employees: 100,
// 					User:      "Test",
// 				})
//
// 				Expect(err).To(HaveOccurred())
// 			})
// 		})
//
// 		When("Login to the app", func() {
// 			It("should fail if invalid credentials", func() {
// 				token, err := db.Login(&types.Login{
// 					Email:    "invalid@whatever.com",
// 					Password: "adfadf",
// 				})
//
// 				Expect(err).NotTo(BeNil())
// 				Expect(token).To(BeEmpty())
// 			})
//
// 			It("should fail if invalid password", func() {
// 				token, err := db.Login(&types.Login{
// 					Email:    "test@test.com",
// 					Password: "invalid",
// 				})
//
// 				Expect(err).NotTo(BeNil())
// 				Expect(token).To(BeEmpty())
// 			})
//
// 			It("should login successfully", func() {
// 				var err error
// 				token, err = db.Login(&types.Login{
// 					Email:    "test@test.com",
// 					Password: "test",
// 				})
//
// 				Expect(err).To(BeNil())
// 				Expect(token).NotTo(BeEmpty())
// 				ctx = context.WithValue(ctx, "token", token)
// 				loggedIn = true
// 			})
// 		})
// 	})
//
// 	When("Settings", func() {
// 		When("Should create a project", func() {
// 			It("should create a project", func() {
// 				project, err := db.CreateProject(types.Project{
// 					Name:      "Test Project",
// 					IsActive:  &trueVal,
// 					CompanyId: payload.CompanyId,
// 					GrossArea: 100,
// 					NetArea:   100,
// 				})
//
// 				Expect(err).NotTo(HaveOccurred())
// 				Expect(project.ID).NotTo(BeEmpty())
// 			})
//
// 			It("should fail if project already exists", func() {
// 				project, err := db.CreateProject(types.Project{
// 					Name:      "Test Project",
// 					IsActive:  &trueVal,
// 					CompanyId: payload.CompanyId,
// 					GrossArea: 100,
// 					NetArea:   100,
// 				})
//
// 				Expect(err).To(HaveOccurred())
// 				Expect(project.ID.String()).To(Equal("00000000-0000-0000-0000-000000000000"))
// 			})
//
// 			It("should create a second project", func() {
// 				project, err := db.CreateProject(types.Project{
// 					Name:      "Test Project 1",
// 					IsActive:  &trueVal,
// 					CompanyId: payload.CompanyId,
// 					GrossArea: 100,
// 					NetArea:   100,
// 				})
//
// 				Expect(err).NotTo(HaveOccurred())
// 				Expect(project.ID).NotTo(BeEmpty())
// 			})
//
// 			It("should get all projects", func() {
// 				projects, err := db.GetAllProjects(payload.CompanyId)
//
// 				Expect(err).NotTo(HaveOccurred())
// 				Expect(len(projects)).To(Equal(2))
//
// 				for _, p := range projects {
// 					log.Printf("%+v", p)
// 					Expect(p.CompanyId).To(Equal(payload.CompanyId))
// 				}
// 			})
// 		})
//
// 		When("Creating a budget item", func() {
// 			It("should create a parent budget item", func() {
// 				err := db.CreateBudgetItem(&types.BudgetItem{
// 					CompanyId:  payload.CompanyId,
// 					Code:       "500",
// 					Name:       "Costo Directo",
// 					ParentId:   nil,
// 					Accumulate: sql.NullBool{Valid: true, Bool: true},
// 				})
//
// 				Expect(err).NotTo(HaveOccurred())
// 			})
//
// 			It("should fail if budget item already exists", func() {
// 				err := db.CreateBudgetItem(&types.BudgetItem{
// 					CompanyId:  payload.CompanyId,
// 					Code:       "500",
// 					Name:       "Costo Directo",
// 					ParentId:   nil,
// 					Accumulate: sql.NullBool{Valid: true, Bool: true},
// 				})
//
// 				Expect(err).To(HaveOccurred())
// 			})
//
// 			It("should create budget item schema", func() {
// 				budgetItems, err := db.GetBudgetItems(payload.CompanyId, "Costo Directo")
// 				Expect(err).NotTo(HaveOccurred())
// 				Expect(len(budgetItems)).To(Equal(1))
//
// 				budgetItem := budgetItems[0]
//
// 				err = db.CreateBudgetItem(&types.BudgetItem{
// 					CompanyId:  payload.CompanyId,
// 					Code:       "500.2",
// 					Name:       "Obra Gruesa",
// 					ParentId:   &budgetItem.ID,
// 					Accumulate: sql.NullBool{Valid: true, Bool: true},
// 				})
//
// 				Expect(err).NotTo(HaveOccurred())
//
// 				err = db.CreateBudgetItem(&types.BudgetItem{
// 					CompanyId:  payload.CompanyId,
// 					Code:       "500.3",
// 					Name:       "Terminaciones",
// 					ParentId:   &budgetItem.ID,
// 					Accumulate: sql.NullBool{Valid: true, Bool: true},
// 				})
//
// 				Expect(err).NotTo(HaveOccurred())
//
// 				budgetItems, err = db.GetBudgetItems(payload.CompanyId, "Obra Gruesa")
// 				Expect(err).NotTo(HaveOccurred())
//
// 				err = db.CreateBudgetItem(&types.BudgetItem{
// 					CompanyId:  payload.CompanyId,
// 					Code:       "500.2.1",
// 					Name:       "Hierro y alambres",
// 					ParentId:   &budgetItems[0].ID,
// 					Accumulate: sql.NullBool{Valid: true, Bool: true},
// 				})
// 				Expect(err).NotTo(HaveOccurred())
//
// 				err = db.CreateBudgetItem(&types.BudgetItem{
// 					CompanyId:  payload.CompanyId,
// 					Code:       "500.2.2",
// 					Name:       "Cemento y hormigones",
// 					ParentId:   &budgetItems[0].ID,
// 					Accumulate: sql.NullBool{Valid: true, Bool: true},
// 				})
// 				Expect(err).NotTo(HaveOccurred())
//
// 				budgetItems, err = db.GetBudgetItems(payload.CompanyId, "Hierro y alambres")
// 				Expect(err).NotTo(HaveOccurred())
//
// 				err = db.CreateBudgetItem(&types.BudgetItem{
// 					CompanyId:  payload.CompanyId,
// 					Code:       "500.2.1.1",
// 					Name:       "Varilla 8 mm",
// 					ParentId:   &budgetItems[0].ID,
// 					Accumulate: sql.NullBool{Valid: true, Bool: true},
// 				})
// 				Expect(err).NotTo(HaveOccurred())
//
// 				err = db.CreateBudgetItem(&types.BudgetItem{
// 					CompanyId:  payload.CompanyId,
// 					Code:       "500.2.1.2",
// 					Name:       "Varilla 10 mm",
// 					ParentId:   &budgetItems[0].ID,
// 					Accumulate: sql.NullBool{Valid: true, Bool: true},
// 				})
// 				Expect(err).NotTo(HaveOccurred())
//
// 				err = db.CreateBudgetItem(&types.BudgetItem{
// 					CompanyId:  payload.CompanyId,
// 					Code:       "500.2.1.3",
// 					Name:       "Varilla 12 mm",
// 					ParentId:   &budgetItems[0].ID,
// 					Accumulate: sql.NullBool{Valid: true, Bool: true},
// 				})
// 				Expect(err).NotTo(HaveOccurred())
//
// 				budgetItems, err = db.GetBudgetItems(payload.CompanyId, "Cemento y hormigones")
// 				Expect(err).NotTo(HaveOccurred())
//
// 				err = db.CreateBudgetItem(&types.BudgetItem{
// 					CompanyId:  payload.CompanyId,
// 					Code:       "500.2.2.1",
// 					Name:       "Cemento",
// 					ParentId:   &budgetItems[0].ID,
// 					Accumulate: sql.NullBool{Valid: true, Bool: true},
// 				})
// 				Expect(err).NotTo(HaveOccurred())
//
// 				err = db.CreateBudgetItem(&types.BudgetItem{
// 					CompanyId:  payload.CompanyId,
// 					Code:       "500.2.2.2",
// 					Name:       "Hormigon 210kg/cm2",
// 					ParentId:   &budgetItems[0].ID,
// 					Accumulate: sql.NullBool{Valid: true, Bool: true},
// 				})
// 				Expect(err).NotTo(HaveOccurred())
// 			})
// 		})
//
// 		When("Suppliers", func() {
// 			It("should create a supplier", func() {
// 				err := db.CreateSupplier(&types.Supplier{
// 					CompanyId:  payload.CompanyId,
// 					SupplierId: "1234",
// 					Name:       "Test Supplier",
// 				})
//
// 				Expect(err).NotTo(HaveOccurred())
// 			})
//
// 			It("should fail if supplier already exists", func() {
// 				err := db.CreateSupplier(&types.Supplier{
// 					CompanyId:  payload.CompanyId,
// 					SupplierId: "1234",
// 					Name:       "Test Supplier",
// 				})
//
// 				log.Println(err.(*pgconn.PgError))
//
// 				Expect(err).To(HaveOccurred())
// 				Expect(err.(*pgconn.PgError).Code).To(Equal("23505"))
// 			})
//
// 			It("should create multiple suppliers", func() {
// 				for i := 0; i < 10; i++ {
// 					err := db.CreateSupplier(&types.Supplier{
// 						CompanyId:  payload.CompanyId,
// 						SupplierId: fmt.Sprintf("%d", i),
// 						Name:       fmt.Sprintf("Supplier %d", i),
// 					})
// 					Expect(err).NotTo(HaveOccurred())
// 				}
// 			})
//
// 			It("should get all suppliers", func() {
// 				suppliers, err := db.GetAllSuppliers(payload.CompanyId, "")
//
// 				Expect(err).NotTo(HaveOccurred())
// 				Expect(len(suppliers)).To(Equal(11))
// 			})
//
// 			It("should update a supplier", func() {
// 				suppliers, err := db.GetAllSuppliers(payload.CompanyId, "")
// 				Expect(err).NotTo(HaveOccurred())
//
// 				supplier, err := db.GetOneSupplier(suppliers[0].ID, payload.CompanyId)
// 				Expect(err).NotTo(HaveOccurred())
//
// 				Expect(supplier.Name).To(Equal("Supplier 0"))
//
// 				err = db.UpdateSupplier(&types.Supplier{
// 					CompanyId:  payload.CompanyId,
// 					ID:         supplier.ID,
// 					SupplierId: "1564",
// 					Name:       "Test Supplier 2",
// 					ContactName: sql.NullString{
// 						Valid:  true,
// 						String: "Test Contact Name",
// 					},
// 					ContactEmail: sql.NullString{
// 						Valid:  true,
// 						String: "contact@test.com",
// 					},
// 					ContactPhone: sql.NullString{
// 						Valid:  true,
// 						String: "123456789",
// 					},
// 				})
//
// 				Expect(err).NotTo(HaveOccurred())
//
// 				supplier, err = db.GetOneSupplier(suppliers[0].ID, payload.CompanyId)
// 				Expect(err).NotTo(HaveOccurred())
// 				Expect(supplier.Name).To(Equal("Test Supplier 2"))
// 				Expect(supplier.ContactName.String).To(Equal("Test Contact Name"))
// 				Expect(supplier.ContactEmail.String).To(Equal("contact@test.com"))
// 				Expect(supplier.ContactPhone.String).To(Equal("123456789"))
// 				Expect(supplier.SupplierId).To(Equal("1564"))
// 			})
// 		})
//
// 		When("Categories", func() {
// 			It("should create a category", func() {
// 				err := db.CreateCategory(types.Category{
// 					CompanyId: payload.CompanyId,
// 					Name:      "Materiales",
// 				})
//
// 				Expect(err).NotTo(HaveOccurred())
// 			})
//
// 			It("should fail if category already exists", func() {
// 				err := db.CreateCategory(types.Category{
// 					CompanyId: payload.CompanyId,
// 					Name:      "Materiales",
// 				})
//
// 				Expect(err).To(HaveOccurred())
// 			})
//
// 			It("should create other categories", func() {
// 				err := db.CreateCategory(types.Category{
// 					CompanyId: payload.CompanyId,
// 					Name:      "Mano de Obra",
// 				})
//
// 				Expect(err).NotTo(HaveOccurred())
// 			})
//
// 			It("should get all categories", func() {
// 				categories, err := db.GetAllCategories(payload.CompanyId)
//
// 				Expect(err).NotTo(HaveOccurred())
// 				Expect(len(categories)).To(Equal(2))
// 			})
//
// 			It("should update a category", func() {
// 				categories, err := db.GetAllCategories(payload.CompanyId)
//
// 				Expect(err).NotTo(HaveOccurred())
// 				Expect(len(categories)).To(Equal(2))
//
// 				category, err := db.GetCategory(categories[1].Id, payload.CompanyId)
//
// 				Expect(err).NotTo(HaveOccurred())
// 				Expect(category.Name).To(Equal("Materiales"))
//
// 				err = db.UpdateCategory(types.Category{
// 					CompanyId: payload.CompanyId,
// 					Id:        category.Id,
// 					Name:      "Updated Materials",
// 				})
//
// 				Expect(err).NotTo(HaveOccurred())
//
// 				category, err = db.GetCategory(categories[1].Id, payload.CompanyId)
//
// 				Expect(err).NotTo(HaveOccurred())
// 				Expect(category.Name).To(Equal("Updated Materials"))
// 			})
//
// 		})
// 	})
// })
