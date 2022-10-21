package repositories

import (
	"context"
	"fmt"
	"mime/multipart"
	"school-module/src/domain/models"
	"strings"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/storage"
	"github.com/google/uuid"
)

type IStudentRepository interface {
	FindAll(ctx context.Context, params *models.StudentParams) ([]models.Student, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Student, error)
	Insert(ctx context.Context, model *models.StudentCreateUpdateDTO) error
	Update(ctx context.Context, id uuid.UUID, model *models.StudentCreateUpdateDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
	UploadDocument(ctx context.Context, id uuid.UUID, file *multipart.File) (string, error)
}

type StudentRepository struct{}

func NewStudentRepository() *StudentRepository {
	return &StudentRepository{}
}

func (r *StudentRepository) FindAll(ctx context.Context, params *models.StudentParams) ([]models.Student, error) {
	if params.Name != "" {
		params.Name = fmt.Sprintf("%%%s%%", strings.ToLower(params.Name))
	}

	const query = `SELECT s.id, s.name, s.email, s.birthday, s.created_at FROM students s WHERE ($1 = '' OR (LOWER(s.name) LIKE $1))`

	return sqlDB.NewQuery[models.Student](ctx, query, params.Name).Many()
}

func (r *StudentRepository) FindById(ctx context.Context, id uuid.UUID) (*models.Student, error) {
	return sqlDB.NewQuery[models.Student](ctx, `SELECT s.id, s.name, s.email, s.birthday, s.created_at FROM students s WHERE s.id = $1`, id).One()
}

func (r *StudentRepository) Insert(ctx context.Context, model *models.StudentCreateUpdateDTO) error {
	const query = `INSERT INTO students (name, email, birthday) VALUES($1, $2, $3)`

	return sqlDB.NewStatement(ctx, query, model.Name, model.Email, model.Birthday).Execute()
}

func (r *StudentRepository) Update(ctx context.Context, id uuid.UUID, model *models.StudentCreateUpdateDTO) error {
	return sqlDB.NewStatement(ctx, `UPDATE students SET name=$1, email=$2, birthday=$3 WHERE id=$4`, model.Name, model.Email, model.Birthday, id).Execute()
}

func (r *StudentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return sqlDB.NewStatement(ctx, `DELETE FROM students WHERE id=$1`, id).Execute()
}

func (r *StudentRepository) UploadDocument(ctx context.Context, id uuid.UUID, file *multipart.File) (string, error) {
	return storage.UploadFile(ctx, "meu-bucket", fmt.Sprintf("STUDENT-DOCUMENT-%s", id.String()), file)
}
