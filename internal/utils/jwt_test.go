package utils_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
)

var invalidJWTToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o"
var expiredJWTToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjVkYzAyZGE1LTFkMWItNDcwZi04NmJhLWRkYTg0ZTY0MWZiNiIsImVtYWlsIjoiYUBhLmNvbSIsImNvbXBhbnlfaWQiOiI5NzAyNTIxNi1jMzRjLTRjMDYtYTc0ZC04NDdiZDcwOThmZjQiLCJuYW1lIjoiQW5kcmVzIiwicm9sZSI6ImFkbWluIiwiaXNfbG9nZ2VkX2luIjp0cnVlLCJpc3N1ZWRfYXQiOiIyMDI0LTA1LTAxVDE4OjQ0OjEyLjU1MTMxNS0wNTowMCIsImV4cGlyZWRfYXQiOiIyMDI0LTA1LTAxVDE4OjQ1OjEyLjU1MTMxNS0wNTowMCJ9.ErdkrSpz_RFOi1IcILJSgTG7CWf_ihG66O0eLaQtywk"

func TestGenerateJWT(t *testing.T) {
	u := types.User{
		Id:        uuid.New(),
		Email:     "a@a.com",
		Name:      "Andres",
		RoleId:    "admin",
		CompanyId: uuid.New(),
	}

	t.Run("should generate JWT", func(t *testing.T) {
		maker, err := utils.NewJWTMaker("supersecretkey")
		assert.Nil(t, err)

		token, err := maker.CreateToken(u, time.Minute)
		assert.Nil(t, err)

		assert.NotEmpty(t, token)
	})

	t.Run("should verify the token", func(t *testing.T) {
		maker, err := utils.NewJWTMaker("supersecretkey")
		assert.Nil(t, err)

		token, err := maker.CreateToken(u, time.Minute)
		assert.Nil(t, err)

		payload, err := maker.VerifyToken(token)
		assert.Nil(t, err)

		assert.Equal(t, u.Id, payload.ID)
		assert.Equal(t, u.Email, payload.Email)
		assert.Equal(t, u.CompanyId, payload.CompanyId)
		assert.Equal(t, u.Name, payload.Name)
		assert.Equal(t, u.RoleId, payload.Role)
	})

	t.Run("should return error if short secret key", func(t *testing.T) {
		_, err := utils.NewJWTMaker("short")
		assert.NotNil(t, err)

		assert.Equal(t, "invalid key size: must be at least 8 characters", err.Error())
	})

	t.Run("should return error if invalid token", func(t *testing.T) {
		maker, err := utils.NewJWTMaker("supersecretkey")
		assert.Nil(t, err)

		_, err = maker.VerifyToken(invalidJWTToken)
		assert.NotNil(t, err)

		assert.Equal(t, "invalid token", err.Error())
	})

	t.Run("should return error if expired token", func(t *testing.T) {
		maker, err := utils.NewJWTMaker("supersecretkey")
		assert.Nil(t, err)

		_, err = maker.VerifyToken(expiredJWTToken)
		assert.NotNil(t, err)

		assert.Equal(t, "token has expired", err.Error())
	})
}

func TestGenerateToken(t *testing.T) {
	u := types.User{
		Id:        uuid.New(),
		Email:     "a@a.com",
		Name:      "Andres",
		RoleId:    "admin",
		CompanyId: uuid.New(),
	}

	t.Run("long secret", func(t *testing.T) {
		t.Setenv("SECRET", "supersecretkey")

		token, err := utils.GenerateToken(u)

		assert.Nil(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("short secret", func(t *testing.T) {
		t.Setenv("SECRET", "short")

		token, err := utils.GenerateToken(u)

		assert.NotNil(t, err)
		assert.Empty(t, token)

		assert.Equal(t, "invalid key size: must be at least 8 characters", err.Error())
	})
}
