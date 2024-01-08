package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views"
	"bca-go-final/internal/views/bca"
	"bca-go-final/internal/views/derrors"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func (s *Server) LoginView(w http.ResponseWriter, r *http.Request) {
	store.Options.HttpOnly = true
	store.Options.Secure = true
	store.Options.SameSite = http.SameSiteStrictMode

	session, _ := store.Get(r, "bca")
	resp := make(map[string]string)
	switch r.Method {
	case http.MethodPost:
		l := &types.Login{}
		_ = l
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
			token, err := s.DB.Login(l)
			if err != nil {
				resp["error"] = "credenciales inválidas"
			} else {
				session.Values["bca"] = token
				session.Save(r, w)
				resp["token"] = token
				http.Redirect(w, r, "/bca", http.StatusSeeOther)
			}
		}

	case http.MethodGet:
		resp = nil

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp["method"] = r.Method
		resp["url"] = r.RequestURI
		component := derrors.NotImplemented(resp)
		component.Render(r.Context(), w)
	}

	if len(resp) == 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	component := views.LoginView(resp)
	component.Render(r.Context(), w)
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "bca")
	session.Values["bca"] = nil
	session.Save(r, w)
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) BcaView(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	component := bca.LandingPage(ctx.Name)
	component.Render(r.Context(), w)
}
