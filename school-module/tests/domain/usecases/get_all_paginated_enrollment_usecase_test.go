package usecases

import (
	"errors"
	"testing"
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/usecases"
	repositoriesmock "github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories/mock"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPaginatedEnrollmentUsecase(t *testing.T) {
	t.Run("Should return new get all paginated enrollment usecase", func(t *testing.T) {
		result := usecases.NewGetAllPaginatedEnrollmentUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.EnrollmentRepository)
	})
}

func TestGetAllPaginatedEnrollmentUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockEnrollmentsRepository := repositoriesmock.NewMockIEnrollmentsRepository(controller)
	usecase := usecases.GetAllPaginatedEnrollmentUsecase{EnrollmentRepository: mockEnrollmentsRepository}
	defer controller.Finish()

	params := &models.EnrollmentPageParams{
		Page:        1,
		Size:        10,
		StudentName: "Test Student",
		CourseName:  "Test Course",
	}

	t.Run("Should return ErrOnFindAllPaginatedEnrollment when occurred error in FindAllPaginated", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnFindAllPaginatedEnrollment)
		mockEnrollmentsRepository.EXPECT().FindAllPaginated(ctx, params).Return(nil, errors.New("mock error in FindAllPaginated")).MaxTimes(1)

		result, err := usecase.Execute(ctx, params)

		assert.EqualError(t, expected, err.Error())
		assert.Nil(t, result)
	})

	t.Run("Should return paginated enrollments when successful", func(t *testing.T) {
		enrollments := []models.Enrollment{
			{
				Student: models.Student{
					ID:        uuid.New(),
					Name:      "Student 1",
					Email:     "student1@example.com",
					Birthday:  types.IsoDate(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
					CreatedAt: time.Now(),
				},
				Course: models.Course{
					ID:        uuid.New(),
					Name:      "Course 1",
					Value:     1000,
					CreatedAt: time.Now(),
				},
				Installments: 3,
				Status:       enums.ADIMPLENTE,
				CreatedAt:    time.Now(),
			},
			{
				Student: models.Student{
					ID:        uuid.New(),
					Name:      "Student 2",
					Email:     "student2@example.com",
					Birthday:  types.IsoDate(time.Date(2001, time.February, 1, 0, 0, 0, 0, time.UTC)),
					CreatedAt: time.Now(),
				},
				Course: models.Course{
					ID:        uuid.New(),
					Name:      "Course 2",
					Value:     2000,
					CreatedAt: time.Now(),
				},
				Installments: 5,
				Status:       enums.INADIMPLENTE,
				CreatedAt:    time.Now(),
			},
		}

		expectedPage := &types.Page[models.Enrollment]{
			TotalItems: 2,
			Items:      enrollments,
		}

		mockEnrollmentsRepository.EXPECT().FindAllPaginated(ctx, params).Return(expectedPage, nil).MaxTimes(1)

		result, err := usecase.Execute(ctx, params)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expectedPage, result)
	})
}
