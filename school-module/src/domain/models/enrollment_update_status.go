package models

import (
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/google/uuid"
)

type EnrollmentUpdateStatus struct {
	StudentID uuid.UUID
	CourseID  uuid.UUID
	Status    enums.EnrollmentStatus
}
