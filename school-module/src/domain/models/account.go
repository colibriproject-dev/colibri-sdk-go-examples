package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID        `json:"id"`
	StudentID    uuid.UUID        `json:"studentId"`
	CourseID     uuid.UUID        `json:"courseId"`
	Installments uint8            `json:"installments"`
	Value        float64          `json:"value"`
	Status       EnrollmentStatus `json:"status"`
	CreatedAt    time.Time        `json:"createdAt"`
}
