//go:generate mockgen -source students_repository.go -destination mock/students_repository_mock.go -package repositoriesmock
package repositories

import (
	"context"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/types"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/google/uuid"
)

const (
	studentBaseQuery          = `SELECT id, name, email, birthday, created_at FROM students`
	findAllStudentsQuery      = studentBaseQuery + ` WHERE ($1 = '' OR (name ILIKE CONCAT('%', $1, '%')))`
	findStudentByIdQuery      = studentBaseQuery + ` WHERE id = $1`
	findStudentByEmailQuery   = studentBaseQuery + ` WHERE email = $1`
	existsStudentByIdQuery    = `SELECT EXISTS(SELECT 1 FROM students WHERE id = $1)`
	existsStudentByEmailQuery = `SELECT EXISTS(SELECT 1 FROM students WHERE email = $1)`
	insertStudentQuery        = `INSERT INTO students (name, email, birthday) VALUES ($1, $2, $3)`
	updateStudentQuery        = `UPDATE students SET name=$1, email=$2, birthday=$3 WHERE id=$4`
	deleteStudentQuery        = `DELETE FROM students WHERE id=$1`
)

type IStudentsRepository interface {
	FindAllPaginated(ctx context.Context, params *models.StudentPageParams) (models.StudentPage, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Student, error)
	FindByEmail(ctx context.Context, email string) (*models.Student, error)
	ExistsById(ctx context.Context, id uuid.UUID) (*bool, error)
	ExistsByEmail(ctx context.Context, email string) (*bool, error)
	Insert(ctx context.Context, model *models.StudentCreate) error
	Update(ctx context.Context, model *models.StudentUpdate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type StudentsDBRepository struct{}

func NewStudentsDBRepository() *StudentsDBRepository {
	return &StudentsDBRepository{}
}

func (r *StudentsDBRepository) FindAllPaginated(ctx context.Context, params *models.StudentPageParams) (models.StudentPage, error) {
	return sqlDB.NewPageQuery[models.Student](ctx,
		types.NewPageRequest(params.Page, params.Size, []types.Sort{types.NewSort(types.ASC, "name")}),
		findAllStudentsQuery,
		params.Name,
	).Execute()
}

func (r *StudentsDBRepository) FindById(ctx context.Context, id uuid.UUID) (*models.Student, error) {
	return sqlDB.NewQuery[models.Student](ctx, findStudentByIdQuery, id).One()
}

func (r *StudentsDBRepository) FindByEmail(ctx context.Context, email string) (*models.Student, error) {
	return sqlDB.NewQuery[models.Student](ctx, findStudentByEmailQuery, email).One()
}

func (r *StudentsDBRepository) ExistsById(ctx context.Context, id uuid.UUID) (*bool, error) {
	return sqlDB.NewQuery[bool](ctx, existsStudentByIdQuery, id).One()
}

func (r *StudentsDBRepository) ExistsByEmail(ctx context.Context, email string) (*bool, error) {
	return sqlDB.NewQuery[bool](ctx, existsStudentByEmailQuery, email).One()
}

func (r *StudentsDBRepository) Insert(ctx context.Context, model *models.StudentCreate) error {
	return sqlDB.NewStatement(ctx, insertStudentQuery, model.Name, model.Email, model.Birthday).Execute()
}

func (r *StudentsDBRepository) Update(ctx context.Context, model *models.StudentUpdate) error {
	return sqlDB.NewStatement(ctx, updateStudentQuery, model.Name, model.Email, model.Birthday, model.ID).Execute()
}

func (r *StudentsDBRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return sqlDB.NewStatement(ctx, deleteStudentQuery, id).Execute()
}
