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

func TestGetAllCourseUsecase(t *testing.T) {
	t.Run("Should return new get all course usecase", func(t *testing.T) {
		result := usecases.NewGetAllCourseUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.CourseRepository)
	})
}

func TestGetAllCourseUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockCoursesRepository := repositoriesmock.NewMockICoursesRepository(controller)
	usecase := usecases.GetAllCourseUsecase{CourseRepository: mockCoursesRepository}
	defer controller.Finish()

	t.Run("Should return ErrOnFindAllCourses when occurred error in FindAll", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnFindAllCourses)
		mockCoursesRepository.EXPECT().FindAll(ctx).Return(nil, errors.New("mock error in FindAll"))

		result, err := usecase.Execute(ctx)

		assert.EqualError(t, expected, err.Error())
		assert.Nil(t, result)
	})

	t.Run("Should return course list", func(t *testing.T) {
		expected := []models.Course{
			{ID: uuid.New(), Name: "Course name 1", Value: 1000, CreatedAt: time.Now()},
			{ID: uuid.New(), Name: "Course name 2", Value: 2000, CreatedAt: time.Now()},
		}
		mockCoursesRepository.EXPECT().FindAll(ctx).Return(expected, nil)

		result, err := usecase.Execute(ctx)

		assert.NoError(t, err)
		assert.EqualValues(t, expected, result)
	})
}
