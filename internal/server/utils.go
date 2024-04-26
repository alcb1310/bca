package server

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

type Result struct {
	n string
	s []types.Select
}

func (s *Server) getSelect(query string, companyId uuid.UUID) []types.Select {
	sel := make([]types.Select, 0)

	switch query {
	case "levels":
		sel = s.DB.Levels(companyId)

	case "projects":
		p := s.DB.GetActiveProjects(companyId, true)
		for _, v := range p {
			x := types.Select{
				Key:   v.ID.String(),
				Value: v.Name,
			}
			sel = append(sel, x)
		}

	case "suppliers":
		sx, _ := s.DB.GetAllSuppliers(companyId, "")
		for _, v := range sx {
			x := types.Select{
				Key:   v.ID.String(),
				Value: v.Name,
			}
			sel = append(sel, x)
		}

	case "budgetitems":
		b := s.DB.GetBudgetItemsByAccumulate(companyId, false)
		for _, v := range b {
			x := types.Select{
				Key:   v.ID.String(),
				Value: v.Name,
			}
			sel = append(sel, x)
		}

	case "categories":
		c, _ := s.DB.GetAllCategories(companyId)
		for _, v := range c {
			x := types.Select{
				Key:   v.Id.String(),
				Value: v.Name,
			}
			sel = append(sel, x)
		}

	case "materials":
		m := s.DB.GetAllMaterials(companyId)
		for _, v := range m {
			x := types.Select{
				Key:   v.Id.String(),
				Value: v.Name,
			}
			sel = append(sel, x)
		}

	case "rubros":
		r, _ := s.DB.GetAllRubros(companyId)
		for _, v := range r {
			x := types.Select{
				Key:   v.Id.String(),
				Value: v.Name,
			}
			sel = append(sel, x)
		}
	}

	return sel
}

func (s *Server) returnAllSelects(query []string, companyId uuid.UUID, flags ...bool) map[string][]types.Select {
	results := make(map[string][]types.Select)
	resultChannel := make(chan Result)
	defer close(resultChannel)

	for _, value := range query {
		go func(v string) {
			results[v] = s.getSelect(v, companyId)
			resultChannel <- Result{v, results[v]}
		}(value)
	}

	for i := 0; i < len(query); i++ {
		r := <-resultChannel
		results[r.n] = r.s
	}

	return results
}
