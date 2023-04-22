package models

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Birthday  time.Time `json:"birthday" validate:"required"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
}

type StudentCreateUpdateDTO struct {
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required"`
	Birthday time.Time `json:"birthday" validate:"required"`
}

type StudentParams struct {
	Name string `form:"name"`
}
