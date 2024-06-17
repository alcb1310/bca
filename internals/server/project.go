package server

import (
	"net/http"

	"github.com/alcb1310/bca/externals/views/bca/parameters/projects"
	"github.com/alcb1310/bca/internals/types"
)

func (s *BCAService) ProjectsPage(w http.ResponseWriter, r *http.Request) error {
	return renderPage(w, r, projects.ProjectLanding())
}

func (s *BCAService) ProjectsTable(w http.ResponseWriter, r *http.Request) error {
	user, _ := getUserFromContext(r)
	retrievedProjects := s.DB.GetAllProjects(user.CompanyID)

	return renderPage(w, r, projects.ProjectTable(retrievedProjects))
}

func (s *BCAService) ProjectsForm(w http.ResponseWriter, r *http.Request) error {
	return renderPage(w, r, projects.ProjectForm())
}
