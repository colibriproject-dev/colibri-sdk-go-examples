package models

import "github.com/google/uuid"

type CourseUpdate struct {
	Name  string    `json:"name" validate:"required"`
	Value float64   `json:"value" validate:"required"`
	ID    uuid.UUID `json:"-"`
}
