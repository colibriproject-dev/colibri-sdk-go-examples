//go:generate mockgen -source enrollment_usecases.go -destination mock/enrollment_usecases_mock.go -package mock
package usecases

import (
	"context"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/producers"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/google/uuid"
)

type IEnrollmentUsecases interface {
	GetPage(ctx context.Context, params *models.EnrollmentPageParamsDTO) (models.EnrollmentPage, error)
	Create(ctx context.Context, model *models.EnrollmentCreateDTO) error
	Delete(ctx context.Context, params *models.EnrollmentDeleteParamsDTO) error
	UpdateStatus(ctx context.Context, studentID, courseID uuid.UUID, status models.EnrollmentStatus) error
}

type EnrollmentUsecases struct {
	Repository repositories.EnrollmentRepository
	Producer   producers.IEnrollmentProducer
}

func NewEnrollmentUsecases() *EnrollmentUsecases {
	return &EnrollmentUsecases{
		Repository: repositories.NewEnrollmentDBRepository(),
		Producer:   producers.NewEnrollmentProducer(),
	}
}

func (u *EnrollmentUsecases) GetPage(ctx context.Context, params *models.EnrollmentPageParamsDTO) (models.EnrollmentPage, error) {
	return u.Repository.FindPage(ctx, params.ToPageRequest(), params.ToFilters())
}

func (u *EnrollmentUsecases) Create(ctx context.Context, model *models.EnrollmentCreateDTO) error {
	if err := u.Repository.Insert(ctx, model); err != nil {
		return err
	}

	result, _ := u.Repository.FindByStudentIdAndCourseId(ctx, model.StudentID, model.CourseID)
	u.Producer.Create(ctx, result)

	return nil
}

func (u *EnrollmentUsecases) Delete(ctx context.Context, params *models.EnrollmentDeleteParamsDTO) error {
	if err := u.Repository.Delete(ctx, params.StudentID, params.CourseID); err != nil {
		return err
	}

	u.Producer.Delete(ctx, &models.Enrollment{
		Student: models.Student{ID: params.StudentID},
		Course:  models.Course{ID: params.CourseID},
	})

	return nil
}

func (u *EnrollmentUsecases) UpdateStatus(ctx context.Context, studentID, courseID uuid.UUID, status models.EnrollmentStatus) error {
	return u.Repository.UpdateStatus(ctx, studentID, courseID, status)
}
