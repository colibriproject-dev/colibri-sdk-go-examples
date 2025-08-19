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

func TestUpdateStudentUsecase(t *testing.T) {
	t.Run("Should return new update student usecase", func(t *testing.T) {
		result := usecases.NewUpdateStudentUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.StudentRepository)
	})
}

func TestUpdateStudentUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockStudentRepository := repositoriesmock.NewMockIStudentsRepository(controller)
	usecase := usecases.UpdateStudentUsecase{StudentRepository: mockStudentRepository}
	defer controller.Finish()

	model := &models.StudentUpdate{
		ID:       uuid.New(),
		Name:     "Updated Student name",
		Email:    "updated@example.com",
		Birthday: types.IsoDate(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
	}

	existingStudent := &models.Student{
		ID:    uuid.New(),
		Name:  model.Name,
		Email: model.Email,
	}

	t.Run("Should return ErrOnExistsStudentById when occurred error in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnExistsStudentById)
		mockStudentRepository.EXPECT().ExistsById(ctx, model.ID).Return(nil, errors.New("mock error in ExistsById")).MaxTimes(1)
		mockStudentRepository.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockStudentRepository.EXPECT().Update(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrStudentNotFound when returns false in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrStudentNotFound)
		mockStudentRepository.EXPECT().ExistsById(ctx, model.ID).Return(&notExists, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockStudentRepository.EXPECT().Update(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrOnFindStudentByEmail when occurred error in FindByEmail", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnFindStudentByEmail)
		mockStudentRepository.EXPECT().ExistsById(ctx, model.ID).Return(&exists, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().FindByEmail(ctx, model.Email).Return(nil, errors.New("mock error in FindByEmail")).MaxTimes(1)
		mockStudentRepository.EXPECT().Update(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrStudentAlreadyExists when email belongs to another student", func(t *testing.T) {
		expected := errors.New(exceptions.ErrStudentAlreadyExists)
		mockStudentRepository.EXPECT().ExistsById(ctx, model.ID).Return(&exists, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().FindByEmail(ctx, model.Email).Return(existingStudent, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().Update(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrOnUpdateStudent when occurred error in Update", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnUpdateStudent)
		mockStudentRepository.EXPECT().ExistsById(ctx, model.ID).Return(&exists, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().FindByEmail(ctx, model.Email).Return(nil, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().Update(ctx, model).Return(errors.New("mock error in Update")).MaxTimes(1)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should update student successfully", func(t *testing.T) {
		mockStudentRepository.EXPECT().ExistsById(ctx, model.ID).Return(&exists, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().FindByEmail(ctx, model.Email).Return(nil, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().Update(ctx, model).Return(nil).MaxTimes(1)

		err := usecase.Execute(ctx, model)

		assert.NoError(t, err)
	})
}
