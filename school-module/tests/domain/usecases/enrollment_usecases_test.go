package usecases

import (
	"errors"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/usecases"
	mockProducers "github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/producers/mock"
	mockRepositories "github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/repositories/mock"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEnrollmentUsecases(t *testing.T) {
	t.Run("Should return new enrollment usecases", func(t *testing.T) {
		result := usecases.NewEnrollmentUsecases()
		assert.NotNil(t, result)
		assert.NotNil(t, result.Repository)
		assert.NotNil(t, result.Producer)
	})
}

func TestGetEnrollmentPage(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockEnrollmentRepository(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in FindPage", func(t *testing.T) {
		expected := errors.New("mock error in FindPage")
		mockRepository.EXPECT().FindPage(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, expected)

		usecase := usecases.EnrollmentUsecases{Repository: mockRepository}
		result, err := usecase.GetPage(ctx, &models.EnrollmentPageParamsDTO{})

		assert.Error(t, expected, err)
		assert.Nil(t, result)
	})

	t.Run("Should return enrollment list", func(t *testing.T) {
		content := []models.Enrollment{
			{
				Student: models.Student{
					ID:        uuid.MustParse("9f9fa978-7df0-4474-b1d4-6be55e0dbd1d"),
					Name:      "Student name 1",
					Email:     "student1@email.com",
					Birthday:  time.Time{},
					CreatedAt: time.Time{},
				},
				Course: models.Course{
					ID:        uuid.MustParse("8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"),
					Name:      "Course name 1",
					Value:     1000,
					CreatedAt: time.Time{},
				},
				Installments: 10,
				Status:       models.INADIMPLENTE,
				CreatedAt:    time.Time{},
			},
		}
		var expected models.EnrollmentPage = &types.Page[models.Enrollment]{TotalElements: 1, Content: content}
		mockRepository.EXPECT().FindPage(gomock.Any(), gomock.Any(), gomock.Any()).Return(expected, nil)

		usecase := usecases.EnrollmentUsecases{Repository: mockRepository}
		result, err := usecase.GetPage(ctx, &models.EnrollmentPageParamsDTO{})

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestCreateEnrollment(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockEnrollmentRepository(controller)
	mockProducer := mockProducers.NewMockIEnrollmentProducer(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in Insert", func(t *testing.T) {
		expected := errors.New("mock error in Insert")
		mockRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expected)

		usecase := usecases.EnrollmentUsecases{Repository: mockRepository}
		err := usecase.Create(ctx, &models.EnrollmentCreateDTO{})

		assert.Error(t, expected, err)
	})

	t.Run("Should create enrollment", func(t *testing.T) {
		mockRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockRepository.EXPECT().FindByStudentIdAndCourseId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.Enrollment{}, nil)
		mockProducer.EXPECT().Create(gomock.Any(), gomock.Any())

		usecase := usecases.EnrollmentUsecases{Repository: mockRepository, Producer: mockProducer}
		err := usecase.Create(ctx, &models.EnrollmentCreateDTO{})

		assert.NoError(t, err)
	})
}

func TestDeleteEnrollment(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockEnrollmentRepository(controller)
	mockProducer := mockProducers.NewMockIEnrollmentProducer(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in Delete", func(t *testing.T) {
		expected := errors.New("mock error in Delete")
		mockRepository.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(expected)

		usecase := usecases.EnrollmentUsecases{Repository: mockRepository}
		err := usecase.Delete(ctx, &models.EnrollmentDeleteParamsDTO{})

		assert.Error(t, expected, err)
	})

	t.Run("Should delete enrollment", func(t *testing.T) {
		mockRepository.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockProducer.EXPECT().Delete(gomock.Any(), gomock.Any())

		usecase := usecases.EnrollmentUsecases{Repository: mockRepository, Producer: mockProducer}
		err := usecase.Delete(ctx, &models.EnrollmentDeleteParamsDTO{})

		assert.NoError(t, err)
	})
}

func TestUpdateEnrollmentStatus(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepository := mockRepositories.NewMockEnrollmentRepository(controller)
	defer controller.Finish()

	t.Run("Should return error when occurred error in UpdateStatus", func(t *testing.T) {
		expected := errors.New("mock error in UpdateStatus")
		mockRepository.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(expected)

		usecase := usecases.EnrollmentUsecases{Repository: mockRepository}
		err := usecase.UpdateStatus(ctx, uuid.New(), uuid.New(), models.ADIMPLENTE)

		assert.Error(t, expected, err)
	})

	t.Run("Should update enrollment status", func(t *testing.T) {
		mockRepository.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		usecase := usecases.EnrollmentUsecases{Repository: mockRepository}
		err := usecase.UpdateStatus(ctx, uuid.New(), uuid.New(), models.ADIMPLENTE)

		assert.NoError(t, err)
	})
}
