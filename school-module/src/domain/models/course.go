package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
}
