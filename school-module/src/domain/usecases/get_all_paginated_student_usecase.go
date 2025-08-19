//go:generate mockgen -source get_all_paginated_student_usecase.go -destination mock/get_all_paginated_student_usecase_mock.go -package usecasesmock
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
	errAnErrorOccurredInGetAllPaginatedStudentUsecaseMsg string = "an error occurred in GetAllPaginatedStudentUsecase"
)

type IGetAllPaginatedStudentUsecase interface {
	Execute(ctx context.Context, params *models.StudentPageParams) (models.StudentPage, error)
}

type GetAllPaginatedStudentUsecase struct {
	StudentRepository repositories.IStudentsRepository
}

func NewGetAllPaginatedStudentUsecase() *GetAllPaginatedStudentUsecase {
	return &GetAllPaginatedStudentUsecase{
		StudentRepository: repositories.NewStudentsDBRepository(),
	}
}

func (u *GetAllPaginatedStudentUsecase) Execute(ctx context.Context, params *models.StudentPageParams) (models.StudentPage, error) {
	result, err := u.StudentRepository.FindAllPaginated(ctx, params)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "StudentRepository.FindAllPaginated").
			AddParam("params", params).
			Msg(errAnErrorOccurredInGetAllPaginatedStudentUsecaseMsg)
		return nil, errors.New(exceptions.ErrOnFindAllPaginatedStudents)
	}

	return result, nil
}
