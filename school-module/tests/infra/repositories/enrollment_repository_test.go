package repositories

import (
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	enrollmentRepository = repositories.NewEnrollmentDBRepository()
	enrollmentMockData   = []models.Enrollment{
		{
			Student:      studentMockData[2],
			Course:       courseMockData[2],
			Installments: 1,
			Status:       models.INADIMPLENTE,
			CreatedAt:    time.Date(2023, time.September, 16, 18, 0, 0, 0, time.FixedZone("", 0)),
		},
		{
			Student:      studentMockData[1],
			Course:       courseMockData[1],
			Installments: 5,
			Status:       models.ADIMPLENTE,
			CreatedAt:    time.Date(2023, time.September, 15, 16, 0, 0, 0, time.FixedZone("", 0)),
		},
		{
			Student:      studentMockData[0],
			Course:       courseMockData[0],
			Installments: 10,
			Status:       models.ADIMPLENTE,
			CreatedAt:    time.Date(2023, time.September, 14, 10, 0, 0, 0, time.FixedZone("", 0)),
		},
	}
)

func TestEnrollmentRepositoryFindPageEnrollments(t *testing.T) {
	pageRequest := types.NewPageRequest(1, 10, []types.Sort{types.NewSort(types.DESC, "e.created_at")})
	datasets := []string{"clear-data.sql", "course-insert.sql", "student-insert.sql", "enrollment-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

	t.Run("Should return enrollments page when exists into enrollments table", func(t *testing.T) {
		expected := enrollmentMockData

		result, err := enrollmentRepository.FindPage(ctx, pageRequest, &models.EnrollmentFilters{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected, result.Content)
		assert.Equal(t, uint64(3), result.TotalElements)
	})

	t.Run("Should return enrollments page filtered by student name when exists into enrollments table", func(t *testing.T) {
		expected := []models.Enrollment{enrollmentMockData[2]}

		result, err := enrollmentRepository.FindPage(ctx, pageRequest, &models.EnrollmentFilters{StudentName: studentMockData[0].Name})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected, result.Content)
		assert.Equal(t, uint64(1), result.TotalElements)
	})

	t.Run("Should return enrollments page filtered by course name when exists into enrollments table", func(t *testing.T) {
		expected := []models.Enrollment{enrollmentMockData[1]}

		result, err := enrollmentRepository.FindPage(ctx, pageRequest, &models.EnrollmentFilters{CourseName: courseMockData[1].Name})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected, result.Content)
		assert.Equal(t, uint64(1), result.TotalElements)
	})
}

func TestEnrollmentRepositoryFindByStudentIdAndCourseId(t *testing.T) {
	datasets := []string{"clear-data.sql", "course-insert.sql", "student-insert.sql", "enrollment-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

	t.Run("Should return nil when studentID not exists into enrollments table", func(t *testing.T) {
		result, err := enrollmentRepository.FindByStudentIdAndCourseId(ctx, uuid.New(), courseMockData[0].ID)

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return nil when studentID not exists into enrollments table", func(t *testing.T) {
		result, err := enrollmentRepository.FindByStudentIdAndCourseId(ctx, studentMockData[0].ID, uuid.New())

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return enrollment when exists into enrollments table", func(t *testing.T) {
		expected := &enrollmentMockData[2]

		result, err := enrollmentRepository.FindByStudentIdAndCourseId(ctx, studentMockData[0].ID, courseMockData[0].ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected, result)
	})
}

func TestEnrollmentRepositoryInsert(t *testing.T) {
	datasets := []string{"clear-data.sql", "course-insert.sql", "student-insert.sql", "enrollment-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

	t.Run("Should return error when enrollment exists into enrollments table", func(t *testing.T) {
		err := enrollmentRepository.Insert(ctx, &models.EnrollmentCreateDTO{
			StudentID:    enrollmentMockData[0].Student.ID,
			CourseID:     enrollmentMockData[0].Course.ID,
			Installments: enrollmentMockData[0].Installments,
		})

		assert.Error(t, err)
	})

	t.Run("Should insert and return enrollment when not exists into enrollments table", func(t *testing.T) {
		model := &models.EnrollmentCreateDTO{
			StudentID:    enrollmentMockData[0].Student.ID,
			CourseID:     enrollmentMockData[1].Course.ID,
			Installments: 10,
		}

		errInsert := enrollmentRepository.Insert(ctx, model)
		result, err := enrollmentRepository.FindByStudentIdAndCourseId(ctx, enrollmentMockData[0].Student.ID, enrollmentMockData[1].Course.ID)

		assert.NoError(t, errInsert)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestEnrollmentRepositoryDelete(t *testing.T) {
	datasets := []string{"clear-data.sql", "course-insert.sql", "student-insert.sql", "enrollment-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

	t.Run("Should delete enrollment when exists into enrollments table", func(t *testing.T) {
		errDelete := enrollmentRepository.Delete(ctx, enrollmentMockData[0].Student.ID, enrollmentMockData[0].Course.ID)
		result, err := enrollmentRepository.FindByStudentIdAndCourseId(ctx, enrollmentMockData[0].Student.ID, enrollmentMockData[0].Course.ID)

		assert.NoError(t, errDelete)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestEnrollmentRepositoryUpdateStatus(t *testing.T) {
	datasets := []string{"clear-data.sql", "course-insert.sql", "student-insert.sql", "enrollment-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

	t.Run("Should update enrollment status when exists into enrollments table", func(t *testing.T) {
		errUpdate := enrollmentRepository.UpdateStatus(ctx, enrollmentMockData[0].Student.ID, enrollmentMockData[0].Course.ID, models.ADIMPLENTE)
		result, err := enrollmentRepository.FindByStudentIdAndCourseId(ctx, enrollmentMockData[0].Student.ID, enrollmentMockData[0].Course.ID)

		assert.NoError(t, errUpdate)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.CreatedAt)
		assert.Equal(t, models.ADIMPLENTE, result.Status)
	})
}
