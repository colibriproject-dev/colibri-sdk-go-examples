//go:generate mockgen -source delete_enrollment_usecase.go -destination mock/delete_enrollment_usecase_mock.go -package usecasesmock
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
	errAnErrorOccurredInDeleteEnrollmentUsecaseMsg string = "an error occurred in DeleteEnrollmentUsecase"
)

type IDeleteEnrollmentUsecase interface {
	Execute(ctx context.Context, params *models.EnrollmentDelete) error
}

type DeleteEnrollmentUsecase struct {
	EnrollmentRepository      repositories.IEnrollmentsRepository
	EnrollmentDeletedProducer producers.IEnrollmentDeletedProducer
}

func NewDeleteEnrollmentUsecase() *DeleteEnrollmentUsecase {
	return &DeleteEnrollmentUsecase{
		EnrollmentRepository:      repositories.NewEnrollmentsDBRepository(),
		EnrollmentDeletedProducer: producers.NewEnrollmentDeletedProducer(),
	}
}

func (u *DeleteEnrollmentUsecase) Execute(ctx context.Context, params *models.EnrollmentDelete) error {
	if err := u.existsEnrollmentByStudentIdAndCourseId(ctx, params); err != nil {
		return err
	}

	if err := u.deleteEnrollment(ctx, params); err != nil {
		return err
	}

	u.sendDeletedEnrollmentNotification(ctx, params)

	return nil
}

func (u *DeleteEnrollmentUsecase) existsEnrollmentByStudentIdAndCourseId(ctx context.Context, params *models.EnrollmentDelete) error {
	exists, err := u.EnrollmentRepository.ExistsByStudentIdAndCourseId(ctx, params.StudentID, params.CourseID)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "EnrollmentRepository.ExistsByStudentIdAndCourseId").
			AddParam("params", params).
			Msg(errAnErrorOccurredInDeleteEnrollmentUsecaseMsg)
		return errors.New(exceptions.ErrOnExistsEnrollmentByStudentIdAndCourseId)
	}

	if exists == nil || !*exists {
		return errors.New(exceptions.ErrEnrollmentNotFound)
	}

	return nil
}

func (u *DeleteEnrollmentUsecase) deleteEnrollment(ctx context.Context, params *models.EnrollmentDelete) error {
	if err := u.EnrollmentRepository.Delete(ctx, params.StudentID, params.CourseID); err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "EnrollmentRepository.Delete").
			AddParam("params", params).
			Msg(errAnErrorOccurredInDeleteEnrollmentUsecaseMsg)
		return errors.New(exceptions.ErrOnDeleteEnrollment)
	}

	return nil
}

func (u *DeleteEnrollmentUsecase) sendDeletedEnrollmentNotification(ctx context.Context, params *models.EnrollmentDelete) {
	if err := u.EnrollmentDeletedProducer.Send(ctx, params); err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "EnrollmentDeletedProducer.Send").
			AddParam("params", params).
			Msg(errAnErrorOccurredInDeleteEnrollmentUsecaseMsg)
	}
}
