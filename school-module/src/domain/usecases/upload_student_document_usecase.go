//go:generate mockgen -source upload_student_document_usecase.go -destination mock/upload_student_document_usecase_mock.go -package usecasesmock
package usecases

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/storages"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/logging"
	"github.com/google/uuid"
)

const (
	errAnErrorOccurredInUploadStudentDocumentUsecaseMsg string = "an error occurred in UploadStudentDocumentUsecase"
)

type IUploadStudentDocumentUsecase interface {
	Execute(ctx context.Context, id uuid.UUID, file *multipart.File) (*models.StudentDocumentUrl, error)
}

type UploadStudentDocumentUsecase struct {
	StudentRepository repositories.IStudentsRepository
	DocumentStorage   storages.IDocumentStorage
}

func NewUploadStudentDocumentUsecase() *UploadStudentDocumentUsecase {
	return &UploadStudentDocumentUsecase{
		StudentRepository: repositories.NewStudentsDBRepository(),
		DocumentStorage:   storages.NewDocumentS3Storage(),
	}
}

func (u *UploadStudentDocumentUsecase) Execute(ctx context.Context, id uuid.UUID, file *multipart.File) (*models.StudentDocumentUrl, error) {
	if err := u.existsStudentById(ctx, id); err != nil {
		return nil, err
	}

	return u.uploadStudentDocument(ctx, id, file)
}

func (u *UploadStudentDocumentUsecase) existsStudentById(ctx context.Context, id uuid.UUID) error {
	exists, err := u.StudentRepository.ExistsById(ctx, id)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "StudentRepository.ExistsById").
			AddParam("id", id).
			Msg(errAnErrorOccurredInUploadStudentDocumentUsecaseMsg)
		return errors.New(exceptions.ErrOnExistsStudentById)
	}

	if (exists == nil) || !*exists {
		return errors.New(exceptions.ErrStudentNotFound)
	}

	return nil
}

func (u *UploadStudentDocumentUsecase) uploadStudentDocument(ctx context.Context, id uuid.UUID, file *multipart.File) (*models.StudentDocumentUrl, error) {
	result, err := u.DocumentStorage.Upload(ctx, id, file)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "DocumentStorage.Upload").
			AddParam("id", id).
			Msg(errAnErrorOccurredInUploadStudentDocumentUsecaseMsg)
		return nil, errors.New(exceptions.ErrOnUploadStudentDocument)
	}

	return &models.StudentDocumentUrl{Url: result}, nil
}
