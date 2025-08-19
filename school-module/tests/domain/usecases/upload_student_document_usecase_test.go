package usecases

import (
	"errors"
	"mime/multipart"
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/usecases"
	repositoriesmock "github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories/mock"
	storagesmock "github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/storages/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUploadStudentDocumentUsecase(t *testing.T) {
	t.Run("Should return new upload student document usecase", func(t *testing.T) {
		result := usecases.NewUploadStudentDocumentUsecase()
		assert.NotNil(t, result)
		assert.NotNil(t, result.StudentRepository)
		assert.NotNil(t, result.DocumentStorage)
	})
}

func TestUploadStudentDocumentUsecase_Execute(t *testing.T) {
	controller := gomock.NewController(t)
	mockStudentRepository := repositoriesmock.NewMockIStudentsRepository(controller)
	mockDocumentStorage := storagesmock.NewMockIDocumentStorage(controller)
	usecase := usecases.UploadStudentDocumentUsecase{
		StudentRepository: mockStudentRepository,
		DocumentStorage:   mockDocumentStorage,
	}
	defer controller.Finish()

	id := uuid.New()
	var file multipart.File

	t.Run("Should return ErrOnExistsStudentById when occurred error in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnExistsStudentById)
		mockStudentRepository.EXPECT().ExistsById(ctx, id).Return(nil, errors.New("mock error in ExistsById")).MaxTimes(1)
		mockDocumentStorage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(0)

		result, err := usecase.Execute(ctx, id, &file)

		assert.EqualError(t, expected, err.Error())
		assert.Nil(t, result)
	})

	t.Run("Should return ErrStudentNotFound when returns false in ExistsById", func(t *testing.T) {
		expected := errors.New(exceptions.ErrStudentNotFound)
		mockStudentRepository.EXPECT().ExistsById(ctx, id).Return(&notExists, nil).MaxTimes(1)
		mockDocumentStorage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(0)

		result, err := usecase.Execute(ctx, id, &file)

		assert.EqualError(t, expected, err.Error())
		assert.Nil(t, result)
	})

	t.Run("Should return ErrOnUploadStudentDocument when occurred error in Upload", func(t *testing.T) {
		expected := errors.New(exceptions.ErrOnUploadStudentDocument)
		mockStudentRepository.EXPECT().ExistsById(ctx, id).Return(&exists, nil).MaxTimes(1)
		mockDocumentStorage.EXPECT().Upload(ctx, id, &file).Return("", errors.New("mock error in Upload"))

		result, err := usecase.Execute(ctx, id, &file)

		assert.EqualError(t, expected, err.Error())
		assert.Nil(t, result)
	})

	t.Run("Should upload document and return URL successfully", func(t *testing.T) {
		expected := &models.StudentDocumentUrl{Url: "https://example.com/document.pdf"}
		mockStudentRepository.EXPECT().ExistsById(ctx, id).Return(&exists, nil).MaxTimes(1)
		mockDocumentStorage.EXPECT().Upload(ctx, id, &file).Return(expected.Url, nil)

		result, err := usecase.Execute(ctx, id, &file)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
	})
}
