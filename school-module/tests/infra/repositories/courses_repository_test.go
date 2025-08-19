package repositories

import (
	"testing"
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	courseDatasets   = []string{"clear-data.sql", "courses-insert.sql"}
	courseRepository = repositories.NewCoursesDBRepository()
	courseMockData   = []models.Course{
		{ID: uuid.MustParse("64bb5356-e900-4929-8d4d-debe31da40cc"), Name: "COURSE TEST 1", Value: 1000.50, CreatedAt: time.Date(2023, time.July, 24, 10, 0, 0, 0, time.FixedZone("", 0))},
		{ID: uuid.MustParse("0319ea44-b0d6-4726-85a0-bbacf16dbe46"), Name: "COURSE TEST 2", Value: 2500.90, CreatedAt: time.Date(2023, time.July, 25, 16, 0, 0, 0, time.FixedZone("", 0))},
		{ID: uuid.MustParse("f600b3a4-fc3e-45e2-9131-1730a58ad9b4"), Name: "COURSE TEST 3", Value: 1800.22, CreatedAt: time.Date(2023, time.July, 26, 18, 0, 0, 0, time.FixedZone("", 0))},
	}
)

func TestCourseRepository_FindAll(t *testing.T) {
	assert.NoError(t, pc.Dataset(basePath, courseDatasets...))

	t.Run("Should return all courses when exists into courses table", func(t *testing.T) {
		expected := courseMockData

		result, err := courseRepository.FindAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
	})
}

func TestCourseRepository_FindById(t *testing.T) {
	assert.NoError(t, pc.Dataset(basePath, courseDatasets...))

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
		assert.EqualValues(t, expected, result)
	})
}

func TestCourseRepository_Insert(t *testing.T) {
	assert.NoError(t, pc.Dataset(basePath, courseDatasets...))

	t.Run("Should return error when name exists into courses table", func(t *testing.T) {
		result, err := courseRepository.Insert(ctx, &models.CourseCreate{Name: courseMockData[0].Name, Value: courseMockData[0].Value})

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should insert and return course when not exists into courses table", func(t *testing.T) {
		model := &models.CourseCreate{Name: "COURSE TEST 999", Value: 500.00}

		result, err := courseRepository.Insert(ctx, model)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
		assert.NotEmpty(t, result.CreatedAt)
		assert.EqualValues(t, model.Name, result.Name)
		assert.EqualValues(t, model.Value, result.Value)
	})
}

func TestCourseRepository_Update(t *testing.T) {
	assert.NoError(t, pc.Dataset(basePath, courseDatasets...))

	t.Run("Should return error when name exists into courses table", func(t *testing.T) {
		err := courseRepository.Update(ctx, &models.CourseUpdate{
			Name:  courseMockData[1].Name,
			Value: courseMockData[1].Value,
			ID:    courseMockData[0].ID,
		})

		assert.Error(t, err)
	})

	t.Run("Should update course when exists into courses table", func(t *testing.T) {
		model := &models.CourseUpdate{Name: "COURSE TEST 999", Value: 500.00, ID: courseMockData[0].ID}

		errUpdate := courseRepository.Update(ctx, model)
		result, err := courseRepository.FindById(ctx, courseMockData[0].ID)

		assert.NoError(t, errUpdate)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.CreatedAt)
		assert.EqualValues(t, courseMockData[0].ID, result.ID)
		assert.EqualValues(t, model.Name, result.Name)
		assert.EqualValues(t, model.Value, result.Value)
	})
}

func TestCourseRepository_Delete(t *testing.T) {
	assert.NoError(t, pc.Dataset(basePath, courseDatasets...))

	t.Run("Should delete course when exists into courses table", func(t *testing.T) {
		errDelete := courseRepository.Delete(ctx, courseMockData[0].ID)
		result, err := courseRepository.FindById(ctx, courseMockData[0].ID)

		assert.NoError(t, errDelete)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}
