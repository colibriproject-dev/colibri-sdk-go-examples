package models

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Birthday  time.Time `json:"birthday"`
	CreatedAt time.Time `json:"createdAt"`
}

type StudentCreateUpdateDTO struct {
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required"`
	Birthday time.Time `json:"birthday" validate:"required"`
}

type StudentParams struct {
	Name string `schema:"name"`
}
