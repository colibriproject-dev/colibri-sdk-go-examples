package usecases

import (
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

func TestStudentUsecases(t *testing.T) {
	t.Run("Should return new student usecases", func(t *testing.T) {
		result := usecases.NewStudentUsecases()
		assert.NotNil(t, result)
		assert.NotNil(t, result.Repository)
		assert.NotNil(t, result.Producer)
	})
}

func TestGetAllStudent(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockStudentRepository(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in FindAll", func(t *testing.T) {
		expected := errors.New("mock error in FindAll")
		mockRepository.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(nil, expected)

		usecase := usecases.StudentUsecases{Repository: mockRepository}
		result, err := usecase.GetAll(ctx, &models.StudentParams{})

		assert.Error(t, expected, err)
		assert.Nil(t, result)
	})

	t.Run("Should return student list", func(t *testing.T) {
		expected := []models.Student{
			{ID: uuid.MustParse("9f9fa978-7df0-4474-b1d4-6be55e0dbd1d"), Name: "Student name 1", Email: "student1@email.com", Birthday: time.Time{}, CreatedAt: time.Time{}},
			{ID: uuid.MustParse("6079def8-f9c2-4258-bbf1-118c1ffb0a67"), Name: "Student name 2", Email: "student2@email.com", Birthday: time.Time{}, CreatedAt: time.Time{}},
		}
		mockRepository.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(expected, nil)

		usecase := usecases.StudentUsecases{Repository: mockRepository}
		result, err := usecase.GetAll(ctx, &models.StudentParams{})

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestStudentGetById(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockStudentRepository(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in FindById", func(t *testing.T) {
		expected := errors.New("mock error in FindById")
		mockRepository.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(nil, expected)

		usecase := usecases.StudentUsecases{Repository: mockRepository}
		result, err := usecase.GetById(ctx, uuid.New())

		assert.Error(t, expected, err)
		assert.Nil(t, result)
	})

	t.Run("Should return student", func(t *testing.T) {
		expected := &models.Student{
			ID:        uuid.MustParse("9f9fa978-7df0-4474-b1d4-6be55e0dbd1d"),
			Name:      "Student name 1",
			Email:     "student1@email.com",
			Birthday:  time.Time{},
			CreatedAt: time.Time{},
		}
		mockRepository.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(expected, nil)

		usecase := usecases.StudentUsecases{Repository: mockRepository}
		result, err := usecase.GetById(ctx, expected.ID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestCreateStudent(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockStudentRepository(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in Insert", func(t *testing.T) {
		expected := errors.New("mock error in Insert")
		mockRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expected)

		usecase := usecases.StudentUsecases{Repository: mockRepository}
		err := usecase.Create(ctx, &models.StudentCreateUpdateDTO{})

		assert.Error(t, expected, err)
	})

	t.Run("Should create student", func(t *testing.T) {
		mockRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		usecase := usecases.StudentUsecases{Repository: mockRepository}
		err := usecase.Create(ctx, &models.StudentCreateUpdateDTO{})

		assert.NoError(t, err)
	})
}

func TestUpdateStudent(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockStudentRepository(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in Update", func(t *testing.T) {
		expected := errors.New("mock error in Update")
		mockRepository.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(expected)

		usecase := usecases.StudentUsecases{Repository: mockRepository}
		err := usecase.Update(ctx, uuid.New(), &models.StudentCreateUpdateDTO{})

		assert.Error(t, expected, err)
	})

	t.Run("Should update student", func(t *testing.T) {
		mockRepository.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		usecase := usecases.StudentUsecases{Repository: mockRepository}
		err := usecase.Update(ctx, uuid.New(), &models.StudentCreateUpdateDTO{})

		assert.NoError(t, err)
	})
}

func TestDeleteStudent(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockStudentRepository(controller)
	mockProducer := mockProducers.NewMockIStudentDeletedProducer(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in Delete", func(t *testing.T) {
		expected := errors.New("mock error in Delete")
		mockRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(expected)

		usecase := usecases.StudentUsecases{Repository: mockRepository}
		err := usecase.Delete(ctx, uuid.New())

		assert.Error(t, expected, err)
	})

	t.Run("Should delete student", func(t *testing.T) {
		mockRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
		mockProducer.EXPECT().Delete(gomock.Any(), gomock.Any())

		usecase := usecases.StudentUsecases{Repository: mockRepository, Producer: mockProducer}
		err := usecase.Delete(ctx, uuid.New())

		assert.NoError(t, err)
	})
}

func TestUploadStudentDocument(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockStudentRepository(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in UploadDocument", func(t *testing.T) {
		expected := errors.New("mock error in UploadDocument")
		mockRepository.EXPECT().UploadDocument(gomock.Any(), gomock.Any(), gomock.Any()).Return("", expected)

		usecase := usecases.StudentUsecases{Repository: mockRepository}
		result, err := usecase.UploadDocument(ctx, uuid.New(), nil)

		assert.Error(t, expected, err)
		assert.Empty(t, result)
	})

	t.Run("Should upload student document and return url", func(t *testing.T) {
		expected := "http://my-file.url"
		mockRepository.EXPECT().UploadDocument(gomock.Any(), gomock.Any(), gomock.Any()).Return(expected, nil)

		usecase := usecases.StudentUsecases{Repository: mockRepository}
		result, err := usecase.UploadDocument(ctx, uuid.New(), nil)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
