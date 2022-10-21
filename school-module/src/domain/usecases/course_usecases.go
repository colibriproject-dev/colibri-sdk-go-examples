package usecases

import (
	"context"
	"school-module/src/domain/models"
	"school-module/src/infra/producers"
	"school-module/src/infra/repositories"

	"github.com/google/uuid"
)

type ICourseUsecase interface {
	GetAll(ctx context.Context) ([]models.Course, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Course, error)
	Create(ctx context.Context, model *models.CourseCreateUpdateDTO) (*models.Course, error)
	Update(ctx context.Context, id uuid.UUID, model *models.CourseCreateUpdateDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CourseUsecase struct {
	Repository repositories.ICourseRepository
	Producer   producers.ICourseProducer
}

func NewCourseUsecase() *CourseUsecase {
	return &CourseUsecase{
		Repository: repositories.NewCourseRepository(),
		Producer:   producers.NewCourseProducer(),
	}
}

func (u *CourseUsecase) GetAll(ctx context.Context) ([]models.Course, error) {
	return u.Repository.FindAll(ctx)
}

func (u *CourseUsecase) GetById(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	return u.Repository.FindById(ctx, id)
}

func (u *CourseUsecase) Create(ctx context.Context, model *models.CourseCreateUpdateDTO) (*models.Course, error) {
	return u.Repository.Insert(ctx, model)
}

func (u *CourseUsecase) Update(ctx context.Context, id uuid.UUID, model *models.CourseCreateUpdateDTO) error {
	return u.Repository.Update(ctx, id, model)
}

func (u *CourseUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.Repository.Delete(ctx, id); err != nil {
		return err
	}

	u.Producer.Delete(ctx, &models.Course{ID: id})

	return nil
}
