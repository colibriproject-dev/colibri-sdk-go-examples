//go:generate mockgen -source create_student_usecase.go -destination mock/create_student_usecase_mock.go -package usecasesmock
package usecases

import (
	"context"
	"errors"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/logging"
)

const (
	errAnErrorOccurredInCreateStudentUsecaseMsg string = "an error occurred in CreateStudentUsecase"
)

type ICreateStudentUsecase interface {
	Execute(ctx context.Context, model *models.StudentCreate) error
}

type CreateStudentUsecase struct {
	Repository repositories.IStudentsRepository
}

func NewCreateStudentUsecase() *CreateStudentUsecase {
	return &CreateStudentUsecase{
		Repository: repositories.NewStudentsDBRepository(),
	}
}

func (u *CreateStudentUsecase) Execute(ctx context.Context, model *models.StudentCreate) error {
	if err := u.existsStudentByEmail(ctx, model); err != nil {
		return err
	}

	return u.insertStudent(ctx, model)
}

func (u *CreateStudentUsecase) existsStudentByEmail(ctx context.Context, model *models.StudentCreate) error {
	exists, err := u.Repository.ExistsByEmail(ctx, model.Email)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "StudentRepository.FindByEmail").
			AddParam("model", model).
			Msg(errAnErrorOccurredInCreateStudentUsecaseMsg)
		return errors.New(exceptions.ErrOnExistsStudentByEmail)
	}

	if exists == nil || *exists {
		return errors.New(exceptions.ErrStudentAlreadyExists)
	}

	return nil
}

func (u *CreateStudentUsecase) insertStudent(ctx context.Context, model *models.StudentCreate) error {
	if err := u.Repository.Insert(ctx, model); err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "StudentRepository.Insert").
			AddParam("model", model).
			Msg(errAnErrorOccurredInCreateStudentUsecaseMsg)
		return errors.New(exceptions.ErrOnInsertStudent)
	}

	return nil
}
