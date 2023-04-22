//go:generate mockgen -source course_repository.go -destination mock/course_repository_mock.go -package mock
package repositories

import (
	"context"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/cacheDB"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/google/uuid"
)

type CourseRepository interface {
	FindAll(ctx context.Context) ([]models.Course, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Course, error)
	Insert(ctx context.Context, model *models.CourseCreateUpdateDTO) (*models.Course, error)
	Update(ctx context.Context, id uuid.UUID, model *models.CourseCreateUpdateDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CourseDBRepository struct{}

func NewCourseDBRepository() *CourseDBRepository {
	return &CourseDBRepository{}
}

var coursesCache = cacheDB.NewCache[models.Course]("Courses", 10*time.Second)

func (r *CourseDBRepository) FindAll(ctx context.Context) ([]models.Course, error) {
	const query = "SELECT id, name, value, created_at FROM courses"
	return sqlDB.NewCachedQuery(ctx, coursesCache, query).Many()
}

func (r *CourseDBRepository) FindById(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	const query = "SELECT id, name, value, created_at FROM courses WHERE id = $1"
	return sqlDB.NewQuery[models.Course](ctx, query, id).One()
}

func (r *CourseDBRepository) Insert(ctx context.Context, model *models.CourseCreateUpdateDTO) (*models.Course, error) {
	const query = "INSERT INTO courses (name, value) VALUES ($1, $2) RETURNING id, name, value, created_at"
	coursesCache.Del(ctx)
	return sqlDB.NewQuery[models.Course](ctx, query, model.Name, model.Value).One()
}

func (r *CourseDBRepository) Update(ctx context.Context, id uuid.UUID, model *models.CourseCreateUpdateDTO) error {
	coursesCache.Del(ctx)
	return sqlDB.NewStatement(ctx, "UPDATE courses SET name=$1, value=$2 WHERE id=$3", model.Name, model.Value, id).Execute()
}

func (r *CourseDBRepository) Delete(ctx context.Context, id uuid.UUID) error {
	coursesCache.Del(ctx)
	return sqlDB.NewStatement(ctx, "DELETE FROM courses WHERE id=$1", id).Execute()
}
