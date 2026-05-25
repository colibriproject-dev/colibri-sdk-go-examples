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
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/monitoring"
	colibrimonitoringbase "github.com/colibriproject-dev/colibri-sdk-go/pkg/base/monitoring/colibri-monitoring-base"
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
	counterSuccess   colibrimonitoringbase.Counter
	counterErr       colibrimonitoringbase.Counter
}

func NewCreateCourseUsecase() *CreateCourseUsecase {
	c := monitoring.Counter("school.course.creation.success", "Contador de criação de cursos com sucesso", "1")
	cErr := monitoring.Counter("school.course.creation.error", "Contador de criação de cursos com erro", "1")
	return &CreateCourseUsecase{
		CourseRepository: repositories.NewCoursesDBRepository(),
		createdProducer:  producers.NewCourseCreatedProducer(),
		counterSuccess:   c,
		counterErr:       cErr,
	}
}

func (u *CreateCourseUsecase) Execute(ctx context.Context, model *models.CourseCreate) (*models.Course, error) {
	if err := u.existsCourseByName(ctx, model); err != nil {
		u.counterErr.Add(ctx, 1, nil)
		return nil, err
	}

	course, err := u.insertCourse(ctx, model)
	if err != nil {
		u.counterErr.Add(ctx, 1, nil)
		return nil, err
	}

	u.counterSuccess.Add(ctx, 1, nil)
	return course, nil
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

	logging.Info(ctx).
		AddParam("model", result).
		Msg("course created")

	if err := u.createdProducer.Send(ctx, result); err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "CourseCreatedProducer.Send").
			AddParam("model", result)
	}

	return result, nil
}
