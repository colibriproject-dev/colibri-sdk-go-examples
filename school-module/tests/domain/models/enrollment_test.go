package models

import (
	"testing"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/types"
	"github.com/stretchr/testify/assert"
)

func TestEnrollment(t *testing.T) {
	t.Run("Should test page params conversions", func(t *testing.T) {
		pageParams := models.EnrollmentPageParamsDTO{
			Page:        1,
			Size:        10,
			StudentName: "test student",
			CourseName:  "test course",
		}

		filters := &models.EnrollmentFilters{
			StudentName: pageParams.StudentName,
			CourseName:  pageParams.CourseName,
		}

		pageRequest := types.NewPageRequest(
			pageParams.Page,
			pageParams.Size,
			[]types.Sort{types.NewSort(types.DESC, "e.created_at")},
		)

		assert.Equal(t, filters, pageParams.ToFilters())
		assert.Equal(t, pageRequest, pageParams.ToPageRequest())
	})
}
