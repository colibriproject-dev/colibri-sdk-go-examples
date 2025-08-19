package repositories

import (
	"testing"
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/stretchr/testify/assert"
)

var (
	enrollmentDatasets   = []string{"clear-data.sql", "courses-insert.sql", "students-insert.sql", "enrollments-insert.sql"}
	enrollmentRepository = repositories.NewEnrollmentsDBRepository()
	enrollmentMockData   = []models.Enrollment{
		{
			Student:      studentMockData[2],
			Course:       courseMockData[2],
			Installments: 1,
			Status:       enums.INADIMPLENTE,
			CreatedAt:    time.Date(2023, time.September, 16, 18, 0, 0, 0, time.FixedZone("", 0)),
		},
		{
			Student:      studentMockData[1],
			Course:       courseMockData[1],
			Installments: 5,
			Status:       enums.ADIMPLENTE,
			CreatedAt:    time.Date(2023, time.September, 15, 16, 0, 0, 0, time.FixedZone("", 0)),
		},
		{
			Student:      studentMockData[0],
			Course:       courseMockData[0],
			Installments: 10,
			Status:       enums.ADIMPLENTE,
			CreatedAt:    time.Date(2023, time.September, 14, 10, 0, 0, 0, time.FixedZone("", 0)),
		},
	}
)

func TestEnrollmentRepository_FindAllPaginated(t *testing.T) {
	assert.NoError(t, pc.Dataset(basePath, enrollmentDatasets...))

	t.Run("Should return enrollments page when exists into enrollments table", func(t *testing.T) {
		expected := enrollmentMockData

		result, err := enrollmentRepository.FindAllPaginated(ctx, &models.EnrollmentPageParams{Page: 1, Size: 10})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result.Items)
		assert.EqualValues(t, uint64(3), result.TotalItems)
	})

	t.Run("Should return enrollments page filtered by student name when exists into enrollments table", func(t *testing.T) {
		expected := []models.Enrollment{enrollmentMockData[2]}

		result, err := enrollmentRepository.FindAllPaginated(ctx, &models.EnrollmentPageParams{Page: 1, Size: 10, StudentName: studentMockData[0].Name})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result.Items)
		assert.EqualValues(t, uint64(1), result.TotalItems)
	})

	t.Run("Should return enrollments page filtered by course name when exists into enrollments table", func(t *testing.T) {
		expected := []models.Enrollment{enrollmentMockData[1]}

		result, err := enrollmentRepository.FindAllPaginated(ctx, &models.EnrollmentPageParams{Page: 1, Size: 10, CourseName: courseMockData[1].Name})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result.Items)
		assert.EqualValues(t, uint64(1), result.TotalItems)
	})
}

func TestEnrollmentRepository_Insert(t *testing.T) {
	assert.NoError(t, pc.Dataset(basePath, enrollmentDatasets...))

	t.Run("Should return error when enrollment exists into enrollments table", func(t *testing.T) {
		result, err := enrollmentRepository.Insert(ctx, &models.EnrollmentCreate{
			StudentID:    enrollmentMockData[0].Student.ID,
			CourseID:     enrollmentMockData[0].Course.ID,
			Installments: enrollmentMockData[0].Installments,
			Status:       enrollmentMockData[0].Status,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should insert and return enrollment when not exists into enrollments table", func(t *testing.T) {
		expected := &models.EnrollmentCreated{
			Student:      models.EnrollmentCreatedStudent{ID: enrollmentMockData[0].Student.ID},
			Course:       models.EnrollmentCreatedCourse{ID: enrollmentMockData[1].Course.ID},
			Installments: 10,
			Status:       enums.ADIMPLENTE,
		}
		model := &models.EnrollmentCreate{
			StudentID:    expected.Student.ID,
			CourseID:     expected.Course.ID,
			Installments: expected.Installments,
			Status:       expected.Status,
		}

		result, err := enrollmentRepository.Insert(ctx, model)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expected.Student.ID, result.Student.ID)
		assert.EqualValues(t, expected.Course.ID, result.Course.ID)
		assert.EqualValues(t, expected.Installments, result.Installments)
		assert.EqualValues(t, expected.Status, result.Status)
		assert.NotEmpty(t, result.CreatedAt)
	})
}

func TestEnrollmentRepository_Delete(t *testing.T) {
	assert.NoError(t, pc.Dataset(basePath, enrollmentDatasets...))

	t.Run("Should delete enrollment when exists into enrollments table", func(t *testing.T) {
		errDelete := enrollmentRepository.Delete(ctx, enrollmentMockData[0].Student.ID, enrollmentMockData[0].Course.ID)
		result, resultErr := enrollmentRepository.FindByStudentIdAndCourseId(ctx,
			enrollmentMockData[0].Student.ID,
			enrollmentMockData[0].Course.ID,
		)

		assert.NoError(t, errDelete)
		assert.NoError(t, resultErr)
		assert.Nil(t, result)
	})
}

func TestEnrollmentRepository_UpdateStatus(t *testing.T) {
	assert.NoError(t, pc.Dataset(basePath, enrollmentDatasets...))

	t.Run("Should update enrollment status when exists into enrollments table", func(t *testing.T) {
		errUpdate := enrollmentRepository.UpdateStatus(ctx, &models.EnrollmentUpdateStatus{
			StudentID: enrollmentMockData[0].Student.ID,
			CourseID:  enrollmentMockData[0].Course.ID,
			Status:    enums.ADIMPLENTE,
		})
		result, err := enrollmentRepository.FindByStudentIdAndCourseId(ctx,
			enrollmentMockData[0].Student.ID,
			enrollmentMockData[0].Course.ID,
		)

		assert.NoError(t, errUpdate)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.CreatedAt)
		assert.EqualValues(t, enums.ADIMPLENTE, result.Status)
	})
}
