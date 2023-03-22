//go:generate mockgen -source course_usecases.go -destination mock/course_usecases_mock.go -package mock
package usecases

import (
	"context"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/producers"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/google/uuid"
)

type ICourseUsecases interface {
	GetAll(ctx context.Context) ([]models.Course, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Course, error)
	Create(ctx context.Context, model *models.CourseCreateUpdateDTO) (*models.Course, error)
	Update(ctx context.Context, id uuid.UUID, model *models.CourseCreateUpdateDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CourseUsecases struct {
	Repository repositories.CourseRepository
	Producer   producers.ICourseProducer
}

func NewCourseUsecases() *CourseUsecases {
	return &CourseUsecases{
		Repository: repositories.NewCourseDBRepository(),
		Producer:   producers.NewCourseProducer(),
	}
}

func (u *CourseUsecases) GetAll(ctx context.Context) ([]models.Course, error) {
	return u.Repository.FindAll(ctx)
}

func (u *CourseUsecases) GetById(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	return u.Repository.FindById(ctx, id)
}

func (u *CourseUsecases) Create(ctx context.Context, model *models.CourseCreateUpdateDTO) (*models.Course, error) {
	return u.Repository.Insert(ctx, model)
}

func (u *CourseUsecases) Update(ctx context.Context, id uuid.UUID, model *models.CourseCreateUpdateDTO) error {
	return u.Repository.Update(ctx, id, model)
}

func (u *CourseUsecases) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.Repository.Delete(ctx, id); err != nil {
		return err
	}

	u.Producer.Delete(ctx, &models.Course{ID: id})

	return nil
}
