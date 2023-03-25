//go:generate mockgen -source student_usecases.go -destination mock/student_usecases_mock.go -package mock
package usecases

import (
	"context"
	"mime/multipart"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/producers"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/google/uuid"
)

type IStudentUsecases interface {
	GetAll(ctx context.Context, params *models.StudentParams) ([]models.Student, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Student, error)
	Create(ctx context.Context, model *models.StudentCreateUpdateDTO) error
	Update(ctx context.Context, id uuid.UUID, model *models.StudentCreateUpdateDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
	UploadDocument(ctx context.Context, id uuid.UUID, file *multipart.File) (string, error)
}

type StudentUsecases struct {
	Repository repositories.StudentRepository
	Producer   producers.IStudentDeletedProducer
}

func NewStudentUsecases() *StudentUsecases {
	return &StudentUsecases{
		Repository: repositories.NewStudentDBRepository(),
		Producer:   producers.NewStudentDeletedProducer(),
	}
}

func (u *StudentUsecases) GetAll(ctx context.Context, params *models.StudentParams) ([]models.Student, error) {
	return u.Repository.FindAll(ctx, params)
}

func (u *StudentUsecases) GetById(ctx context.Context, id uuid.UUID) (*models.Student, error) {
	return u.Repository.FindById(ctx, id)
}

func (u *StudentUsecases) Create(ctx context.Context, model *models.StudentCreateUpdateDTO) error {
	return u.Repository.Insert(ctx, model)
}

func (u *StudentUsecases) Update(ctx context.Context, id uuid.UUID, model *models.StudentCreateUpdateDTO) error {
	return u.Repository.Update(ctx, id, model)
}

func (u *StudentUsecases) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.Repository.Delete(ctx, id); err != nil {
		return err
	}

	u.Producer.Delete(ctx, &models.Student{ID: id})

	return nil
}

func (u *StudentUsecases) UploadDocument(ctx context.Context, id uuid.UUID, file *multipart.File) (string, error) {
	return u.Repository.UploadDocument(ctx, id, file)
}
