package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"bca-go-final/internal/views"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.Use(s.authVerify)

	r.HandleFunc("/", s.HelloWorldHandler)
	r.HandleFunc("/api/login", s.Login)       // Fully tested
	r.HandleFunc("/api/register", s.Register) // Fully tested
	r.HandleFunc("/bca/dummy", s.loadDummyDataHandler)

	// views

	r.HandleFunc("/login", s.LoginView)   // Fully tested
	r.HandleFunc("/bca", s.BcaView)       // Fully tested
	r.HandleFunc("/bca/logout", s.Logout) // Fully tested

	r.HandleFunc("/bca/transacciones/presupuesto", s.Budget) // Fully tested
	r.HandleFunc("/bca/transacciones/facturas", s.Invoice)   // Fully tested
	r.HandleFunc("/bca/transacciones/facturas/crear", s.InvoiceAdd)
	r.HandleFunc("/bca/transacciones/cierre", s.Closure) // Fully tested

	r.HandleFunc("/bca/reportes/actual", s.Actual)                                               // Fully tested
	r.HandleFunc("/bca/reportes/actual/generar", s.ActualGenerate)                               // Fully tested
	r.HandleFunc("/bca/reportes/cuadre", s.Balance)                                              // Fully tested
	r.HandleFunc("/bca/reportes/historico", s.Historic)                                          // Fully tested
	r.HandleFunc("/bca/reportes/gastado", s.Spent)                                               // Fully tested
	r.HandleFunc("/bca/reportes/gastado/{projectId}/{budgetItemId}/{date}", s.SpentByBudgetItem) // Trouble reading URL parameters

	r.HandleFunc("/bca/configuracion/partidas", s.BudgetItems)   // Fully tested
	r.HandleFunc("/bca/configuracion/proveedores", s.Suppliers)  // Fully tested
	r.HandleFunc("/bca/configuracion/proyectos", s.Projects)     // Fully tested
	r.HandleFunc("/bca/configuracion/categorias", s.Categories)  // Fully tested
	r.HandleFunc("/bca/configuracion/materiales", s.Materiales)  // Fully tested
	r.HandleFunc("/bca/configuracion/rubros", s.Rubros)          // Fully tested
	r.HandleFunc("/bca/configuracion/rubros/crear", s.RubrosAdd) // Fully tested

	r.HandleFunc("/bca/user/perfil", s.Profile) // Fully tested
	r.HandleFunc("/bca/user/admin", s.Admin)
	r.HandleFunc("/bca/user/cambio", s.ChangePassword)

	r.HandleFunc("/bca/costo-unitario/cantidades", s.UnitQuantity)
	r.HandleFunc("/bca/costo-unitario/analisis", s.UnitAnalysis)

	// excel
	r.HandleFunc("/bca/reportes/excel/cuadre", s.BalanceExcel)
	r.HandleFunc("/bca/reportes/excel/actual", s.ActualExcel)
	r.HandleFunc("/bca/reportes/excel/historico", s.HistoricExcel)
	r.HandleFunc("/bca/reportes/excel/gastado", s.SpentExcel)

	// partials

	r.HandleFunc("/bca/partials/users", s.UsersTable)
	r.HandleFunc("/bca/partials/users/add", s.UserAdd)
	r.HandleFunc("/bca/partials/users/edit/{id}", s.UserEdit)
	r.HandleFunc("/bca/partials/users/{id}", s.SingleUser)

	r.HandleFunc("/bca/partials/projects", s.ProjectsTable)
	r.HandleFunc("/bca/partials/projects/add", s.ProjectAdd)
	r.HandleFunc("/bca/partials/projects/edit/{id}", s.ProjectEditSave)
	r.HandleFunc("/bca/partials/projects/{id}", s.ProjectEdit)

	r.HandleFunc("/bca/partials/suppliers", s.SuppliersTable)
	r.HandleFunc("/bca/partials/suppliers/add", s.SupplierAdd)
	r.HandleFunc("/bca/partials/suppliers/edit/{id}", s.SuppliersEditSave)
	r.HandleFunc("/bca/partials/suppliers/{id}", s.SuppliersEdit)

	r.HandleFunc("/bca/partials/budget-item", s.BudgetItemsTable)
	r.HandleFunc("/bca/partials/budget-item/add", s.BudgetItemAdd)
	r.HandleFunc("/bca/partials/budget-item/{id}", s.BudgetItemEdit)

	r.HandleFunc("/bca/partials/budgets", s.BudgetsTable)
	r.HandleFunc("/bca/partials/budgets/add", s.BudgetAdd)
	r.HandleFunc("/bca/partials/budgets/{projectId}/{budgetItemId}", s.BudgetEdit)

	r.HandleFunc("/bca/partials/invoices", s.InvoicesTable)
	r.HandleFunc("/bca/partials/invoices/{id}", s.InvoiceEdit)

	r.HandleFunc("/bca/partials/invoices/{invoiceId}/details", s.DetailsTable)
	r.HandleFunc("/bca/partials/invoices/{invoiceId}/details/add", s.DetailsAdd)
	r.HandleFunc("/bca/partials/invoices/{invoiceId}/details/{budgetItemId}", s.DetailsEdit)

	r.HandleFunc("/bca/partials/categories", s.CategoriesTable)
	r.HandleFunc("/bca/partials/categories/add", s.CategoryAdd)
	r.HandleFunc("/bca/partials/categories/{id}", s.EditCategory)

	r.HandleFunc("/bca/partials/materiales", s.MaterialsTable)
	r.HandleFunc("/bca/partials/materiales/add", s.MaterialsAdd)
	r.HandleFunc("/bca/partials/materiales/{id}", s.MaterialsEdit)

	r.HandleFunc("/bca/partials/rubros", s.RubrosTable)
	r.HandleFunc("/bca/partials/rubros/{id}", s.MaterialsByItem)
	r.HandleFunc("/bca/partials/rubros/{id}/material", s.MaterialByItemForm)
	r.HandleFunc("/bca/partials/rubros/{id}/material/{materialId}", s.MaterialItemsOperations)

	r.HandleFunc("/bca/partials/cantidades", s.CantidadesTable)
	r.HandleFunc("/bca/partials/cantidades/add", s.CantidadesAdd)
	r.HandleFunc("/bca/partials/cantidades/{id}", s.CantidadesEdit)

	r.HandleFunc("/bca/partials/analisis", s.AnalysisTable)

	// This should be the last route for static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	component := views.WelcomeView()
	component.Render(r.Context(), w)
}
