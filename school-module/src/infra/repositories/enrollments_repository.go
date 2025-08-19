//go:generate mockgen -source enrollments_repository.go -destination mock/enrollments_repository_mock.go -package repositoriesmock
package repositories

import (
	"context"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/types"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/google/uuid"
)

const (
	enrollmentBaseQuery = `
		SELECT 
			s.id, s.name, s.email, s.birthday, s.created_at,
			c.id, c.name, c.value, c.created_at,
			e.installments, e.status, e.created_at
		FROM enrollments e
		JOIN students s ON e.student_id = s.id
		JOIN courses c ON e.course_id = c.id`

	findAllPaginatedEnrollmentQuery = enrollmentBaseQuery + `
		WHERE 1=1
 		AND ($1 = '' OR (s.name ILIKE CONCAT('%', $1, '%')))
 		AND ($2 = '' OR (c.name ILIKE CONCAT('%', $2, '%')))`

	findEnrollmentByStudentIdAndCourseIdQuery = enrollmentBaseQuery + `
		WHERE e.student_id = $1
		AND e.course_id = $2`

	existsEnrollmentByStudentIdAndCourseIdQuery = `
		SELECT EXISTS (
			SELECT 1 FROM enrollments e
			WHERE e.student_id = $1 AND e.course_id = $2
		)`

	insertEnrollmentQuery = `
		INSERT INTO enrollments(student_id, course_id, installments, status)
		VALUES ($1, $2, $3, $4)
		RETURNING student_id, course_id, installments, status, created_at`

	deleteEnrollmentQuery       = `DELETE FROM enrollments WHERE student_id = $1 AND course_id = $2`
	updateEnrollmentStatusQuery = `UPDATE enrollments SET status = $3 WHERE student_id = $1 AND course_id = $2`
)

type IEnrollmentsRepository interface {
	FindAllPaginated(ctx context.Context, params *models.EnrollmentPageParams) (models.EnrollmentPage, error)
	FindByStudentIdAndCourseId(ctx context.Context, studentID, courseID uuid.UUID) (*models.Enrollment, error)
	ExistsByStudentIdAndCourseId(ctx context.Context, studentID, courseID uuid.UUID) (*bool, error)
	Insert(ctx context.Context, model *models.EnrollmentCreate) (*models.EnrollmentCreated, error)
	Delete(ctx context.Context, studentID, courseID uuid.UUID) error
	UpdateStatus(ctx context.Context, model *models.EnrollmentUpdateStatus) error
}

type EnrollmentsDBRepository struct{}

func NewEnrollmentsDBRepository() *EnrollmentsDBRepository {
	return &EnrollmentsDBRepository{}
}

func (r *EnrollmentsDBRepository) FindAllPaginated(ctx context.Context, params *models.EnrollmentPageParams) (models.EnrollmentPage, error) {
	return sqlDB.NewPageQuery[models.Enrollment](
		ctx,
		types.NewPageRequest(params.Page, params.Size, []types.Sort{types.NewSort(types.DESC, "e.created_at")}),
		findAllPaginatedEnrollmentQuery,
		params.StudentName,
		params.CourseName,
	).Execute()
}

func (r *EnrollmentsDBRepository) FindByStudentIdAndCourseId(ctx context.Context, studentID, courseID uuid.UUID) (*models.Enrollment, error) {
	return sqlDB.NewQuery[models.Enrollment](ctx, findEnrollmentByStudentIdAndCourseIdQuery, studentID, courseID).One()
}

func (r *EnrollmentsDBRepository) ExistsByStudentIdAndCourseId(ctx context.Context, studentID, courseID uuid.UUID) (*bool, error) {
	return sqlDB.NewQuery[bool](ctx, existsEnrollmentByStudentIdAndCourseIdQuery, studentID, courseID).One()
}

func (r *EnrollmentsDBRepository) Insert(ctx context.Context, model *models.EnrollmentCreate) (*models.EnrollmentCreated, error) {
	return sqlDB.NewQuery[models.EnrollmentCreated](ctx,
		insertEnrollmentQuery,
		model.StudentID,
		model.CourseID,
		model.Installments,
		model.Status,
	).One()
}

func (r *EnrollmentsDBRepository) Delete(ctx context.Context, studentID, courseID uuid.UUID) error {
	return sqlDB.NewStatement(ctx, deleteEnrollmentQuery, studentID, courseID).Execute()
}

func (r *EnrollmentsDBRepository) UpdateStatus(ctx context.Context, model *models.EnrollmentUpdateStatus) error {
	return sqlDB.NewStatement(ctx,
		updateEnrollmentStatusQuery,
		model.StudentID,
		model.CourseID,
		model.Status,
	).Execute()
}
