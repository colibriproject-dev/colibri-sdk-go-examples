//go:generate mockgen -source get_student_by_id_usecase.go -destination mock/get_student_by_id_usecase_mock.go -package usecasesmock
package usecases

import (
	"context"
	"errors"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/logging"
	"github.com/google/uuid"
)

const (
	errAnErrorOccurredInGetStudentByIdUsecaseMsg string = "an error occurred in GetStudentByIdUsecase"
)

type IGetStudentByIdUsecase interface {
	Execute(ctx context.Context, id uuid.UUID) (*models.Student, error)
}

type GetStudentByIdUsecase struct {
	StudentRepository repositories.IStudentsRepository
}

func NewGetStudentByIdUsecase() *GetStudentByIdUsecase {
	return &GetStudentByIdUsecase{
		StudentRepository: repositories.NewStudentsDBRepository(),
	}
}

func (u *GetStudentByIdUsecase) Execute(ctx context.Context, id uuid.UUID) (*models.Student, error) {
	result, err := u.StudentRepository.FindById(ctx, id)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "StudentRepository.FindById").
			AddParam("id", id).
			Msg(errAnErrorOccurredInGetStudentByIdUsecaseMsg)
		return nil, errors.New(exceptions.ErrOnFindStudentById)
	}

	if result == nil {
		return nil, errors.New(exceptions.ErrStudentNotFound)
	}

	return result, nil
}
