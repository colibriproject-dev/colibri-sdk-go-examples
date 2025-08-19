package usecases

import (
	"errors"
	"testing"
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/usecases"
	repositoriesmock "github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories/mock"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPaginatedStudentUsecase(t *testing.T) {
	t.Run("Should return new get all paginated student usecase", func(t *testing.T) {
		result := usecases.NewGetAllPaginatedStudentUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.StudentRepository)
	})
}

func TestGetAllPaginatedStudentUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockStudentsRepository := repositoriesmock.NewMockIStudentsRepository(controller)
	usecase := usecases.GetAllPaginatedStudentUsecase{StudentRepository: mockStudentsRepository}
	defer controller.Finish()

	params := &models.StudentPageParams{
		Page: 1,
		Size: 10,
		Name: "Test",
	}

	t.Run("Should return ErrOnFindAllPaginatedStudents when occurred error in FindAllPaginated", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnFindAllPaginatedStudents)
		mockStudentsRepository.EXPECT().FindAllPaginated(ctx, params).Return(nil, errors.New("mock error in FindAllPaginated")).MaxTimes(1)

		result, err := usecase.Execute(ctx, params)

		assert.EqualError(t, expected, err.Error())
		assert.Nil(t, result)
	})

	t.Run("Should return paginated students when successful", func(t *testing.T) {
		students := []models.Student{
			{
				ID:        uuid.New(),
				Name:      "Student 1",
				Email:     "student1@example.com",
				Birthday:  types.IsoDate(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				Name:      "Student 2",
				Email:     "student2@example.com",
				Birthday:  types.IsoDate(time.Date(2001, time.February, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt: time.Now(),
			},
		}

		expectedPage := &types.Page[models.Student]{
			TotalItems: 2,
			Items:      students,
		}

		mockStudentsRepository.EXPECT().FindAllPaginated(ctx, params).Return(expectedPage, nil).MaxTimes(1)

		result, err := usecase.Execute(ctx, params)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expectedPage, result)
	})
}
