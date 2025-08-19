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

func TestCreateCourseUsecase(t *testing.T) {
	t.Run("Should return new create course usecase", func(t *testing.T) {
		result := usecases.NewCreateCourseUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.CourseRepository)
	})
}

func TestCreateCourseUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockCoursesRepository := repositoriesmock.NewMockICoursesRepository(controller)
	usecase := usecases.CreateCourseUsecase{CourseRepository: mockCoursesRepository}
	defer controller.Finish()

	model := &models.CourseCreate{
		Name:  "Course name",
		Value: 1000,
	}

	t.Run("Should return ErrOnExistsCourseByName when occurred error in ExistsByName", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnExistsCourseByName)
		mockCoursesRepository.EXPECT().ExistsByName(ctx, model.Name).Return(nil, errors.New("mock error in ExistsByName")).MaxTimes(1)
		mockCoursesRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).MaxTimes(0)

		result, err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
		assert.Nil(t, result)
	})

	t.Run("Should return ErrCourseAlreadyExists when returns true in ExistsByName", func(t *testing.T) {
		expected := errors.New(exceptions.ErrCourseAlreadyExists)
		mockCoursesRepository.EXPECT().ExistsByName(ctx, model.Name).Return(&exists, nil).MaxTimes(1)
		mockCoursesRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).MaxTimes(0)

		result, err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
		assert.Nil(t, result)
	})

	t.Run("Should return ErrOnInsertCourse when occurred error in Insert", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnInsertCourse)
		mockCoursesRepository.EXPECT().ExistsByName(ctx, model.Name).Return(&notExists, nil).MaxTimes(1)
		mockCoursesRepository.EXPECT().Insert(ctx, model).Return(nil, errors.New("mock error in Insert"))

		result, err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
		assert.Nil(t, result)
	})

	t.Run("Should create course and return it", func(t *testing.T) {
		expected := &models.Course{
			ID:        uuid.New(),
			Name:      model.Name,
			Value:     model.Value,
			CreatedAt: time.Now(),
		}
		mockCoursesRepository.EXPECT().ExistsByName(ctx, model.Name).Return(&notExists, nil).MaxTimes(1)
		mockCoursesRepository.EXPECT().Insert(ctx, model).Return(expected, nil)

		result, err := usecase.Execute(ctx, model)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
	})
}
