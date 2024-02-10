package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
)

func (s *Server) authVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/bca") {
			next.ServeHTTP(w, r)
			return
		}
		secretKey := os.Getenv("SECRET")
		maker, err := utils.NewJWTMaker(secretKey)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		session, _ := store.Get(r, "bca")
		token, ok := session.Values["bca"].(string)

		if !ok || token == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenData, err := maker.VerifyToken(token)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		marshalStr, _ := json.Marshal(tokenData)
		ctx := r.Context()
		ctx = context.WithValue(r.Context(), "token", marshalStr)
		r = r.Clone(ctx)

		if !s.DB.IsLoggedIn(token, tokenData.ID) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		u := &types.User{
			Id:        tokenData.ID,
			Name:      tokenData.Name,
			Email:     tokenData.Email,
			CompanyId: tokenData.CompanyId,
			RoleId:    tokenData.Role,
		}

		token, err = utils.GenerateToken(*u)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		session.Values["bca"] = token
		session.Save(r, w)

		if err := s.DB.RegenerateToken(token, u.Id); err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
