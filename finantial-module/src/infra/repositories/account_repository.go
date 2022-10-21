package repositories

import (
	"context"
	"finantial-module/src/domain/models"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/google/uuid"
)

type AccountRepository interface {
	FindAll(ctx context.Context) ([]models.Account, error)
	Insert(ctx context.Context, model *models.Account) error
	UpdateStatus(ctx context.Context, account *models.Account) error
	DeleteByStudentAndCourse(ctx context.Context, studentId, courseId uuid.UUID) error
	DeleteByCourse(ctx context.Context, courseId uuid.UUID) error
	DeleteByStudent(ctx context.Context, studentId uuid.UUID) error
}

type AccountDBRepository struct{}

func NewAccountDBRepository() *AccountDBRepository {
	return &AccountDBRepository{}
}

func (r *AccountDBRepository) FindAll(ctx context.Context) ([]models.Account, error) {
	const query = `SELECT a.id, a.student_id, a.course_id, a.installments, a.value, a.status, a.created_at FROM accounts a`

	return sqlDB.NewQuery[models.Account](ctx, query).Many()
}

func (r *AccountDBRepository) Insert(ctx context.Context, model *models.Account) error {
	const query = `INSERT INTO accounts (id, student_id, course_id, installments, value, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	return sqlDB.NewStatement(ctx, query,
		model.ID, model.StudentID, model.CourseID, model.Installments, model.Value, model.Status, model.CreatedAt,
	).Execute()
}

func (r *AccountDBRepository) UpdateStatus(ctx context.Context, account *models.Account) error {
	const query = `UPDATE accounts SET status = $3 WHERE student_id = $1 AND course_id = $2`

	return sqlDB.NewStatement(ctx, query, account.StudentID, account.CourseID, account.Status).Execute()
}

func (r *AccountDBRepository) DeleteByStudentAndCourse(ctx context.Context, studentId, courseId uuid.UUID) error {
	const query = `DELETE FROM accounts WHERE student_id = $1 AND course_id = $2`

	return sqlDB.NewStatement(ctx, query, studentId, courseId).Execute()
}

func (r *AccountDBRepository) DeleteByCourse(ctx context.Context, courseId uuid.UUID) error {
	const query = `DELETE FROM accounts WHERE course_id = $1`

	return sqlDB.NewStatement(ctx, query, courseId).Execute()
}

func (r *AccountDBRepository) DeleteByStudent(ctx context.Context, studentId uuid.UUID) error {
	const query = `DELETE FROM accounts WHERE student_id = $1`

	return sqlDB.NewStatement(ctx, query, studentId).Execute()
}
