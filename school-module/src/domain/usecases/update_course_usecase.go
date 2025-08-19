//go:generate mockgen -source update_course_usecase.go -destination mock/update_course_usecase_mock.go -package usecasesmock
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
	errAnErrorOccurredInUpdateCourseUsecaseMsg string = "an error occurred in UpdateCourseUsecase"
)

type IUpdateCourseUsecase interface {
	Execute(ctx context.Context, model *models.CourseUpdate) error
}

type UpdateCourseUsecase struct {
	CourseRepository repositories.ICoursesRepository
}

func NewUpdateCourseUsecase() *UpdateCourseUsecase {
	return &UpdateCourseUsecase{
		CourseRepository: repositories.NewCoursesDBRepository(),
	}
}

func (u *UpdateCourseUsecase) Execute(ctx context.Context, model *models.CourseUpdate) error {
	if err := u.existsCourseById(ctx, model); err != nil {
		return err
	}

	if err := u.findCourseByName(ctx, model); err != nil {
		return err
	}

	return u.UpdateCourse(ctx, model)
}

func (u *UpdateCourseUsecase) existsCourseById(ctx context.Context, model *models.CourseUpdate) error {
	exists, err := u.CourseRepository.ExistsById(ctx, model.ID)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "CourseRepository.ExistsById").
			AddParam("model", model).
			Msg(errAnErrorOccurredInUpdateCourseUsecaseMsg)
		return errors.New(exceptions.ErrOnExistsCourseById)
	}

	if (exists == nil) || !*exists {
		return errors.New(exceptions.ErrCourseNotFound)
	}

	return nil
}

func (u *UpdateCourseUsecase) findCourseByName(ctx context.Context, model *models.CourseUpdate) error {
	result, err := u.CourseRepository.FindByName(ctx, model.Name)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "CourseRepository.FindByName").
			AddParam("model", model).
			Msg(errAnErrorOccurredInUpdateCourseUsecaseMsg)
		return errors.New(exceptions.ErrOnFindCourseByName)
	}

	if (result != nil) && (result.Name == model.Name) && (result.ID != model.ID) {
		return errors.New(exceptions.ErrCourseAlreadyExists)
	}

	return nil
}

func (u *UpdateCourseUsecase) UpdateCourse(ctx context.Context, model *models.CourseUpdate) error {
	if err := u.CourseRepository.Update(ctx, model); err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "CourseRepository.Update").
			AddParam("model", model).
			Msg(errAnErrorOccurredInUpdateCourseUsecaseMsg)
		return errors.New(exceptions.ErrOnUpdateCourse)
	}

	return nil
}
