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

func TestGetStudentByIdUsecase(t *testing.T) {
	t.Run("Should return new get student by id usecase", func(t *testing.T) {
		result := usecases.NewGetStudentByIdUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.StudentRepository)
	})
}

func TestGetStudentByIdUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockStudentsRepository := repositoriesmock.NewMockIStudentsRepository(controller)
	usecase := usecases.GetStudentByIdUsecase{StudentRepository: mockStudentsRepository}
	defer controller.Finish()

	id := uuid.New()

	t.Run("Should return ErrOnFindStudentById when occurred error in FindById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnFindStudentById)
		mockStudentsRepository.EXPECT().FindById(ctx, id).Return(nil, errors.New("mock error in FindById")).MaxTimes(1)

		result, err := usecase.Execute(ctx, id)

		assert.EqualError(t, expected, err.Error())
		assert.Nil(t, result)
	})

	t.Run("Should return ErrStudentNotFound when returns nil in FindById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrStudentNotFound)
		mockStudentsRepository.EXPECT().FindById(ctx, id).Return(nil, nil).MaxTimes(1)

		result, err := usecase.Execute(ctx, id)

		assert.EqualError(t, expected, err.Error())
		assert.Nil(t, result)
	})

	t.Run("Should return student when found", func(t *testing.T) {
		expectedStudent := &models.Student{
			ID:        id,
			Name:      "Student name",
			Email:     "student@example.com",
			Birthday:  types.IsoDate(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
			CreatedAt: time.Now(),
		}
		mockStudentsRepository.EXPECT().FindById(ctx, id).Return(expectedStudent, nil).MaxTimes(1)

		result, err := usecase.Execute(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expectedStudent, result)
	})
}
