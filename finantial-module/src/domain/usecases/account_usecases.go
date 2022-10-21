package usecases

import (
	"context"
	"finantial-module/src/domain/models"
	"finantial-module/src/infra/repositories"
	"time"

	"github.com/google/uuid"
)

type AccountUsecases interface {
	GetAll(ctx context.Context) ([]models.Account, error)
	Create(ctx context.Context, model *models.Account) error
	DeleteByStudentAndCourse(ctx context.Context, studentId, courseId uuid.UUID) error
	DeleteByCourse(ctx context.Context, courseId uuid.UUID) error
	DeleteByStudent(ctx context.Context, studentId uuid.UUID) error
}

type AccountUsecase struct {
	InvoiceUsecases InvoiceUsecases
	Repository      repositories.AccountRepository
}

func NewAccountUsecase() *AccountUsecase {
	return &AccountUsecase{
		InvoiceUsecases: NewInvoiceUsecase(),
		Repository:      repositories.NewAccountDBRepository(),
	}
}

func (u *AccountUsecase) GetAll(ctx context.Context) ([]models.Account, error) {
	return u.Repository.FindAll(ctx)
}

func (u *AccountUsecase) Create(ctx context.Context, model *models.Account) error {
	model.ID = uuid.New()
	model.Status = models.ADIMPLENTE
	model.CreatedAt = time.Now()

	if err := u.Repository.Insert(ctx, model); err != nil {
		return err
	}

	if err := u.InvoiceUsecases.Create(ctx, model); err != nil {
		return err
	}

	return nil
}

func (u *AccountUsecase) DeleteByStudentAndCourse(ctx context.Context, studentId, courseId uuid.UUID) error {
	return u.Repository.DeleteByStudentAndCourse(ctx, studentId, courseId)
}

func (u *AccountUsecase) DeleteByCourse(ctx context.Context, courseId uuid.UUID) error {
	return u.Repository.DeleteByCourse(ctx, courseId)
}

func (u *AccountUsecase) DeleteByStudent(ctx context.Context, studentId uuid.UUID) error {
	return u.Repository.DeleteByStudent(ctx, studentId)
}
