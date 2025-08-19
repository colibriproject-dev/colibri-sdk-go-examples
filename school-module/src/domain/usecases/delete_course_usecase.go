//go:generate mockgen -source delete_course_usecase.go -destination mock/delete_course_usecase_mock.go -package usecasesmock
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
	errAnErrorOccurredInDeleteCourseUsecaseMsg string = "an error occurred in DeleteCourseUsecase"
)

type IDeleteCourseUsecase interface {
	Execute(ctx context.Context, id uuid.UUID) error
}

type DeleteCourseUsecase struct {
	CourseRepository      repositories.ICoursesRepository
	CourseDeletedProducer producers.ICourseDeletedProducer
}

func NewDeleteCourseUsecase() *DeleteCourseUsecase {
	return &DeleteCourseUsecase{
		CourseRepository:      repositories.NewCoursesDBRepository(),
		CourseDeletedProducer: producers.NewCourseDeletedProducer(),
	}
}

func (u *DeleteCourseUsecase) Execute(ctx context.Context, id uuid.UUID) error {
	if err := u.existsCourseById(ctx, id); err != nil {
		return err
	}

	if err := u.deleteCourseById(ctx, id); err != nil {
		return err
	}

	u.sendDeletedCourseNotification(ctx, id)

	return nil
}

func (u *DeleteCourseUsecase) existsCourseById(ctx context.Context, id uuid.UUID) error {
	exists, err := u.CourseRepository.ExistsById(ctx, id)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "CourseRepository.ExistsById").
			AddParam("id", id).
			Msg(errAnErrorOccurredInDeleteCourseUsecaseMsg)
		return errors.New(exceptions.ErrOnExistsCourseById)
	}

	if (exists == nil) || !*exists {
		return errors.New(exceptions.ErrCourseNotFound)
	}

	return nil
}

func (u *DeleteCourseUsecase) deleteCourseById(ctx context.Context, id uuid.UUID) error {
	if err := u.CourseRepository.Delete(ctx, id); err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "CourseRepository.Delete").
			AddParam("id", id).
			Msg(errAnErrorOccurredInDeleteCourseUsecaseMsg)
		return errors.New(exceptions.ErrOnDeleteCourse)
	}

	return nil
}

func (u *DeleteCourseUsecase) sendDeletedCourseNotification(ctx context.Context, id uuid.UUID) {
	if err := u.CourseDeletedProducer.Send(ctx, &models.CourseDelete{ID: id}); err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "CourseDeletedProducer.Send").
			AddParam("id", id).
			Msg(errAnErrorOccurredInDeleteCourseUsecaseMsg)
	}
}
