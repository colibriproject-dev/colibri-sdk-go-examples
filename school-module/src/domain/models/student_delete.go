package models

import (
	"github.com/google/uuid"
)

type StudentDelete struct {
	ID uuid.UUID `json:"id"`
}
