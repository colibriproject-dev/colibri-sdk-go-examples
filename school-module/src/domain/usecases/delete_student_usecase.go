//go:generate mockgen -source delete_student_usecase.go -destination mock/delete_student_usecase_mock.go -package usecasesmock
package usecases

import (
	"context"
	"errors"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/producers"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/logging"
	"github.com/google/uuid"
)

const (
	errAnErrorOccurredInDeleteStudentUsecaseMsg string = "an error occurred in DeleteStudentUsecase"
)

type IDeleteStudentUsecase interface {
	Execute(ctx context.Context, id uuid.UUID) error
}

type DeleteStudentUsecase struct {
	Repository             repositories.IStudentsRepository
	StudentDeletedProducer producers.IStudentDeletedProducer
}

func NewDeleteStudentUsecase() *DeleteStudentUsecase {
	return &DeleteStudentUsecase{
		Repository:             repositories.NewStudentsDBRepository(),
		StudentDeletedProducer: producers.NewStudentDeletedProducer(),
	}
}

func (u *DeleteStudentUsecase) Execute(ctx context.Context, id uuid.UUID) error {
	if err := u.existsStudentById(ctx, id); err != nil {
		return err
	}

	if err := u.deleteStudentById(ctx, id); err != nil {
		return err
	}

	u.sendDeletedStudentNotification(ctx, id)

	return nil
}

func (u *DeleteStudentUsecase) existsStudentById(ctx context.Context, id uuid.UUID) error {
	exists, err := u.Repository.ExistsById(ctx, id)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "StudentRepository.ExistsById").
			AddParam("id", id).
			Msg(errAnErrorOccurredInDeleteStudentUsecaseMsg)
		return errors.New(exceptions.ErrOnExistsStudentById)
	}

	if exists == nil || !*exists {
		return errors.New(exceptions.ErrStudentNotFound)
	}

	return nil
}

func (u *DeleteStudentUsecase) deleteStudentById(ctx context.Context, id uuid.UUID) error {
	if err := u.Repository.Delete(ctx, id); err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "StudentRepository.Delete").
			AddParam("id", id).
			Msg(errAnErrorOccurredInDeleteStudentUsecaseMsg)
		return errors.New(exceptions.ErrOnDeleteStudent)
	}

	return nil
}

func (u *DeleteStudentUsecase) sendDeletedStudentNotification(ctx context.Context, id uuid.UUID) {
	if err := u.StudentDeletedProducer.Send(ctx, &models.StudentDelete{ID: id}); err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "StudentDeletedProducer.Send").
			AddParam("id", id).
			Msg(errAnErrorOccurredInDeleteStudentUsecaseMsg)
	}
}
