package models

import (
	"github.com/google/uuid"
)

type EnrollmentDelete struct {
	StudentID uuid.UUID `form:"studentId" validate:"required"`
	CourseID  uuid.UUID `form:"courseId" validate:"required"`
}
