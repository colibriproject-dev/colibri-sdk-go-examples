package usecases

import (
	"context"
	"mime/multipart"
	"school-module/src/domain/models"
	"school-module/src/infra/producers"
	"school-module/src/infra/repositories"

	"github.com/google/uuid"
)

type IStudentUsecase interface {
	GetAll(ctx context.Context, params *models.StudentParams) ([]models.Student, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Student, error)
	Create(ctx context.Context, model *models.StudentCreateUpdateDTO) error
	Update(ctx context.Context, id uuid.UUID, model *models.StudentCreateUpdateDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
	UploadDocument(ctx context.Context, id uuid.UUID, file *multipart.File) (string, error)
}

type StudentUsecase struct {
	Repository repositories.IStudentRepository
	Producer   producers.IStudentProducer
}

func NewStudentUsecase() *StudentUsecase {
	return &StudentUsecase{
		Repository: repositories.NewStudentRepository(),
		Producer:   producers.NewStudentProducer(),
	}
}

func (u *StudentUsecase) GetAll(ctx context.Context, params *models.StudentParams) ([]models.Student, error) {
	return u.Repository.FindAll(ctx, params)
}

func (u *StudentUsecase) GetById(ctx context.Context, id uuid.UUID) (*models.Student, error) {
	return u.Repository.FindById(ctx, id)
}

func (u *StudentUsecase) Create(ctx context.Context, model *models.StudentCreateUpdateDTO) error {
	return u.Repository.Insert(ctx, model)
}

func (u *StudentUsecase) Update(ctx context.Context, id uuid.UUID, model *models.StudentCreateUpdateDTO) error {
	return u.Repository.Update(ctx, id, model)
}

func (u *StudentUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.Repository.Delete(ctx, id); err != nil {
		return err
	}

	u.Producer.Delete(ctx, &models.Student{ID: id})

	return nil
}

func (u *StudentUsecase) UploadDocument(ctx context.Context, id uuid.UUID, file *multipart.File) (string, error) {
	return u.Repository.UploadDocument(ctx, id, file)
}
