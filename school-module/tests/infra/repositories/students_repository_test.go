package repositories

import (
	"testing"
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	studentDatasets   = []string{"clear-data.sql", "students-insert.sql"}
	studentRepository = repositories.NewStudentsDBRepository()
	studentMockData   = []models.Student{
		{ID: uuid.MustParse("53bb5356-e900-4929-8d4d-debe31da40bb"), Name: "STUDENT TEST 1", Email: "student1@email.com", Birthday: types.IsoDate(time.Date(2000, time.August, 24, 0, 0, 0, 0, time.FixedZone("", 0))), CreatedAt: time.Date(2023, time.August, 24, 9, 0, 0, 0, time.FixedZone("", 0))},
		{ID: uuid.MustParse("9219ea44-b0d6-4726-85a0-bbacf16dbe35"), Name: "STUDENT TEST 2", Email: "student2@email.com", Birthday: types.IsoDate(time.Date(2001, time.August, 25, 0, 0, 0, 0, time.FixedZone("", 0))), CreatedAt: time.Date(2023, time.August, 25, 15, 0, 0, 0, time.FixedZone("", 0))},
		{ID: uuid.MustParse("f590b3a4-fc3e-45e2-9131-1730a58ad9a3"), Name: "STUDENT TEST 3", Email: "student3@email.com", Birthday: types.IsoDate(time.Date(2002, time.August, 26, 0, 0, 0, 0, time.FixedZone("", 0))), CreatedAt: time.Date(2023, time.August, 26, 17, 0, 0, 0, time.FixedZone("", 0))},
	}
)

func TestStudentRepository_FindAll(t *testing.T) {
	assert.NoError(t, pc.Dataset(basePath, studentDatasets...))

	t.Run("Should return all students when exists into students table", func(t *testing.T) {
		expected := studentMockData

		result, err := studentRepository.FindAllPaginated(ctx, &models.StudentPageParams{Page: 1, Size: 10})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result.Items)
		assert.EqualValues(t, uint64(3), result.TotalItems)
	})

	t.Run("Should return filtered students when exists into students table", func(t *testing.T) {
		expected := []models.Student{studentMockData[0]}

		result, err := studentRepository.FindAllPaginated(ctx, &models.StudentPageParams{Page: 1, Size: 10, Name: "1"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result.Items)
		assert.EqualValues(t, uint64(1), result.TotalItems)
	})
}

func TestStudentRepository_FindById(t *testing.T) {
	assert.NoError(t, pc.Dataset(basePath, studentDatasets...))

	t.Run("Should return nil when id not exists into students table", func(t *testing.T) {
		result, err := studentRepository.FindById(ctx, uuid.New())

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return student when exists into students table", func(t *testing.T) {
		expected := &studentMockData[0]

		result, err := studentRepository.FindById(ctx, expected.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
	})
}

func TestStudentRepository_Insert(t *testing.T) {
	assert.NoError(t, pc.Dataset(basePath, studentDatasets...))

	t.Run("Should return error when name exists into students table", func(t *testing.T) {
		err := studentRepository.Insert(ctx, &models.StudentCreate{
			Name:     studentMockData[0].Name,
			Email:    studentMockData[0].Email,
			Birthday: studentMockData[0].Birthday,
		})

		assert.Error(t, err)
	})

	t.Run("Should insert and return student when not exists into students table", func(t *testing.T) {
		model := &models.StudentCreate{
			Name:     "COURSE TEST 999",
			Email:    "student999@email.com",
			Birthday: types.IsoDate(time.Date(2000, time.August, 10, 0, 0, 0, 0, time.FixedZone("", 0))),
		}

		errInsert := studentRepository.Insert(ctx, model)
		result, resultErr := studentRepository.FindByEmail(ctx, model.Email)

		assert.NoError(t, errInsert)
		assert.NoError(t, resultErr)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
		assert.NotEmpty(t, result.CreatedAt)
		assert.EqualValues(t, model.Name, result.Name)
		assert.EqualValues(t, model.Email, result.Email)
		assert.EqualValues(t, model.Birthday, result.Birthday)
	})
}

func TestStudentRepository_Update(t *testing.T) {
	assert.NoError(t, pc.Dataset(basePath, studentDatasets...))

	t.Run("Should return error when name exists into students table", func(t *testing.T) {
		err := studentRepository.Update(ctx, &models.StudentUpdate{
			Name:     studentMockData[1].Name,
			Email:    studentMockData[1].Email,
			Birthday: studentMockData[1].Birthday,
			ID:       studentMockData[0].ID,
		})

		assert.Error(t, err)
	})

	t.Run("Should update student when exists into students table", func(t *testing.T) {
		model := &models.StudentUpdate{
			Name:     "STUDENT TEST 999",
			Email:    "student1@email.com",
			Birthday: types.IsoDate(time.Date(2000, time.August, 24, 0, 0, 0, 0, time.FixedZone("", 0))),
			ID:       studentMockData[0].ID,
		}

		errUpdate := studentRepository.Update(ctx, model)
		result, err := studentRepository.FindById(ctx, studentMockData[0].ID)

		assert.NoError(t, errUpdate)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.CreatedAt)
		assert.EqualValues(t, studentMockData[0].ID, result.ID)
		assert.EqualValues(t, model.Name, result.Name)
		assert.EqualValues(t, model.Email, result.Email)
		assert.EqualValues(t, model.Birthday, result.Birthday)
	})
}

func TestStudentRepository_Delete(t *testing.T) {
	assert.NoError(t, pc.Dataset(basePath, studentDatasets...))

	t.Run("Should delete student when exists into students table", func(t *testing.T) {
		errDelete := studentRepository.Delete(ctx, studentMockData[0].ID)
		result, err := studentRepository.FindById(ctx, studentMockData[0].ID)

		assert.NoError(t, errDelete)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}
