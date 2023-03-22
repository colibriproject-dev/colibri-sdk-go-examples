package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/usecases"
	mockProducers "github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/producers/mock"
	mockRepositories "github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/repositories/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.Background()
)

func TestCourseUsecases(t *testing.T) {
	t.Run("Should return new course usecases", func(t *testing.T) {
		result := usecases.NewCourseUsecases()
		assert.NotNil(t, result)
		assert.NotNil(t, result.Repository)
		assert.NotNil(t, result.Producer)
	})
}

func TestGetAllCourse(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockCourseRepository(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in FindAll", func(t *testing.T) {
		expected := errors.New("mock error in FindAll")
		mockRepository.EXPECT().FindAll(gomock.Any()).Return(nil, expected)

		usecase := usecases.CourseUsecases{Repository: mockRepository}
		result, err := usecase.GetAll(ctx)

		assert.Error(t, expected, err)
		assert.Nil(t, result)
	})

	t.Run("Should return course list", func(t *testing.T) {
		expected := []models.Course{
			{ID: uuid.MustParse("8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"), Name: "Course name 1", Value: 1000, CreatedAt: time.Time{}},
			{ID: uuid.MustParse("5979def8-f9c2-4258-bbf1-118c1ffb0a56"), Name: "Course name 2", Value: 2000, CreatedAt: time.Time{}},
		}
		mockRepository.EXPECT().FindAll(gomock.Any()).Return(expected, nil)

		usecase := usecases.CourseUsecases{Repository: mockRepository}
		result, err := usecase.GetAll(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestCourseGetById(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockCourseRepository(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in FindById", func(t *testing.T) {
		expected := errors.New("mock error in FindById")
		mockRepository.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(nil, expected)

		usecase := usecases.CourseUsecases{Repository: mockRepository}
		result, err := usecase.GetById(ctx, uuid.New())

		assert.Error(t, expected, err)
		assert.Nil(t, result)
	})

	t.Run("Should return course", func(t *testing.T) {
		expected := &models.Course{
			ID:        uuid.MustParse("8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"),
			Name:      "Course name 1",
			Value:     1000,
			CreatedAt: time.Time{},
		}
		mockRepository.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(expected, nil)

		usecase := usecases.CourseUsecases{Repository: mockRepository}
		result, err := usecase.GetById(ctx, expected.ID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestCreateCourse(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockCourseRepository(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in Insert", func(t *testing.T) {
		expected := errors.New("mock error in Insert")
		mockRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil, expected)

		usecase := usecases.CourseUsecases{Repository: mockRepository}
		result, err := usecase.Create(ctx, &models.CourseCreateUpdateDTO{})

		assert.Error(t, expected, err)
		assert.Nil(t, result)
	})

	t.Run("Should create course", func(t *testing.T) {
		mockRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(&models.Course{}, nil)

		usecase := usecases.CourseUsecases{Repository: mockRepository}
		result, err := usecase.Create(ctx, &models.CourseCreateUpdateDTO{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestUpdateCourse(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockCourseRepository(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in Update", func(t *testing.T) {
		expected := errors.New("mock error in Update")
		mockRepository.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(expected)

		usecase := usecases.CourseUsecases{Repository: mockRepository}
		err := usecase.Update(ctx, uuid.New(), &models.CourseCreateUpdateDTO{})

		assert.Error(t, expected, err)
	})

	t.Run("Should update course", func(t *testing.T) {
		mockRepository.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		usecase := usecases.CourseUsecases{Repository: mockRepository}
		err := usecase.Update(ctx, uuid.New(), &models.CourseCreateUpdateDTO{})

		assert.NoError(t, err)
	})
}

func TestDeleteCourse(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockCourseRepository(controller)
	mockProducer := mockProducers.NewMockICourseProducer(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in Delete", func(t *testing.T) {
		expected := errors.New("mock error in Delete")
		mockRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(expected)

		usecase := usecases.CourseUsecases{Repository: mockRepository}
		err := usecase.Delete(ctx, uuid.New())

		assert.Error(t, expected, err)
	})

	t.Run("Should delete course", func(t *testing.T) {
		mockRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
		mockProducer.EXPECT().Delete(gomock.Any(), gomock.Any())

		usecase := usecases.CourseUsecases{Repository: mockRepository, Producer: mockProducer}
		err := usecase.Delete(ctx, uuid.New())

		assert.NoError(t, err)
	})
}
