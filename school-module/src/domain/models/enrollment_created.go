package models

import (
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/google/uuid"
)

type EnrollmentCreatedStudent struct {
	ID uuid.UUID `json:"id"`
}

type EnrollmentCreatedCourse struct {
	ID uuid.UUID `json:"id"`
}

type EnrollmentCreated struct {
	Student      EnrollmentCreatedStudent `json:"student"`
	Course       EnrollmentCreatedCourse  `json:"course"`
	Installments uint8                    `json:"installments"`
	Status       enums.EnrollmentStatus   `json:"status"`
	CreatedAt    time.Time                `json:"createdAt"`
}
