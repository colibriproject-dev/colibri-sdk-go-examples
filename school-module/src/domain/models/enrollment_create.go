package models

import (
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/google/uuid"
)

type EnrollmentCreate struct {
	StudentID    uuid.UUID              `json:"studentId" validate:"required"`
	CourseID     uuid.UUID              `json:"courseId" validate:"required"`
	Installments uint8                  `json:"installments" validate:"required"`
	Status       enums.EnrollmentStatus `json:"-"`
}
