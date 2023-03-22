package models

import (
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/types"
	"github.com/google/uuid"
)

type EnrollmentPage *types.Page[Enrollment]

type Enrollment struct {
	Student      Student          `json:"student"`
	Course       Course           `json:"course"`
	Installments uint8            `json:"installments"`
	Status       EnrollmentStatus `json:"status"`
	CreatedAt    time.Time        `json:"createdAt"`
}

type EnrollmentCreateDTO struct {
	StudentID    uuid.UUID `json:"studentId" validate:"required"`
	CourseID     uuid.UUID `json:"courseId" validate:"required"`
	Installments uint8     `json:"installments" validate:"required"`
}

type EnrollmentDeleteParamsDTO struct {
	StudentID uuid.UUID `form:"studentId" validate:"required"`
	CourseID  uuid.UUID `form:"courseId" validate:"required"`
}

type EnrollmentPageParamsDTO struct {
	Page        uint16 `form:"page" validate:"required"`
	Size        uint16 `form:"pageSize" validate:"required"`
	StudentName string `form:"studentName"`
	CourseName  string `form:"courseName"`
}

type EnrollmentFilters struct {
	StudentName string
	CourseName  string
}

func (p *EnrollmentPageParamsDTO) ToFilters() *EnrollmentFilters {
	return &EnrollmentFilters{
		CourseName:  p.CourseName,
		StudentName: p.StudentName,
	}
}

func (p *EnrollmentPageParamsDTO) ToPageRequest() *types.PageRequest {
	return types.NewPageRequest(p.Page, p.Size, []types.Sort{types.NewSort(types.DESC, "e.created_at")})
}
