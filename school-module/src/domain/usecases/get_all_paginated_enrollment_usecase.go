//go:generate mockgen -source get_all_paginated_enrollment_usecase.go -destination mock/get_all_paginated_enrollment_usecase_mock.go -package usecasesmock
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
	errAnErrorOccurredInGetAllPaginatedEnrollmentUsecaseMsg string = "an error occurred in GetAllPaginatedEnrollmentUsecase"
)

type IGetAllPaginatedEnrollmentUsecase interface {
	Execute(ctx context.Context, params *models.EnrollmentPageParams) (models.EnrollmentPage, error)
}

type GetAllPaginatedEnrollmentUsecase struct {
	EnrollmentRepository repositories.IEnrollmentsRepository
}

func NewGetAllPaginatedEnrollmentUsecase() *GetAllPaginatedEnrollmentUsecase {
	return &GetAllPaginatedEnrollmentUsecase{
		EnrollmentRepository: repositories.NewEnrollmentsDBRepository(),
	}
}

func (u *GetAllPaginatedEnrollmentUsecase) Execute(ctx context.Context, params *models.EnrollmentPageParams) (models.EnrollmentPage, error) {
	result, err := u.EnrollmentRepository.FindAllPaginated(ctx, params)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "EnrollmentRepository.FindAllPaginated").
			AddParam("params", params).
			Msg(errAnErrorOccurredInGetAllPaginatedEnrollmentUsecaseMsg)
		return nil, errors.New(exceptions.ErrOnFindAllPaginatedEnrollment)
	}

	return result, nil
}
