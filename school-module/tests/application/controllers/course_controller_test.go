package controllers

import (
	"errors"
	"net/http"
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

func TestCourseController(t *testing.T) {
	t.Run("Should return new course controller", func(t *testing.T) {
		result := controllers.NewCourseController()

		assert.NotNil(t, result)
		assert.NotNil(t, result.Usecase)
		assert.NotNil(t, result.Routes())
	})
}

func TestGetAllCourses(t *testing.T) {
	controller := gomock.NewController(t)
	usecaseMock := mock.NewMockICourseUsecases(controller)
	restController := controllers.CourseController{Usecase: usecaseMock}
	defer controller.Finish()

	t.Run("Should return error when occurred error in GetAll", func(t *testing.T) {
		errMock := errors.New("mock GetAll")
		usecaseMock.EXPECT().GetAll(gomock.Any()).Return(nil, errMock)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   "/courses",
			Url:    "/courses",
		}, restController.GetAll)

		var result restserver.Error
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, errMock.Error(), result.Error)
	})

	t.Run("Should return all courses", func(t *testing.T) {
		expected := []models.Course{
			{ID: uuid.MustParse("8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"), Name: "Course name 1", Value: 1000, CreatedAt: time.Time{}},
			{ID: uuid.MustParse("5979def8-f9c2-4258-bbf1-118c1ffb0a56"), Name: "Course name 2", Value: 2000, CreatedAt: time.Time{}},
			{ID: uuid.MustParse("1b61d66e-1e6b-48fa-921c-eb96d009f8f9"), Name: "Course name 3", Value: 2500, CreatedAt: time.Time{}},
		}
		usecaseMock.EXPECT().GetAll(gomock.Any()).Return(expected, nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   "/courses",
			Url:    "/courses",
		}, restController.GetAll)

		var result []models.Course
		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, expected, result)
	})
}

func TestGetCourseById(t *testing.T) {
	controller := gomock.NewController(t)
	usecaseMock := mock.NewMockICourseUsecases(controller)
	restController := controllers.CourseController{Usecase: usecaseMock}
	defer controller.Finish()

	t.Run("Should return error when i try get invalid id path param", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   "/courses/{id}",
			Url:    "/courses/abc",
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
			Path:   "/courses/{id}",
			Url:    "/courses/8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
		}, restController.GetById)

		var result restserver.Error
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, errMock.Error(), result.Error)
	})

	t.Run("Should return course by id", func(t *testing.T) {
		expected := models.Course{
			ID:        uuid.MustParse("8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"),
			Name:      "Course name 1",
			Value:     1000,
			CreatedAt: time.Time{},
		}
		usecaseMock.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(&expected, nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   "/courses/{id}",
			Url:    "/courses/8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
		}, restController.GetById)

		var result models.Course
		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, expected, result)
	})
}

func TestPostCourse(t *testing.T) {
	bodyDTO := `{ "name": "Course name", "value": 1000 }`
	controller := gomock.NewController(t)
	usecaseMock := mock.NewMockICourseUsecases(controller)
	restController := controllers.CourseController{Usecase: usecaseMock}
	defer controller.Finish()

	t.Run("Should return error when occurred error in DecodeBody (nobody)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/courses",
			Url:    "/courses",
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (name is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/courses",
			Url:    "/courses",
			Body:   `{ "value": 1000 }`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (value is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/courses",
			Url:    "/courses",
			Body:   `{ "name": "Course name" }`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (value is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/courses",
			Url:    "/courses",
			Body:   `{ "name": "Course name", "value": "1000" }`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in Create", func(t *testing.T) {
		errMock := errors.New("mock Create")
		usecaseMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errMock)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/courses",
			Url:    "/courses",
			Body:   bodyDTO,
		}, restController.Post)

		var result restserver.Error
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, errMock.Error(), result.Error)
	})

	t.Run("Should create course", func(t *testing.T) {
		expected := models.Course{ID: uuid.MustParse("8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"), Name: "Course name", Value: 1000, CreatedAt: time.Time{}}
		usecaseMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&expected, nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/courses",
			Url:    "/courses",
			Body:   bodyDTO,
		}, restController.Post)

		var result models.Course
		assert.Equal(t, http.StatusCreated, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, expected, result)
	})
}

func TestPutCourse(t *testing.T) {
	bodyDTO := `{ "name": "Course name updated", "value": 1200 }`
	controller := gomock.NewController(t)
	usecaseMock := mock.NewMockICourseUsecases(controller)
	restController := controllers.CourseController{Usecase: usecaseMock}
	defer controller.Finish()

	t.Run("Should return error when i try get invalid id path param", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/courses/{id}",
			Url:    "/courses/abc",
		}, restController.Put)

		var result restserver.Error
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return error when occurred error in DecodeBody (nobody)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/courses/{id}",
			Url:    "/courses/8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
		}, restController.Put)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (name is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/courses/{id}",
			Url:    "/courses/8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
			Body:   `{ "value": 1200 }`,
		}, restController.Put)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (value is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/courses/{id}",
			Url:    "/courses/8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
			Body:   `{ "name": "Course name updated" }`,
		}, restController.Put)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (value is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/courses/{id}",
			Url:    "/courses/8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
			Body:   `{ "name": "Course name updated", "value": "1200" }`,
		}, restController.Put)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in Update", func(t *testing.T) {
		errMock := errors.New("mock Update")
		usecaseMock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errMock)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/courses/{id}",
			Url:    "/courses/8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
			Body:   bodyDTO,
		}, restController.Put)

		var result restserver.Error
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, errMock.Error(), result.Error)
	})

	t.Run("Should update course", func(t *testing.T) {
		usecaseMock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   "/courses/{id}",
			Url:    "/courses/8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
			Body:   bodyDTO,
		}, restController.Put)

		assert.Equal(t, http.StatusNoContent, response.StatusCode())
	})
}

func TestDeleteCourse(t *testing.T) {
	controller := gomock.NewController(t)
	usecaseMock := mock.NewMockICourseUsecases(controller)
	restController := controllers.CourseController{Usecase: usecaseMock}
	defer controller.Finish()

	t.Run("Should return error when i try get invalid id path param", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   "/courses/{id}",
			Url:    "/courses/abc",
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
			Path:   "/courses/{id}",
			Url:    "/courses/8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
		}, restController.Delete)

		var result restserver.Error
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, errMock.Error(), result.Error)
	})

	t.Run("Should delete course", func(t *testing.T) {
		usecaseMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   "/courses/{id}",
			Url:    "/courses/8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
		}, restController.Delete)

		assert.Equal(t, http.StatusNoContent, response.StatusCode())
	})
}
