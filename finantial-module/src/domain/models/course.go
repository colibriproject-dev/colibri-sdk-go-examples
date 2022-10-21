package models

import (
	"github.com/google/uuid"
)

type Course struct {
	ID    uuid.UUID `json:"id"`
	Value float64   `json:"value"`
}
