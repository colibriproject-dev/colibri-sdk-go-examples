package repositories

import (
	"mime/multipart"
	"os"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	studentRepository = repositories.NewStudentDBRepository()
	studentMockData   = []models.Student{
		{ID: uuid.MustParse("53bb5356-e900-4929-8d4d-debe31da40bb"), Name: "STUDENT TEST 1", Email: "student1@email.com", Birthday: time.Date(2000, time.August, 24, 0, 0, 0, 0, time.FixedZone("", 0)), CreatedAt: time.Date(2023, time.August, 24, 9, 0, 0, 0, time.FixedZone("", 0))},
		{ID: uuid.MustParse("9219ea44-b0d6-4726-85a0-bbacf16dbe35"), Name: "STUDENT TEST 2", Email: "student2@email.com", Birthday: time.Date(2001, time.August, 25, 0, 0, 0, 0, time.FixedZone("", 0)), CreatedAt: time.Date(2023, time.August, 25, 15, 0, 0, 0, time.FixedZone("", 0))},
		{ID: uuid.MustParse("f590b3a4-fc3e-45e2-9131-1730a58ad9a3"), Name: "STUDENT TEST 3", Email: "student3@email.com", Birthday: time.Date(2002, time.August, 26, 0, 0, 0, 0, time.FixedZone("", 0)), CreatedAt: time.Date(2023, time.August, 26, 17, 0, 0, 0, time.FixedZone("", 0))},
	}
)

func TestStudentRepositoryFindAllStudents(t *testing.T) {
	datasets := []string{"clear-data.sql", "student-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

	t.Run("Should return all students when exists into students table", func(t *testing.T) {
		expected := studentMockData

		result, err := studentRepository.FindAll(ctx, &models.StudentParams{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected, result)
	})

	t.Run("Should return filtered students when exists into students table", func(t *testing.T) {
		expected := []models.Student{studentMockData[0]}

		result, err := studentRepository.FindAll(ctx, &models.StudentParams{Name: "1"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected, result)
	})
}

func TestStudentRepositoryFindById(t *testing.T) {
	datasets := []string{"clear-data.sql", "student-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

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
		assert.Equal(t, expected, result)
	})
}

func TestStudentRepositoryInsert(t *testing.T) {
	datasets := []string{"clear-data.sql", "student-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

	t.Run("Should return error when name exists into students table", func(t *testing.T) {
		err := studentRepository.Insert(ctx, &models.StudentCreateUpdateDTO{
			Name:     studentMockData[0].Name,
			Email:    studentMockData[0].Email,
			Birthday: studentMockData[0].Birthday,
		})

		assert.Error(t, err)
	})

	t.Run("Should insert and return student when not exists into students table", func(t *testing.T) {
		model := &models.StudentCreateUpdateDTO{Name: "COURSE TEST 999", Email: "student999@email.com", Birthday: time.Date(2000, time.August, 10, 0, 0, 0, 0, time.FixedZone("", 0))}

		errInsert := studentRepository.Insert(ctx, model)
		result, err := studentRepository.FindAll(ctx, &models.StudentParams{Name: model.Name})

		assert.NoError(t, errInsert)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.NotEmpty(t, result[0].ID)
		assert.NotEmpty(t, result[0].CreatedAt)
		assert.Equal(t, model.Name, result[0].Name)
		assert.Equal(t, model.Email, result[0].Email)
		assert.Equal(t, model.Birthday, result[0].Birthday)
	})
}

func TestStudentRepositoryUpdate(t *testing.T) {
	datasets := []string{"clear-data.sql", "student-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

	t.Run("Should return error when name exists into students table", func(t *testing.T) {
		err := studentRepository.Update(ctx, studentMockData[0].ID, &models.StudentCreateUpdateDTO{
			Name:     studentMockData[1].Name,
			Email:    studentMockData[1].Email,
			Birthday: studentMockData[1].Birthday,
		})

		assert.Error(t, err)
	})

	t.Run("Should update student when exists into students table", func(t *testing.T) {
		model := &models.StudentCreateUpdateDTO{
			Name:     "STUDENT TEST 999",
			Email:    "student1@email.com",
			Birthday: time.Date(2000, time.August, 24, 0, 0, 0, 0, time.FixedZone("", 0)),
		}

		errUpdate := studentRepository.Update(ctx, studentMockData[0].ID, model)
		result, err := studentRepository.FindById(ctx, studentMockData[0].ID)

		assert.NoError(t, errUpdate)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.CreatedAt)
		assert.Equal(t, studentMockData[0].ID, result.ID)
		assert.Equal(t, model.Name, result.Name)
		assert.Equal(t, model.Email, result.Email)
		assert.Equal(t, model.Birthday, result.Birthday)
	})
}

func TestStudentRepositoryDelete(t *testing.T) {
	datasets := []string{"clear-data.sql", "student-insert.sql"}
	err := pc.Dataset(basePath, datasets...)
	assert.NoError(t, err)

	t.Run("Should delete student when exists into students table", func(t *testing.T) {
		errDelete := studentRepository.Delete(ctx, studentMockData[0].ID)
		result, err := studentRepository.FindById(ctx, studentMockData[0].ID)

		assert.NoError(t, errDelete)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestStudentRepositoryUploadDocument(t *testing.T) {
	t.Run("Should upload document to storage", func(t *testing.T) {
		var file multipart.File
		file, err := os.Open(test.MountAbsolutPath("../../../development-environment/files/img.png"))
		assert.NoError(t, err)

		result, err := studentRepository.UploadDocument(ctx, uuid.New(), &file)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}
