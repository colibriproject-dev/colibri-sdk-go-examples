//go:generate mockgen -source update_student_usecase.go -destination mock/update_student_usecase_mock.go -package usecasesmock
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
	errAnErrorOccurredInUpdateStudentUsecaseMsg string = "an error occurred in UpdateStudentUsecase"
)

type IUpdateStudentUsecase interface {
	Execute(ctx context.Context, model *models.StudentUpdate) error
}

type UpdateStudentUsecase struct {
	StudentRepository repositories.IStudentsRepository
}

func NewUpdateStudentUsecase() *UpdateStudentUsecase {
	return &UpdateStudentUsecase{
		StudentRepository: repositories.NewStudentsDBRepository(),
	}
}

func (u *UpdateStudentUsecase) Execute(ctx context.Context, model *models.StudentUpdate) error {
	if err := u.existsStudentById(ctx, model); err != nil {
		return err
	}
	if err := u.findStudentByEmail(ctx, model); err != nil {
		return err
	}

	return u.updateStudent(ctx, model)
}

func (u *UpdateStudentUsecase) existsStudentById(ctx context.Context, model *models.StudentUpdate) error {
	exists, err := u.StudentRepository.ExistsById(ctx, model.ID)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "StudentRepository.ExistsById").
			AddParam("model", model).
			Msg(errAnErrorOccurredInUpdateStudentUsecaseMsg)
		return errors.New(exceptions.ErrOnExistsStudentById)
	}

	if (exists == nil) || !*exists {
		return errors.New(exceptions.ErrStudentNotFound)
	}

	return nil
}

func (u *UpdateStudentUsecase) findStudentByEmail(ctx context.Context, model *models.StudentUpdate) error {
	result, err := u.StudentRepository.FindByEmail(ctx, model.Email)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "StudentRepository.FindByEmail").
			AddParam("model", model).
			Msg(errAnErrorOccurredInUpdateStudentUsecaseMsg)
		return errors.New(exceptions.ErrOnFindStudentByEmail)
	}

	if (result != nil) && (result.Name == model.Name) && (result.ID != model.ID) {
		return errors.New(exceptions.ErrStudentAlreadyExists)
	}

	return nil
}

func (u *UpdateStudentUsecase) updateStudent(ctx context.Context, model *models.StudentUpdate) error {
	if err := u.StudentRepository.Update(ctx, model); err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "StudentRepository.Update").
			AddParam("model", model).
			Msg(errAnErrorOccurredInUpdateStudentUsecaseMsg)
		return errors.New(exceptions.ErrOnUpdateStudent)
	}

	return nil
}
