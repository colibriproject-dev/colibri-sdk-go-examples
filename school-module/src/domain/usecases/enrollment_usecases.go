package usecases

import (
	"context"
	"school-module/src/domain/models"
	"school-module/src/infra/producers"
	"school-module/src/infra/repositories"

	"github.com/google/uuid"
)

type IEnrollmentUsecase interface {
	GetAll(ctx context.Context, studentName, courseName string) ([]models.Enrollment, error)
	Create(ctx context.Context, model *models.EnrollmentCreateDTO) error
	Delete(ctx context.Context, studentId, courseId uuid.UUID) error
	UpdateStatus(ctx context.Context, student_id, course_id uuid.UUID, status models.EnrollmentStatus) error
}

type EnrollmentUsecase struct {
	Repository repositories.IEnrollmentRepository
	Producer   producers.IEnrollmentProducer
}

func NewEnrollmentUsecase() *EnrollmentUsecase {
	return &EnrollmentUsecase{
		Repository: repositories.NewEnrollmentRepository(),
		Producer:   producers.NewEnrollmentProducer(),
	}
}

func (u *EnrollmentUsecase) GetAll(ctx context.Context, studentName, courseName string) ([]models.Enrollment, error) {
	return u.Repository.FindAll(ctx, studentName, courseName)
}

func (u *EnrollmentUsecase) Create(ctx context.Context, model *models.EnrollmentCreateDTO) error {
	if err := u.Repository.Insert(ctx, model); err != nil {
		return err
	}

	result, _ := u.Repository.FindByStudentIdAndCourseId(ctx, model.Student.ID, model.Course.ID)
	u.Producer.Create(ctx, result)

	return nil
}

func (u *EnrollmentUsecase) Delete(ctx context.Context, studentId, courseId uuid.UUID) error {
	if err := u.Repository.Delete(ctx, studentId, courseId); err != nil {
		return err
	}

	u.Producer.Delete(ctx, &models.Enrollment{
		Student: models.Student{ID: studentId},
		Course:  models.Course{ID: courseId},
	})

	return nil
}

func (u *EnrollmentUsecase) UpdateStatus(ctx context.Context, student_id, course_id uuid.UUID, status models.EnrollmentStatus) error {
	return u.Repository.UpdateStatus(ctx, student_id, course_id, status)
}
