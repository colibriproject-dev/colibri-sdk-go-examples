//go:generate mockgen -source get_all_course_usecase.go -destination mock/get_all_course_usecase_mock.go -package usecasesmock
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
	errAnErrorOccurredInGetAllCourseUsecaseMsg string = "an error occurred in GetAllCourseUsecase"
)

type IGetAllCourseUsecase interface {
	Execute(ctx context.Context) ([]models.Course, error)
}

type GetAllCourseUsecase struct {
	CourseRepository repositories.ICoursesRepository
}

func NewGetAllCourseUsecase() *GetAllCourseUsecase {
	return &GetAllCourseUsecase{
		CourseRepository: repositories.NewCoursesDBRepository(),
	}
}

func (u *GetAllCourseUsecase) Execute(ctx context.Context) ([]models.Course, error) {
	result, err := u.CourseRepository.FindAll(ctx)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "CourseRepository.FindAll").
			Msg(errAnErrorOccurredInGetAllCourseUsecaseMsg)
		return nil, errors.New(exceptions.ErrOnFindAllCourses)
	}

	return result, nil
}
