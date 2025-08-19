package models

import (
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/types"
)

type EnrollmentPage *types.Page[Enrollment]

type EnrollmentPageParams struct {
	Page        uint16 `form:"page" validate:"required"`
	Size        uint16 `form:"pageSize" validate:"required"`
	StudentName string `form:"studentName"`
	CourseName  string `form:"courseName"`
}
