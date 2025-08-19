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

func TestDeleteCourseUsecase(t *testing.T) {
	t.Run("Should return new delete course usecase", func(t *testing.T) {
		result := usecases.NewDeleteCourseUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.CourseRepository)
		assert.NotNil(t, result.CourseDeletedProducer)
	})
}

func TestDeleteCourseUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockCourseRepository := repositoriesmock.NewMockICoursesRepository(controller)
	mockCourseDeletedProducer := producersmock.NewMockICourseDeletedProducer(controller)
	usecase := usecases.DeleteCourseUsecase{CourseRepository: mockCourseRepository, CourseDeletedProducer: mockCourseDeletedProducer}
	defer controller.Finish()

	id := uuid.New()

	t.Run("Should return ErrOnExistsCourseById when occurred error in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnExistsCourseById)
		mockCourseRepository.EXPECT().ExistsById(ctx, id).Return(nil, errors.New("mock error in ExistsById")).MaxTimes(1)
		mockCourseRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockCourseDeletedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, id)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrCourseNotFound when return nil in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrCourseNotFound)
		mockCourseRepository.EXPECT().ExistsById(ctx, id).Return(nil, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockCourseDeletedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, id)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrCourseNotFound when return false in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrCourseNotFound)
		mockCourseRepository.EXPECT().ExistsById(ctx, id).Return(&notExists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockCourseDeletedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, id)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrOnDeleteCourse when occurred error in Delete", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnDeleteCourse)
		mockCourseRepository.EXPECT().ExistsById(ctx, id).Return(&exists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().Delete(ctx, id).Return(errors.New("mock error in Delete")).MaxTimes(1)
		mockCourseDeletedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, id)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should delete course and send deleted course notification", func(t *testing.T) {
		mockCourseRepository.EXPECT().ExistsById(ctx, id).Return(&exists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().Delete(ctx, id).Return(nil).MaxTimes(1)
		mockCourseDeletedProducer.EXPECT().Send(ctx, &models.CourseDelete{ID: id}).Return(nil).MaxTimes(1)

		err := usecase.Execute(ctx, id)

		assert.NoError(t, err)
	})
}
