package usecases

import (
	"errors"
	"testing"
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/usecases"
	repositoriesmock "github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetCourseByIdUsecase(t *testing.T) {
	t.Run("Should return new get course by id usecase", func(t *testing.T) {
		result := usecases.NewGetCourseByIdUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.CourseRepository)
	})
}

func TestGetCourseByIdUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockCoursesRepository := repositoriesmock.NewMockICoursesRepository(controller)
	usecase := usecases.GetCourseByIdUsecase{CourseRepository: mockCoursesRepository}
	defer controller.Finish()

	id := uuid.New()

	t.Run("Should return ErrOnFindCourseById when occurred error in FindById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnFindCourseById)
		mockCoursesRepository.EXPECT().FindById(ctx, id).Return(nil, errors.New("mock error in FindById"))

		result, err := usecase.Execute(ctx, id)

		assert.EqualError(t, expected, err.Error())
		assert.Nil(t, result)
	})

	t.Run("Should return ErrCourseNotFound when return nil in FindById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrCourseNotFound)
		mockCoursesRepository.EXPECT().FindById(ctx, id).Return(nil, nil)

		result, err := usecase.Execute(ctx, id)

		assert.EqualError(t, expected, err.Error())
		assert.Nil(t, result)
	})

	t.Run("Should return course", func(t *testing.T) {
		expected := &models.Course{
			ID:        id,
			Name:      "Course name 1",
			Value:     1000,
			CreatedAt: time.Now(),
		}
		mockCoursesRepository.EXPECT().FindById(ctx, id).Return(expected, nil)

		result, err := usecase.Execute(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
	})
}
