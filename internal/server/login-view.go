package server

import (
	"net/http"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views"
	"bca-go-final/internal/views/bca"
)

func (s *Server) LoginView(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	l := &types.Login{}
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		component := views.LoginView(resp)
		component.Render(r.Context(), w)
		return
	}

	l.Email = r.PostFormValue("email")
	l.Password = r.PostFormValue("password")

	if !utils.IsValidEmail(l.Email) {
		resp["error"] = "credenciales inválidas"
	}

	if l.Password == "" {
		resp["error"] = "credenciales inválidas"
	}

	if len(resp) == 0 {
		_, u, err := s.DB.Login(l)
		if err != nil {
			resp["error"] = "credenciales inválidas"
		} else {
			_, tokenString, _ := s.TokenAuth.Encode(map[string]interface{}{"id": u.Id, "name": u.Name, "email": u.Email, "company_id": u.CompanyId, "role": u.RoleId})
			http.SetCookie(w, &http.Cookie{
				Name:  "jwt",
				Value: tokenString,
				Path:  "/",
			})

			resp["token"] = tokenString
			http.Redirect(w, r, "/bca", http.StatusSeeOther)
			return
		}
	}

	if len(resp) == 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	component := views.LoginView(resp)
	component.Render(r.Context(), w)
}

func (s *Server) DisplayLogin(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)

	w.WriteHeader(http.StatusOK)
	component := views.LoginView(resp)
	component.Render(r.Context(), w)
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	// session, _ := store.Get(r, "bca")
	// session.Values["bca"] = nil
	// session.Save(r, w)
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) BcaView(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	component := bca.LandingPage(ctx.Name)
	component.Render(r.Context(), w)
}
