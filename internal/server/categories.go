package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/settings/partials"
)

func (s *Server) CategoriesTable(w http.ResponseWriter, r *http.Request) {
	var err error
	ctxPayload, _ := utils.GetMyPaload(r)

	if r.Method == http.MethodPost {
		r.ParseForm()
		c := types.Category{
			Name:      r.Form.Get("name"),
			CompanyId: ctxPayload.CompanyId,
		}
		err = s.DB.CreateCategory(c)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(fmt.Sprintf("La categoria %s ya existe", c.Name)))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}

	categories, _ := s.DB.GetAllCategories(ctxPayload.CompanyId)
	component := partials.CategoriesTable(categories)

	component.Render(r.Context(), w)
}

func (s *Server) CategoryAdd(w http.ResponseWriter, r *http.Request) {
	component := partials.EditCategory(nil)
	component.Render(r.Context(), w)
}

func (s *Server) EditCategory(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)
	parsedId, _ := utils.ValidateUUID(mux.Vars(r)["id"], "categoria")
	c, _ := s.DB.GetCategory(parsedId, ctxPayload.CompanyId)

	switch r.Method {
	case http.MethodGet:
		component := partials.EditCategory(&c)
		component.Render(r.Context(), w)

	case http.MethodPut:
		r.ParseForm()
		cat := types.Category{
			Id:        parsedId,
			Name:      r.Form.Get("name"),
			CompanyId: ctxPayload.CompanyId,
		}

		if err := s.DB.UpdateCategory(cat); err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(fmt.Sprintf("La categoria %s ya existe", cat.Name)))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		categories, _ := s.DB.GetAllCategories(ctxPayload.CompanyId)
		component := partials.CategoriesTable(categories)

		component.Render(r.Context(), w)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
