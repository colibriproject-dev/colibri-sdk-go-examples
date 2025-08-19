//go:generate mockgen -source courses_repository.go -destination mock/courses_repository_mock.go -package repositoriesmock
package repositories

import (
	"context"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/google/uuid"
)

const (
	courseBaseQuery         = `SELECT id, name, value, created_at FROM courses`
	findAllCoursesQuery     = courseBaseQuery
	findCourseByIdQuery     = courseBaseQuery + ` WHERE id = $1`
	findCourseByNameQuery   = courseBaseQuery + ` WHERE name = $1`
	existsCourseByIdQuery   = `SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)`
	existsCourseByNameQuery = `SELECT EXISTS(SELECT 1 FROM courses WHERE name = $1)`
	insertCourseQuery       = `INSERT INTO courses (name, value) VALUES ($1, $2) RETURNING id, name, value, created_at`
	updateCourseQuery       = `UPDATE courses SET name=$1, value=$2 WHERE id=$3`
	deleteCourseQuery       = `DELETE FROM courses WHERE id=$1`
)

type ICoursesRepository interface {
	FindAll(ctx context.Context) ([]models.Course, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Course, error)
	FindByName(ctx context.Context, name string) (*models.Course, error)
	ExistsById(ctx context.Context, id uuid.UUID) (*bool, error)
	ExistsByName(ctx context.Context, name string) (*bool, error)
	Insert(ctx context.Context, model *models.CourseCreate) (*models.Course, error)
	Update(ctx context.Context, model *models.CourseUpdate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CoursesDBRepository struct{}

func NewCoursesDBRepository() *CoursesDBRepository {
	return &CoursesDBRepository{}
}

func (r *CoursesDBRepository) FindAll(ctx context.Context) ([]models.Course, error) {
	return sqlDB.NewQuery[models.Course](ctx, findAllCoursesQuery).Many()
}

func (r *CoursesDBRepository) FindById(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	return sqlDB.NewQuery[models.Course](ctx, findCourseByIdQuery, id).One()
}

func (r *CoursesDBRepository) FindByName(ctx context.Context, name string) (*models.Course, error) {
	return sqlDB.NewQuery[models.Course](ctx, findCourseByNameQuery, name).One()
}

func (r *CoursesDBRepository) ExistsById(ctx context.Context, id uuid.UUID) (*bool, error) {
	return sqlDB.NewQuery[bool](ctx, existsCourseByIdQuery, id).One()
}

func (r *CoursesDBRepository) ExistsByName(ctx context.Context, name string) (*bool, error) {
	return sqlDB.NewQuery[bool](ctx, existsCourseByNameQuery, name).One()
}

func (r *CoursesDBRepository) Insert(ctx context.Context, model *models.CourseCreate) (*models.Course, error) {
	return sqlDB.NewQuery[models.Course](ctx, insertCourseQuery, model.Name, model.Value).One()
}

func (r *CoursesDBRepository) Update(ctx context.Context, model *models.CourseUpdate) error {
	return sqlDB.NewStatement(ctx, updateCourseQuery, model.Name, model.Value, model.ID).Execute()
}

func (r *CoursesDBRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return sqlDB.NewStatement(ctx, deleteCourseQuery, id).Execute()
}
