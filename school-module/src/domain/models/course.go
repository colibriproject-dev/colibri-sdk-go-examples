package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Value     float64   `json:"value" validate:"required"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
}

type CourseCreateUpdateDTO struct {
	Name  string  `json:"name" validate:"required"`
	Value float64 `json:"value" validate:"required"`
}
