//go:generate mockgen -source create_enrollment_usecase.go -destination mock/create_enrollment_usecase_mock.go -package usecasesmock
package usecases

import (
	"context"
	"errors"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/producers"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/logging"
)

const (
	errAnErrorOccurredInCreateEnrollmentUsecaseMsg string = "an error occurred in CreateEnrollmentUsecase"
)

type ICreateEnrollmentUsecase interface {
	Execute(ctx context.Context, model *models.EnrollmentCreate) error
}

type CreateEnrollmentUsecase struct {
	CourseRepository          repositories.ICoursesRepository
	StudentRepository         repositories.IStudentsRepository
	EnrollmentRepository      repositories.IEnrollmentsRepository
	EnrollmentCreatedProducer producers.IEnrollmentCreatedProducer
}

func NewCreateEnrollmentUsecase() *CreateEnrollmentUsecase {
	return &CreateEnrollmentUsecase{
		CourseRepository:          repositories.NewCoursesDBRepository(),
		StudentRepository:         repositories.NewStudentsDBRepository(),
		EnrollmentRepository:      repositories.NewEnrollmentsDBRepository(),
		EnrollmentCreatedProducer: producers.NewEnrollmentCreatedProducer(),
	}
}

func (u *CreateEnrollmentUsecase) Execute(ctx context.Context, model *models.EnrollmentCreate) error {
	model.Status = enums.ADIMPLENTE

	if err := u.existsEnrollmentByStudentIdAndCourseId(ctx, model); err != nil {
		return err
	}

	if err := u.existsCourseById(ctx, model); err != nil {
		return err
	}

	if err := u.existsStudentById(ctx, model); err != nil {
		return err
	}

	result, err := u.insertEnrollment(ctx, model)
	if err != nil {
		return err
	}

	u.sendCreatedEnrollmentNotification(ctx, result)

	return nil
}

func (u *CreateEnrollmentUsecase) existsEnrollmentByStudentIdAndCourseId(ctx context.Context, model *models.EnrollmentCreate) error {
	exists, err := u.EnrollmentRepository.ExistsByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "EnrollmentRepository.ExistsByStudentIdAndCourseId").
			AddParam("model", model).
			Msg(errAnErrorOccurredInCreateEnrollmentUsecaseMsg)
		return errors.New(exceptions.ErrOnExistsEnrollmentByStudentIdAndCourseId)
	}

	if exists != nil && *exists {
		return errors.New(exceptions.ErrEnrollmentAlreadyExists)
	}

	return nil
}

func (u *CreateEnrollmentUsecase) existsCourseById(ctx context.Context, model *models.EnrollmentCreate) error {
	exists, err := u.CourseRepository.ExistsById(ctx, model.CourseID)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "CourseRepository.ExistsById").
			AddParam("model", model).
			Msg(errAnErrorOccurredInCreateEnrollmentUsecaseMsg)
		return errors.New(exceptions.ErrOnExistsCourseById)
	}

	if exists != nil && !*exists {
		return errors.New(exceptions.ErrCourseNotFound)
	}

	return nil
}

func (u *CreateEnrollmentUsecase) existsStudentById(ctx context.Context, model *models.EnrollmentCreate) error {
	exists, err := u.StudentRepository.ExistsById(ctx, model.StudentID)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "StudentRepository.ExistsById").
			AddParam("model", model).
			Msg(errAnErrorOccurredInCreateEnrollmentUsecaseMsg)
		return errors.New(exceptions.ErrOnExistsStudentById)
	}

	if exists != nil && !*exists {
		return errors.New(exceptions.ErrStudentNotFound)
	}

	return nil
}

func (u *CreateEnrollmentUsecase) insertEnrollment(ctx context.Context, model *models.EnrollmentCreate) (*models.EnrollmentCreated, error) {
	result, err := u.EnrollmentRepository.Insert(ctx, model)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "EnrollmentRepository.Insert").
			AddParam("model", model).
			Msg(errAnErrorOccurredInCreateEnrollmentUsecaseMsg)
		return nil, errors.New(exceptions.ErrOnInsertEnrollment)
	}

	return result, nil
}

func (u *CreateEnrollmentUsecase) sendCreatedEnrollmentNotification(ctx context.Context, enrollmentCreated *models.EnrollmentCreated) {
	if err := u.EnrollmentCreatedProducer.Send(ctx, enrollmentCreated); err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "EnrollmentCreatedProducer.Send").
			AddParam("model", enrollmentCreated).
			Msg(errAnErrorOccurredInCreateEnrollmentUsecaseMsg)
	}
}
