//go:generate mockgen -source enrollment_repository.go -destination mock/enrollment_repository_mock.go -package mock
package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/types"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/google/uuid"
)

type EnrollmentRepository interface {
	FindPage(ctx context.Context, page *types.PageRequest, filters *models.EnrollmentFilters) (models.EnrollmentPage, error)
	FindByStudentIdAndCourseId(ctx context.Context, studentId, courseId uuid.UUID) (*models.Enrollment, error)
	Insert(ctx context.Context, model *models.EnrollmentCreateDTO) error
	Delete(ctx context.Context, studentID, courseID uuid.UUID) error
	UpdateStatus(ctx context.Context, studentID, courseID uuid.UUID, status models.EnrollmentStatus) error
}

type EnrollmentDBRepository struct{}

func NewEnrollmentDBRepository() *EnrollmentDBRepository {
	return &EnrollmentDBRepository{}
}

func (r *EnrollmentDBRepository) FindPage(ctx context.Context, page *types.PageRequest, filters *models.EnrollmentFilters) (models.EnrollmentPage, error) {
	if filters.StudentName != "" {
		filters.StudentName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.StudentName))
	}

	if filters.CourseName != "" {
		filters.CourseName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.CourseName))
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

	return sqlDB.NewPageQuery[models.Enrollment](ctx, page, query, filters.StudentName, filters.CourseName).Execute()
}

func (r *EnrollmentDBRepository) FindByStudentIdAndCourseId(ctx context.Context, studentID, courseID uuid.UUID) (*models.Enrollment, error) {
	const query = `SELECT 
	  s.id, s.name, s.email, s.birthday, s.created_at,
	  c.id, c.name, c.value, c.created_at,
	  e.installments, e.status, e.created_at
	FROM enrollments e
	JOIN students s ON e.student_id = s.id
	JOIN courses c ON e.course_id = c.id
	WHERE e.student_id = $1
	AND e.course_id = $2`

	return sqlDB.NewQuery[models.Enrollment](ctx, query, studentID, courseID).One()
}

func (r *EnrollmentDBRepository) Insert(ctx context.Context, model *models.EnrollmentCreateDTO) error {
	const query = `INSERT INTO enrollments(student_id, course_id, installments, status) VALUES($1, $2, $3, $4)`

	return sqlDB.NewStatement(ctx, query, model.StudentID, model.CourseID, model.Installments, models.ADIMPLENTE).Execute()
}

func (r *EnrollmentDBRepository) Delete(ctx context.Context, studentID, courseID uuid.UUID) error {
	const query = `DELETE FROM enrollments WHERE student_id = $1 AND course_id = $2`

	return sqlDB.NewStatement(ctx, query, studentID, courseID).Execute()
}

func (r *EnrollmentDBRepository) UpdateStatus(ctx context.Context, studentID, courseID uuid.UUID, status models.EnrollmentStatus) error {
	const query = `UPDATE enrollments SET status = $3 WHERE student_id = $1 AND course_id = $2`

	return sqlDB.NewStatement(ctx, query, studentID, courseID, status).Execute()
}
