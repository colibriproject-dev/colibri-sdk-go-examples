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

func TestDeleteEnrollmentUsecase(t *testing.T) {
	t.Run("Should return new delete enrollment usecase", func(t *testing.T) {
		result := usecases.NewDeleteEnrollmentUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.EnrollmentRepository)
		assert.NotNil(t, result.EnrollmentDeletedProducer)
	})
}

func TestDeleteEnrollmentUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockEnrollmentRepository := repositoriesmock.NewMockIEnrollmentsRepository(controller)
	mockEnrollmentDeletedProducer := producersmock.NewMockIEnrollmentDeletedProducer(controller)
	usecase := usecases.DeleteEnrollmentUsecase{
		EnrollmentRepository:      mockEnrollmentRepository,
		EnrollmentDeletedProducer: mockEnrollmentDeletedProducer,
	}
	defer controller.Finish()

	params := &models.EnrollmentDelete{
		StudentID: uuid.New(),
		CourseID:  uuid.New(),
	}

	t.Run("Should return ErrOnExistsEnrollmentByStudentIdAndCourseId when occurred error in ExistsByStudentIdAndCourseId", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnExistsEnrollmentByStudentIdAndCourseId)
		mockEnrollmentRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, params.StudentID, params.CourseID).Return(nil, errors.New("mock error in ExistsByStudentIdAndCourseId")).MaxTimes(1)
		mockEnrollmentRepository.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(0)
		mockEnrollmentDeletedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, params)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrEnrollmentNotFound when returns false in ExistsByStudentIdAndCourseId", func(t *testing.T) {
		expected := errors.New(exceptions.ErrEnrollmentNotFound)
		mockEnrollmentRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, params.StudentID, params.CourseID).Return(&notExists, nil).MaxTimes(1)
		mockEnrollmentRepository.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(0)
		mockEnrollmentDeletedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, params)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrOnDeleteEnrollment when occurred error in Delete", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnDeleteEnrollment)
		mockEnrollmentRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, params.StudentID, params.CourseID).Return(&exists, nil).MaxTimes(1)
		mockEnrollmentRepository.EXPECT().Delete(ctx, params.StudentID, params.CourseID).Return(errors.New("mock error in Delete"))
		mockEnrollmentDeletedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, params)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should delete enrollment and send notification successfully", func(t *testing.T) {
		mockEnrollmentRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, params.StudentID, params.CourseID).Return(&exists, nil).MaxTimes(1)
		mockEnrollmentRepository.EXPECT().Delete(ctx, params.StudentID, params.CourseID).Return(nil)
		mockEnrollmentDeletedProducer.EXPECT().Send(ctx, params).Return(nil)

		err := usecase.Execute(ctx, params)

		assert.NoError(t, err)
	})
}
