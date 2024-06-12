package utils

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/alcb1310/bca/internals/types"
)

var ja = jwtauth.New("HS256", []byte(os.Getenv("SECRET")), nil)

// GenerateJWT: Generates a new JWT token based on a given user
func GenerateJWT(u types.User) (string, error) {
	_, tokenString, err := ja.Encode(map[string]interface{}{"id": u.ID, "email": u.Email, "company_id": u.CompanyID, "name": u.Name})

	return tokenString, err
}

// ValidateJWTToken: Validates a given JWT token
func ValidateJWTToken(w http.ResponseWriter, r *http.Request, tokenString string) (types.User, error) {
	user := types.User{}
	token, err := jwtauth.VerifyToken(ja, tokenString)
	if err != nil {
		slog.Error("AuthMiddleware", "error", err)
		return user, err
	}

	u, ok := token.(jwt.Token)

	if !ok {
		slog.Error("AuthMiddleware", "error", err)
		return user, err
	}

	claims, err := u.AsMap(context.Background())
	if err != nil {
		slog.Error("AuthMiddleware", "error", err)
		return user, err
	}

	user.ID = uuid.MustParse(claims["id"].(string))
	user.Email = claims["email"].(string)
	user.Name = claims["name"].(string)
	user.CompanyID = uuid.MustParse(claims["company_id"].(string))

	return user, nil
}
