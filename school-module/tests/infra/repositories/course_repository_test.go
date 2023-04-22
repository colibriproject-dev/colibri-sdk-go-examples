package repositories

import (
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	courseRepository = repositories.NewCourseDBRepository()
	courseMockData   = []models.Course{
		{ID: uuid.MustParse("64bb5356-e900-4929-8d4d-debe31da40cc"), Name: "COURSE TEST 1", Value: 1000.50, CreatedAt: time.Date(2023, time.July, 24, 10, 0, 0, 0, time.FixedZone("", 0))},
		{ID: uuid.MustParse("0319ea44-b0d6-4726-85a0-bbacf16dbe46"), Name: "COURSE TEST 2", Value: 2500.90, CreatedAt: time.Date(2023, time.July, 25, 16, 0, 0, 0, time.FixedZone("", 0))},
		{ID: uuid.MustParse("f600b3a4-fc3e-45e2-9131-1730a58ad9b4"), Name: "COURSE TEST 3", Value: 1800.22, CreatedAt: time.Date(2023, time.July, 26, 18, 0, 0, 0, time.FixedZone("", 0))},
	}
)

func TestCourseRepositoryFindAllCourses(t *testing.T) {
	datasets := []string{"clear-data.sql", "course-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

	t.Run("Should return all courses when exists into courses table", func(t *testing.T) {
		expected := courseMockData

		result, err := courseRepository.FindAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected, result)
	})
}

func TestCourseRepositoryFindById(t *testing.T) {
	datasets := []string{"clear-data.sql", "course-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

	t.Run("Should return nil when id not exists into courses table", func(t *testing.T) {
		result, err := courseRepository.FindById(ctx, uuid.New())

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return course when exists into courses table", func(t *testing.T) {
		expected := &courseMockData[0]

		result, err := courseRepository.FindById(ctx, expected.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected, result)
	})
}

func TestCourseRepositoryInsert(t *testing.T) {
	datasets := []string{"clear-data.sql", "course-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

	t.Run("Should return error when name exists into courses table", func(t *testing.T) {
		result, err := courseRepository.Insert(ctx, &models.CourseCreateUpdateDTO{Name: courseMockData[0].Name, Value: courseMockData[0].Value})

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should insert and return course when not exists into courses table", func(t *testing.T) {
		model := &models.CourseCreateUpdateDTO{Name: "COURSE TEST 999", Value: 500.00}

		result, err := courseRepository.Insert(ctx, model)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
		assert.NotEmpty(t, result.CreatedAt)
		assert.Equal(t, model.Name, result.Name)
		assert.Equal(t, model.Value, result.Value)
	})
}

func TestCourseRepositoryUpdate(t *testing.T) {
	datasets := []string{"clear-data.sql", "course-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

	t.Run("Should return error when name exists into courses table", func(t *testing.T) {
		err := courseRepository.Update(ctx, courseMockData[0].ID, &models.CourseCreateUpdateDTO{Name: courseMockData[1].Name, Value: courseMockData[1].Value})

		assert.Error(t, err)
	})

	t.Run("Should update course when exists into courses table", func(t *testing.T) {
		model := &models.CourseCreateUpdateDTO{Name: "COURSE TEST 999", Value: 500.00}

		errUpdate := courseRepository.Update(ctx, courseMockData[0].ID, model)
		result, err := courseRepository.FindById(ctx, courseMockData[0].ID)

		assert.NoError(t, errUpdate)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.CreatedAt)
		assert.Equal(t, courseMockData[0].ID, result.ID)
		assert.Equal(t, model.Name, result.Name)
		assert.Equal(t, model.Value, result.Value)
	})
}

func TestCourseRepositoryDelete(t *testing.T) {
	datasets := []string{"clear-data.sql", "course-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

	t.Run("Should delete course when exists into courses table", func(t *testing.T) {
		errDelete := courseRepository.Delete(ctx, courseMockData[0].ID)
		result, err := courseRepository.FindById(ctx, courseMockData[0].ID)

		assert.NoError(t, errDelete)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}
