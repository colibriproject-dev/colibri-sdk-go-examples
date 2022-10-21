package repositories

import (
	"context"
	"fmt"
	"school-module/src/domain/models"
	"strings"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/google/uuid"
)

type IEnrollmentRepository interface {
	FindAll(ctx context.Context, studentName, courseName string) ([]models.Enrollment, error)
	FindByStudentIdAndCourseId(ctx context.Context, studentId, courseId uuid.UUID) (*models.Enrollment, error)
	Insert(ctx context.Context, model *models.EnrollmentCreateDTO) error
	Delete(ctx context.Context, student_id, course_id uuid.UUID) error
	UpdateStatus(ctx context.Context, student_id, course_id uuid.UUID, status models.EnrollmentStatus) error
}

type EnrollmentRepository struct{}

func NewEnrollmentRepository() *EnrollmentRepository {
	return &EnrollmentRepository{}
}

func (r *EnrollmentRepository) FindAll(ctx context.Context, studentName, courseName string) ([]models.Enrollment, error) {
	if studentName != "" {
		studentName = fmt.Sprintf("%%%s%%", strings.ToLower(studentName))
	}

	if courseName != "" {
		courseName = fmt.Sprintf("%%%s%%", strings.ToLower(courseName))
	}

	const query = `SELECT 
	  s.id, s.name, s.email, s.birthday, s.created_at,
	  c.id, c.name, c.value, c.created_at,
	  e.installments, e.status, e.created_at
	FROM enrollments e
	JOIN students s ON e.student_id = s.id
	JOIN courses c ON e.course_id = c.id
	WHERE 1=1
	AND ($1 = '' OR (LOWER(s.name) LIKE $1))
	AND ($2 = '' OR (LOWER(c.name) LIKE $2))`

	return sqlDB.NewQuery[models.Enrollment](ctx, query, studentName, courseName).Many()
}

func (r *EnrollmentRepository) FindByStudentIdAndCourseId(ctx context.Context, studentId, courseId uuid.UUID) (*models.Enrollment, error) {
	const query = `SELECT 
	  s.id, s.name, s.email, s.birthday, s.created_at,
	  c.id, c.name, c.value, c.created_at,
	  e.installments, e.status, e.created_at
	FROM enrollments e
	JOIN students s ON e.student_id = s.id
	JOIN courses c ON e.course_id = c.id
	WHERE e.student_id = $1
	AND e.course_id = $2`

	return sqlDB.NewQuery[models.Enrollment](ctx, query, studentId, courseId).One()
}

func (r *EnrollmentRepository) Insert(ctx context.Context, model *models.EnrollmentCreateDTO) error {
	const query = `INSERT INTO enrollments(student_id, course_id, installments, status) VALUES($1, $2, $3, $4)`

	return sqlDB.NewStatement(ctx, query, model.Student.ID, model.Course.ID, model.Installments, models.ADIMPLENTE).Execute()
}

func (r *EnrollmentRepository) Delete(ctx context.Context, student_id, course_id uuid.UUID) error {
	const query = `DELETE FROM enrollments WHERE student_id = $1 AND course_id = $2`

	return sqlDB.NewStatement(ctx, query, student_id, course_id).Execute()
}

func (r *EnrollmentRepository) UpdateStatus(ctx context.Context, student_id, course_id uuid.UUID, status models.EnrollmentStatus) error {
	const query = `UPDATE enrollments SET status = $3 WHERE student_id = $1 AND course_id = $2`

	return sqlDB.NewStatement(ctx, query, student_id, course_id, status).Execute()
}
