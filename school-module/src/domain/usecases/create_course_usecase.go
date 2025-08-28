//go:generate mockgen -source create_course_usecase.go -destination mock/create_course_usecase_mock.go -package usecasesmock
package usecases

import (
	"context"
	"errors"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/producers"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/logging"
)

const (
	errAnErrorOccurredInCreateCourseUsecaseMsg string = "an error occurred in CreateCourseUsecase"
)

type ICreateCourseUsecase interface {
	Execute(ctx context.Context, model *models.CourseCreate) (*models.Course, error)
}

type CreateCourseUsecase struct {
	CourseRepository repositories.ICoursesRepository
	createdProducer  producers.ICourseCreatedProducer
}

func NewCreateCourseUsecase() *CreateCourseUsecase {
	return &CreateCourseUsecase{
		CourseRepository: repositories.NewCoursesDBRepository(),
		createdProducer:  producers.NewCourseCreatedProducer(),
	}
}

func (u *CreateCourseUsecase) Execute(ctx context.Context, model *models.CourseCreate) (*models.Course, error) {
	if err := u.existsCourseByName(ctx, model); err != nil {
		return nil, err
	}

	return u.insertCourse(ctx, model)
}

func (u *CreateCourseUsecase) existsCourseByName(ctx context.Context, model *models.CourseCreate) error {
	exists, err := u.CourseRepository.ExistsByName(ctx, model.Name)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "CourseRepository.ExistsByName").
			AddParam("model", model).
			Msg(errAnErrorOccurredInCreateCourseUsecaseMsg)
		return errors.New(exceptions.ErrOnExistsCourseByName)
	}

	if exists != nil && *exists {
		return errors.New(exceptions.ErrCourseAlreadyExists)
	}

	return nil
}

func (u *CreateCourseUsecase) insertCourse(ctx context.Context, model *models.CourseCreate) (*models.Course, error) {
	result, err := u.CourseRepository.Insert(ctx, model)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "CourseRepository.Insert").
			AddParam("model", model).
			Msg(errAnErrorOccurredInCreateCourseUsecaseMsg)
		return nil, errors.New(exceptions.ErrOnInsertCourse)
	}

	if err := u.createdProducer.Send(ctx, result); err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "CourseCreatedProducer.Send").
			AddParam("model", result)
	}

	return result, nil
}
