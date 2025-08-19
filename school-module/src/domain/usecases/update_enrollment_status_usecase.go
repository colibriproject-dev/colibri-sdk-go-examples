//go:generate mockgen -source update_enrollment_status_usecase.go -destination mock/update_enrollment_status_usecase_mock.go -package usecasesmock
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
	errAnErrorOccurredInUpdateEnrollmentStatusUsecaseMsg string = "an error occurred in UpdateEnrollmentStatusUsecase"
)

type IUpdateEnrollmentStatusUsecase interface {
	Execute(ctx context.Context, model *models.EnrollmentUpdateStatus) error
}

type UpdateEnrollmentStatusUsecase struct {
	EnrollmentRepository repositories.IEnrollmentsRepository
}

func NewUpdateEnrollmentStatusUsecase() *UpdateEnrollmentStatusUsecase {
	return &UpdateEnrollmentStatusUsecase{
		EnrollmentRepository: repositories.NewEnrollmentsDBRepository(),
	}
}

func (u *UpdateEnrollmentStatusUsecase) Execute(ctx context.Context, model *models.EnrollmentUpdateStatus) error {
	if err := u.existsEnrollmentByStudentIdAndCourseId(ctx, model); err != nil {
		return err
	}

	return u.updateEnrollmentStatus(ctx, model)
}

func (u *UpdateEnrollmentStatusUsecase) existsEnrollmentByStudentIdAndCourseId(ctx context.Context, model *models.EnrollmentUpdateStatus) error {
	exists, err := u.EnrollmentRepository.ExistsByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "EnrollmentRepository.ExistsByStudentIdAndCourseId").
			AddParam("model", model).
			Msg(errAnErrorOccurredInUpdateEnrollmentStatusUsecaseMsg)
		return errors.New(exceptions.ErrOnExistsEnrollmentByStudentIdAndCourseId)
	}

	if exists == nil || !*exists {
		return errors.New(exceptions.ErrEnrollmentNotFound)
	}

	return nil
}

func (u *UpdateEnrollmentStatusUsecase) updateEnrollmentStatus(ctx context.Context, model *models.EnrollmentUpdateStatus) error {
	if err := u.EnrollmentRepository.UpdateStatus(ctx, model); err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "EnrollmentRepository.UpdateStatus").
			AddParam("model", model).
			Msg(errAnErrorOccurredInUpdateEnrollmentStatusUsecaseMsg)
		return errors.New(exceptions.ErrOnUpdateEnrollmentStatus)
	}

	return nil
}
