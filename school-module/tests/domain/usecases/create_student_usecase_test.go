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
	"github.com/stretchr/testify/assert"
)

func TestCreateStudentUsecase(t *testing.T) {
	t.Run("Should return new create student usecase", func(t *testing.T) {
		result := usecases.NewCreateStudentUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.Repository)
	})
}

func TestCreateStudentUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockStudentsRepository := repositoriesmock.NewMockIStudentsRepository(controller)
	usecase := usecases.CreateStudentUsecase{Repository: mockStudentsRepository}
	defer controller.Finish()

	model := &models.StudentCreate{
		Name:     "Student name",
		Email:    "student@example.com",
		Birthday: types.IsoDate(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
	}

	t.Run("Should return ErrOnExistsStudentByEmail when occurred error in ExistsByEmail", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnExistsStudentByEmail)
		mockStudentsRepository.EXPECT().ExistsByEmail(ctx, model.Email).Return(nil, errors.New("mock error in ExistsByEmail")).MaxTimes(1)
		mockStudentsRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrStudentAlreadyExists when returns true in ExistsByEmail", func(t *testing.T) {
		expected := errors.New(exceptions.ErrStudentAlreadyExists)
		mockStudentsRepository.EXPECT().ExistsByEmail(ctx, model.Email).Return(&exists, nil).MaxTimes(1)
		mockStudentsRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrOnInsertStudent when occurred error in Insert", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnInsertStudent)
		mockStudentsRepository.EXPECT().ExistsByEmail(ctx, model.Email).Return(&notExists, nil).MaxTimes(1)
		mockStudentsRepository.EXPECT().Insert(ctx, model).Return(errors.New("mock error in Insert"))

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should create student successfully", func(t *testing.T) {
		mockStudentsRepository.EXPECT().ExistsByEmail(ctx, model.Email).Return(&notExists, nil).MaxTimes(1)
		mockStudentsRepository.EXPECT().Insert(ctx, model).Return(nil)

		err := usecase.Execute(ctx, model)

		assert.NoError(t, err)
	})
}
