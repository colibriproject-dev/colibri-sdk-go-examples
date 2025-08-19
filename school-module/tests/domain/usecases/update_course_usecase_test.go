package usecases

import (
	"errors"
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/usecases"
	repositoriesmock "github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCourseUsecase(t *testing.T) {
	t.Run("Should return new update course usecase", func(t *testing.T) {
		result := usecases.NewUpdateCourseUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.CourseRepository)
	})
}

func TestUpdateCourseUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockCourseRepository := repositoriesmock.NewMockICoursesRepository(controller)
	usecase := usecases.UpdateCourseUsecase{CourseRepository: mockCourseRepository}
	defer controller.Finish()

	model := &models.CourseUpdate{
		Name:  "Course name",
		Value: 1000,
		ID:    uuid.New(),
	}

	t.Run("Should return ErrOnExistsCourseById when occurred error in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnExistsCourseById)
		mockCourseRepository.EXPECT().ExistsById(ctx, model.ID).Return(nil, errors.New("mock error in ExistsById")).MaxTimes(1)
		mockCourseRepository.EXPECT().FindByName(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockCourseRepository.EXPECT().Update(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrCourseNotFound when returns nil in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrCourseNotFound)
		mockCourseRepository.EXPECT().ExistsById(ctx, model.ID).Return(nil, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().FindByName(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockCourseRepository.EXPECT().Update(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrCourseNotFound when returns false in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrCourseNotFound)
		mockCourseRepository.EXPECT().ExistsById(ctx, model.ID).Return(&notExists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().FindByName(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockCourseRepository.EXPECT().Update(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrOnFindCourseByName when occurred error in FindByName", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnFindCourseByName)
		mockCourseRepository.EXPECT().ExistsById(ctx, model.ID).Return(&exists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().FindByName(ctx, model.Name).Return(nil, errors.New("mock error in FindByName")).MaxTimes(1)
		mockCourseRepository.EXPECT().Update(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrCourseAlreadyExists when occurred error in FindByName", func(t *testing.T) {
		expected := errors.New(exceptions.ErrCourseAlreadyExists)
		mockCourseRepository.EXPECT().ExistsById(ctx, model.ID).Return(&exists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().FindByName(ctx, model.Name).Return(&models.Course{ID: uuid.New(), Name: model.Name}, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().Update(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrOnUpdateCourse when occurred error in Update", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnUpdateCourse)
		mockCourseRepository.EXPECT().ExistsById(ctx, model.ID).Return(&exists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().FindByName(ctx, model.Name).Return(nil, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().Update(ctx, model).Return(errors.New("mock error in Update")).MaxTimes(1)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should update course", func(t *testing.T) {
		mockCourseRepository.EXPECT().ExistsById(ctx, model.ID).Return(&exists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().FindByName(ctx, model.Name).Return(nil, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().Update(ctx, model).Return(nil).MaxTimes(1)

		err := usecase.Execute(ctx, model)

		assert.NoError(t, err)
	})
}
