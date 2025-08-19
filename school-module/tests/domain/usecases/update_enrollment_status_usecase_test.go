package usecases

import (
	"errors"
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/usecases"
	repositoriesmock "github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateEnrollmentStatusUsecase(t *testing.T) {
	t.Run("Should return new update enrollment status usecase", func(t *testing.T) {
		result := usecases.NewUpdateEnrollmentStatusUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.EnrollmentRepository)
	})
}

func TestUpdateEnrollmentStatusUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockEnrollmentsRepository := repositoriesmock.NewMockIEnrollmentsRepository(controller)
	usecase := usecases.UpdateEnrollmentStatusUsecase{EnrollmentRepository: mockEnrollmentsRepository}
	defer controller.Finish()

	model := &models.EnrollmentUpdateStatus{
		StudentID: uuid.New(),
		CourseID:  uuid.New(),
		Status:    enums.INADIMPLENTE,
	}

	t.Run("Should return ErrOnExistsEnrollmentByStudentIdAndCourseId when occurred error in ExistsByStudentIdAndCourseId", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnExistsEnrollmentByStudentIdAndCourseId)
		mockEnrollmentsRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID).Return(nil, errors.New("mock error in ExistsByStudentIdAndCourseId")).MaxTimes(1)
		mockEnrollmentsRepository.EXPECT().UpdateStatus(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrEnrollmentNotFound when returns false in ExistsByStudentIdAndCourseId", func(t *testing.T) {
		expected := errors.New(exceptions.ErrEnrollmentNotFound)
		mockEnrollmentsRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID).Return(&notExists, nil).MaxTimes(1)
		mockEnrollmentsRepository.EXPECT().UpdateStatus(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrOnUpdateEnrollmentStatus when occurred error in UpdateStatus", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnUpdateEnrollmentStatus)
		mockEnrollmentsRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID).Return(&exists, nil).MaxTimes(1)
		mockEnrollmentsRepository.EXPECT().UpdateStatus(ctx, model).Return(errors.New("mock error in UpdateStatus"))

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should update enrollment status successfully", func(t *testing.T) {
		mockEnrollmentsRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID).Return(&exists, nil).MaxTimes(1)
		mockEnrollmentsRepository.EXPECT().UpdateStatus(ctx, model).Return(nil)

		err := usecase.Execute(ctx, model)

		assert.NoError(t, err)
	})
}
