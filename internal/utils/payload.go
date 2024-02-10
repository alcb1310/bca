package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"bca-go-final/internal/types"
)

func GetMyPaload(r *http.Request) (types.ContextPayload, error) {
	ctx := r.Context()
	val := ctx.Value("token")

	x, ok := val.([]byte)
	if !ok {
		return types.ContextPayload{}, errors.New("Unable to load context")
	}
	var p types.ContextPayload
	err := json.Unmarshal(x, &p)
	if err != nil {
		return types.ContextPayload{}, errors.New("Unable to parse context")
	}
	return p, nil
}
