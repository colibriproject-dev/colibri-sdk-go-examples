package models

import (
	"github.com/google/uuid"
)

type CourseDelete struct {
	ID uuid.UUID `json:"id"`
}
