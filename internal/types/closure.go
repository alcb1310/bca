package types

import (
	"time"

	"github.com/google/uuid"
)

type Closure struct {
	ProjectId uuid.UUID `json:"project_id"`
	Date      time.Time `json:"date"`
}
