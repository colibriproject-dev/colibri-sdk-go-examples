package controllers

import (
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/application/controllers"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/usecases/mock"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/restserver"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStudentController(t *testing.T) {
	t.Run("Should return new student controller", func(t *testing.T) {
		result := controllers.NewStudentController()

		assert.NotNil(t, result)
		assert.NotNil(t, result.Usecase)
		assert.NotNil(t, result.Routes())
	})
}

func TestGetAllStudents(t *testing.T) {
	controller := gomock.NewController(t)
	usecaseMock := mock.NewMockIStudentUsecases(controller)
	restController := controllers.StudentController{Usecase: usecaseMock}
	defer controller.Finish()

	t.Run("Should return error when occurred error in GetAll", func(t *testing.T) {
		errMock := errors.New("mock GetAll")
		usecaseMock.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(nil, errMock)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   "/students",
			Url:    "/students",
		}, restController.GetAll)

		var result restserver.Error
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, errMock.Error(), result.Error)
	})

	t.Run("Should return all students", func(t *testing.T) {
		expected := []models.Student{
			{ID: uuid.MustParse("9f9fa978-7df0-4474-b1d4-6be55e0dbd1d"), Name: "Student name 1", Email: "student1@email.com", Birthday: time.Time{}, CreatedAt: time.Time{}},
			{ID: uuid.MustParse("6079def8-f9c2-4258-bbf1-118c1ffb0a67"), Name: "Student name 2", Email: "student2@email.com", Birthday: time.Time{}, CreatedAt: time.Time{}},
			{ID: uuid.MustParse("2c61d66e-1e6b-48fa-921c-eb96d009f8e8"), Name: "Student name 3", Email: "student3@email.com", Birthday: time.Time{}, CreatedAt: time.Time{}},
		}
		usecaseMock.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(expected, nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   "/students",
			Url:    "/students",
		}, restController.GetAll)

		var result []models.Student
		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, expected, result)
	})
}

func TestGetStudentById(t *testing.T) {
	controller := gomock.NewController(t)
	usecaseMock := mock.NewMockIStudentUsecases(controller)
	restController := controllers.StudentController{Usecase: usecaseMock}
	defer controller.Finish()

	t.Run("Should return error when i try get invalid id path param", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   "/students/{id}",
			Url:    "/students/abc",
		}, restController.GetById)

		var result restserver.Error
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return error when occurred error in GetById", func(t *testing.T) {
		errMock := errors.New("mock GetById")
		usecaseMock.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(nil, errMock)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   "/students/{id}",
			Url:    "/students/9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
		}, restController.GetById)

		var result restserver.Error
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, errMock.Error(), result.Error)
	})

	t.Run("Should return student by id", func(t *testing.T) {
		expected := models.Student{ID: uuid.MustParse("9f9fa978-7df0-4474-b1d4-6be55e0dbd1d"), Name: "Student name 1", Email: "student1@email.com", Birthday: time.Time{}, CreatedAt: time.Time{}}
		usecaseMock.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(&expected, nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   "/students/{id}",
			Url:    "/students/9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
		}, restController.GetById)

		var result models.Student
		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, expected, result)
	})
}

func TestPostStudent(t *testing.T) {
	bodyDTO := `{ "name": "Student name", "email": "student@email.com", "birthday": "1990-12-31T00:00:00Z" }`
	controller := gomock.NewController(t)
	usecaseMock := mock.NewMockIStudentUsecases(controller)
	restController := controllers.StudentController{Usecase: usecaseMock}
	defer controller.Finish()

	t.Run("Should return error when occurred error in DecodeBody (nobody)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/students",
			Url:    "/students",
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (name is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/students",
			Url:    "/students",
			Body:   `{ "email": "student@email.com", "birthday": "1990-12-31T00:00:00Z" }`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (email is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/students",
			Url:    "/students",
			Body:   `{ "name": "Student name", "birthday": "1990-12-31T00:00:00Z" }`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (birthday is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/students",
			Url:    "/students",
			Body:   `{ "name": "Student name", "email": "student@email.com" }`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (birthday is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/students",
			Url:    "/students",
			Body:   `{ "name": "Student name", "email": "student@email.com", "birthday": "invalid" }`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in Create", func(t *testing.T) {
		errMock := errors.New("mock Create")
		usecaseMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errMock)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/students",
			Url:    "/students",
			Body:   bodyDTO,
		}, restController.Post)

		var result restserver.Error
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, errMock.Error(), result.Error)
	})

	t.Run("Should create student", func(t *testing.T) {
		usecaseMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/students",
			Url:    "/students",
			Body:   bodyDTO,
		}, restController.Post)

		assert.Equal(t, http.StatusCreated, response.StatusCode())
	})
}

func TestPutStudent(t *testing.T) {
	bodyDTO := `{ "name": "Student name updated", "email": "student.updated@email.com", "birthday": "1990-02-21T00:00:00Z" }`
	controller := gomock.NewController(t)
	usecaseMock := mock.NewMockIStudentUsecases(controller)
	restController := controllers.StudentController{Usecase: usecaseMock}
	defer controller.Finish()

	t.Run("Should return error when i try get invalid id path param", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/students/{id}",
			Url:    "/students/abc",
		}, restController.Put)

		var result restserver.Error
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return error when occurred error in DecodeBody (nobody)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/students/{id}",
			Url:    "/students/9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
		}, restController.Put)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (name is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/students/{id}",
			Url:    "/students/9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
			Body:   `{ "email": "student@email.com", "birthday": "1990-12-31T00:00:00Z" }`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (email is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/students/{id}",
			Url:    "/students/9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
			Body:   `{ "name": "Student name", "birthday": "1990-12-31T00:00:00Z" }`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (birthday is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/students/{id}",
			Url:    "/students/9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
			Body:   `{ "name": "Student name", "email": "student@email.com" }`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (birthday is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/students/{id}",
			Url:    "/students/9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
			Body:   `{ "name": "Student name", "email": "student@email.com", "birthday": "invalid" }`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in Update", func(t *testing.T) {
		errMock := errors.New("mock Update")
		usecaseMock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errMock)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/students/{id}",
			Url:    "/students/9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
			Body:   bodyDTO,
		}, restController.Put)

		var result restserver.Error
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, errMock.Error(), result.Error)
	})

	t.Run("Should update student", func(t *testing.T) {
		usecaseMock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/students/{id}",
			Url:    "/students/9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
			Body:   bodyDTO,
		}, restController.Put)

		assert.Equal(t, http.StatusNoContent, response.StatusCode())
	})
}

func TestDeleteStudent(t *testing.T) {
	controller := gomock.NewController(t)
	usecaseMock := mock.NewMockIStudentUsecases(controller)
	restController := controllers.StudentController{Usecase: usecaseMock}
	defer controller.Finish()

	t.Run("Should return error when i try get invalid id path param", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   "/students/{id}",
			Url:    "/students/abc",
		}, restController.Delete)

		var result restserver.Error
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return error when occurred error in Delete", func(t *testing.T) {
		errMock := errors.New("mock Delete")
		usecaseMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errMock)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   "/students/{id}",
			Url:    "/students/9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
		}, restController.Delete)

		var result restserver.Error
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, errMock.Error(), result.Error)
	})

	t.Run("Should delete student", func(t *testing.T) {
		usecaseMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   "/students/{id}",
			Url:    "/students/9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
		}, restController.Delete)

		assert.Equal(t, http.StatusNoContent, response.StatusCode())
	})
}

func TestPostStudentDocument(t *testing.T) {
	path := "/students/{id}/upload-document"
	controller := gomock.NewController(t)
	usecaseMock := mock.NewMockIStudentUsecases(controller)
	restController := controllers.StudentController{Usecase: usecaseMock}
	defer controller.Finish()

	t.Run("Should return error when i try get invalid id path param", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    "/students/abc/upload-document",
		}, restController.PostDocument)

		var result restserver.Error
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return error when occurred error in FormFile", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method:  http.MethodPost,
			Headers: map[string]string{"Content-Type": "multipart/form-data"},
			Path:    path,
			Url:     "/students/9f9fa978-7df0-4474-b1d4-6be55e0dbd1d/upload-document",
		}, restController.PostDocument)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
	})

	t.Run("Should  return error when occurred error in FormFile", func(t *testing.T) {
		errMock := errors.New("mock GetById")
		usecaseMock.EXPECT().UploadDocument(gomock.Any(), gomock.Any(), gomock.Any()).Return("", errMock)

		file, err := os.Open("../../../development-environment/files/img.png")

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method:        http.MethodPost,
			Path:          path,
			Url:           "/students/9f9fa978-7df0-4474-b1d4-6be55e0dbd1d/upload-document",
			UploadFile:    file,
			FormFileField: "file",
			FormFileName:  "img.png",
		}, restController.PostDocument)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
	})

	t.Run("Should upload document by student", func(t *testing.T) {
		file, err := os.Open("../../../development-environment/files/img.png")
		assert.NoError(t, err)

		usecaseMock.EXPECT().UploadDocument(gomock.Any(), uuid.MustParse("9f9fa978-7df0-4474-b1d4-6be55e0dbd1d"), gomock.Any()).Return("/file/img.png", nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method:        http.MethodPost,
			Path:          path,
			Url:           "/students/9f9fa978-7df0-4474-b1d4-6be55e0dbd1d/upload-document",
			UploadFile:    file,
			FormFileField: "file",
			FormFileName:  "img.png",
		}, restController.PostDocument)

		assert.Equal(t, http.StatusOK, response.StatusCode())
	})
}
