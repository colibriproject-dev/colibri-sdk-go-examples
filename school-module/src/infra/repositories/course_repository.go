package repositories

import (
	"context"
	"school-module/src/domain/models"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/cacheDB"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/google/uuid"
)

type ICourseRepository interface {
	FindAll(ctx context.Context) ([]models.Course, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Course, error)
	Insert(ctx context.Context, model *models.CourseCreateUpdateDTO) (*models.Course, error)
	Update(ctx context.Context, id uuid.UUID, model *models.CourseCreateUpdateDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CourseRepository struct{}

func NewCourseRepository() *CourseRepository {
	return &CourseRepository{}
}

var coursesCache = cacheDB.NewCache[models.Course]("Courses", 10*time.Second)

func (r *CourseRepository) FindAll(ctx context.Context) ([]models.Course, error) {
	const query = "SELECT id, name, value, created_at FROM courses"
	return sqlDB.NewCachedQuery(ctx, coursesCache, query).Many()
}

func (r *CourseRepository) FindById(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	const query = "SELECT id, name, value, created_at FROM courses WHERE id = $1"
	return sqlDB.NewQuery[models.Course](ctx, query, id).One()
}

func (r *CourseRepository) Insert(ctx context.Context, model *models.CourseCreateUpdateDTO) (*models.Course, error) {
	const query = "INSERT INTO courses (name, value) VALUES ($1, $2) RETURNING id, name, value, created_at"
	coursesCache.Del(ctx)
	return sqlDB.NewQuery[models.Course](ctx, query, model.Name, model.Value).One()
}

func (r *CourseRepository) Update(ctx context.Context, id uuid.UUID, model *models.CourseCreateUpdateDTO) error {
	coursesCache.Del(ctx)
	return sqlDB.NewStatement(ctx, "UPDATE courses SET name=$1, value=$2 WHERE id=$3", model.Name, model.Value, id).Execute()
}

func (r *CourseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	coursesCache.Del(ctx)
	return sqlDB.NewStatement(ctx, "DELETE FROM courses WHERE id=$1", id).Execute()
}
