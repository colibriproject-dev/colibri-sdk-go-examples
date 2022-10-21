package models

import (
	"time"
)

type Enrollment struct {
	Student      Student          `json:"student"`
	Course       Course           `json:"course"`
	Installments uint8            `json:"installments"`
	Status       EnrollmentStatus `json:"status"`
	CreatedAt    time.Time        `json:"createdAt"`
}

type EnrollmentCreateDTO struct {
	Student      Student `json:"student" validate:"required"`
	Course       Course  `json:"course" validate:"required"`
	Installments uint8   `json:"installments" validate:"required"`
}
