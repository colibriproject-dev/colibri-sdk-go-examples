package usecases

import (
	"errors"
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/usecases"
	producersmock "github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/producers/mock"
	repositoriesmock "github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteStudentUsecase(t *testing.T) {
	t.Run("Should return new delete student usecase", func(t *testing.T) {
		result := usecases.NewDeleteStudentUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.Repository)
		assert.NotNil(t, result.StudentDeletedProducer)
	})
}

func TestDeleteStudentUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockStudentRepository := repositoriesmock.NewMockIStudentsRepository(controller)
	mockStudentDeletedProducer := producersmock.NewMockIStudentDeletedProducer(controller)
	usecase := usecases.DeleteStudentUsecase{
		Repository:             mockStudentRepository,
		StudentDeletedProducer: mockStudentDeletedProducer,
	}
	defer controller.Finish()

	id := uuid.New()

	t.Run("Should return ErrOnExistsStudentById when occurred error in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnExistsStudentById)
		mockStudentRepository.EXPECT().ExistsById(ctx, id).Return(nil, errors.New("mock error in ExistsById")).MaxTimes(1)
		mockStudentRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockStudentDeletedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, id)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrStudentNotFound when returns false in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrStudentNotFound)
		mockStudentRepository.EXPECT().ExistsById(ctx, id).Return(&notExists, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockStudentDeletedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, id)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrOnDeleteStudent when occurred error in Delete", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnDeleteStudent)
		mockStudentRepository.EXPECT().ExistsById(ctx, id).Return(&exists, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().Delete(ctx, id).Return(errors.New("mock error in Delete"))
		mockStudentDeletedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, id)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should delete student and send notification successfully", func(t *testing.T) {
		mockStudentRepository.EXPECT().ExistsById(ctx, id).Return(&exists, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().Delete(ctx, id).Return(nil)
		mockStudentDeletedProducer.EXPECT().Send(ctx, &models.StudentDelete{ID: id}).Return(nil)

		err := usecase.Execute(ctx, id)

		assert.NoError(t, err)
	})
}
