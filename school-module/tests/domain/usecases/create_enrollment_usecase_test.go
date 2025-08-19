package usecases

import (
	"errors"
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/usecases"
	producersmock "github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/producers/mock"
	repositoriesmock "github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateEnrollmentUsecase(t *testing.T) {
	t.Run("Should return new create enrollment usecase", func(t *testing.T) {
		result := usecases.NewCreateEnrollmentUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.EnrollmentRepository)
		assert.NotNil(t, result.EnrollmentCreatedProducer)
	})
}

func TestCreateEnrollmentUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockEnrollmentRepository := repositoriesmock.NewMockIEnrollmentsRepository(controller)
	mockCourseRepository := repositoriesmock.NewMockICoursesRepository(controller)
	mockStudentRepository := repositoriesmock.NewMockIStudentsRepository(controller)
	mockEnrollmentCreatedProducer := producersmock.NewMockIEnrollmentCreatedProducer(controller)
	usecase := usecases.CreateEnrollmentUsecase{
		EnrollmentRepository:      mockEnrollmentRepository,
		CourseRepository:          mockCourseRepository,
		StudentRepository:         mockStudentRepository,
		EnrollmentCreatedProducer: mockEnrollmentCreatedProducer,
	}
	defer controller.Finish()

	model := &models.EnrollmentCreate{
		StudentID:    uuid.New(),
		CourseID:     uuid.New(),
		Installments: 1,
		Status:       enums.ADIMPLENTE,
	}

	enrollmentCreated := &models.EnrollmentCreated{
		Student:      models.EnrollmentCreatedStudent{ID: model.StudentID},
		Course:       models.EnrollmentCreatedCourse{ID: model.CourseID},
		Installments: model.Installments,
		Status:       model.Status,
	}

	t.Run("Should return ErrOnExistsEnrollmentByStudentIdAndCourseId when occurred error in ExistsByStudentIdAndCourseId", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnExistsEnrollmentByStudentIdAndCourseId)
		mockEnrollmentRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID).Return(nil, errors.New("mock error in ExistsByStudentIdAndCourseId")).MaxTimes(1)
		mockCourseRepository.EXPECT().ExistsById(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockStudentRepository.EXPECT().ExistsById(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockEnrollmentRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockEnrollmentCreatedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrEnrollmentAlreadyExists when returns true in ExistsByStudentIdAndCourseId", func(t *testing.T) {
		expected := errors.New(exceptions.ErrEnrollmentAlreadyExists)
		mockEnrollmentRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID).Return(&exists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().ExistsById(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockStudentRepository.EXPECT().ExistsById(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockEnrollmentRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockEnrollmentCreatedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrOnExistsCourseById when occurred error in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnExistsCourseById)
		mockEnrollmentRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID).Return(&notExists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().ExistsById(ctx, model.CourseID).Return(nil, errors.New("mock error in ExistsById")).MaxTimes(1)
		mockStudentRepository.EXPECT().ExistsById(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockEnrollmentRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockEnrollmentCreatedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrCourseNotFound when returns false in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrCourseNotFound)
		mockEnrollmentRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID).Return(&notExists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().ExistsById(ctx, model.CourseID).Return(&notExists, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().ExistsById(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockEnrollmentRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockEnrollmentCreatedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrOnExistsStudentById when occurred error in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnExistsStudentById)
		mockEnrollmentRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID).Return(&notExists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().ExistsById(ctx, model.CourseID).Return(&exists, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().ExistsById(ctx, model.StudentID).Return(nil, errors.New("mock error in ExistsById")).MaxTimes(1)
		mockEnrollmentRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockEnrollmentCreatedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrStudentNotFound when returns false in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrStudentNotFound)
		mockEnrollmentRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID).Return(&notExists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().ExistsById(ctx, model.CourseID).Return(&exists, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().ExistsById(ctx, model.StudentID).Return(&notExists, nil).MaxTimes(1)
		mockEnrollmentRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).MaxTimes(0)
		mockEnrollmentCreatedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should return ErrOnInsertEnrollment when occurred error in Insert", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnInsertEnrollment)
		mockEnrollmentRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID).Return(&notExists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().ExistsById(ctx, model.CourseID).Return(&exists, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().ExistsById(ctx, model.StudentID).Return(&exists, nil).MaxTimes(1)
		mockEnrollmentRepository.EXPECT().Insert(ctx, model).Return(nil, errors.New("mock error in Insert")).MaxTimes(1)
		mockEnrollmentCreatedProducer.EXPECT().Send(gomock.Any(), gomock.Any()).MaxTimes(0)

		err := usecase.Execute(ctx, model)

		assert.EqualError(t, expected, err.Error())
	})

	t.Run("Should create enrollment and send enrollment created notification", func(t *testing.T) {
		mockEnrollmentRepository.EXPECT().ExistsByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID).Return(&notExists, nil).MaxTimes(1)
		mockCourseRepository.EXPECT().ExistsById(ctx, model.CourseID).Return(&exists, nil).MaxTimes(1)
		mockStudentRepository.EXPECT().ExistsById(ctx, model.StudentID).Return(&exists, nil).MaxTimes(1)
		mockEnrollmentRepository.EXPECT().Insert(ctx, model).Return(enrollmentCreated, nil).MaxTimes(1)
		mockEnrollmentCreatedProducer.EXPECT().Send(ctx, enrollmentCreated).Return(nil).MaxTimes(1)

		err := usecase.Execute(ctx, model)

		assert.NoError(t, err)
	})
}
