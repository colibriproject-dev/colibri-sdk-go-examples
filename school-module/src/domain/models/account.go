package models

import (
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID              `json:"id" validate:"required"`
	StudentID    uuid.UUID              `json:"studentId" validate:"required"`
	CourseID     uuid.UUID              `json:"courseId" validate:"required"`
	Installments uint8                  `json:"installments" validate:"required"`
	Value        float64                `json:"value" validate:"required"`
	Status       enums.EnrollmentStatus `json:"status" validate:"required,oneOfEnrollmentStatus"`
	CreatedAt    time.Time              `json:"createdAt" validate:"required"`
}

func (m *Account) ToEnrollmentUpdateStatus() *EnrollmentUpdateStatus {
	return &EnrollmentUpdateStatus{
		StudentID: m.StudentID,
		CourseID:  m.CourseID,
		Status:    m.Status,
	}
}
