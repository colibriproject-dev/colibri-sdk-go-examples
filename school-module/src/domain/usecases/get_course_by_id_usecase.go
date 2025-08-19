//go:generate mockgen -source get_course_by_id_usecase.go -destination mock/get_course_by_id_usecase_mock.go -package usecasesmock
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
	errAnErrorOccurredInGetCourseByIdUsecaseMsg string = "an error occurred in GetCourseByIdUsecase"
)

type IGetCourseByIdUsecase interface {
	Execute(ctx context.Context, id uuid.UUID) (*models.Course, error)
}

type GetCourseByIdUsecase struct {
	CourseRepository repositories.ICoursesRepository
}

func NewGetCourseByIdUsecase() *GetCourseByIdUsecase {
	return &GetCourseByIdUsecase{
		CourseRepository: repositories.NewCoursesDBRepository(),
	}
}

func (u *GetCourseByIdUsecase) Execute(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	result, err := u.CourseRepository.FindById(ctx, id)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "CourseRepository.FindById").
			AddParam("id", id).
			Msg(errAnErrorOccurredInGetCourseByIdUsecaseMsg)
		return nil, errors.New(exceptions.ErrOnFindCourseById)
	}

	if result == nil {
		return nil, errors.New(exceptions.ErrCourseNotFound)
	}

	return result, nil
}
