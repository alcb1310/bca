package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/views/bca/transaction/partials/details"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) DetailsTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	id := mux.Vars(r)["invoiceId"]
	parsedId, _ := uuid.Parse(id)

	det, err := s.DB.GetAllDetails(parsedId, ctx.CompanyId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		// TODO: Validate data
		// TODO: save to the database
		// TODO: return details table and display new total
	}

	component := details.InvoiceDetailsTable(det)
	component.Render(r.Context(), w)
}

func (s *Server) DetailsAdd(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)

	budgetItems := []types.Select{}
	bi := s.DB.GetBudgetItemsByAccumulate(ctx.CompanyId, false)

	for _, v := range bi {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		budgetItems = append(budgetItems, x)
	}

	component := details.EditDetails(budgetItems)
	component.Render(r.Context(), w)
}
