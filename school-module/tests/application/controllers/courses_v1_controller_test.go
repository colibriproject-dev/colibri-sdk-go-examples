package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/application/controllers"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	usecasesmock "github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/usecases/mock"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/web/restserver"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCourseV1Controller(t *testing.T) {
	t.Run("Should return new courses v1 controller", func(t *testing.T) {
		result := controllers.NewCoursesV1Controller()

		assert.NotNil(t, result)
		assert.NotNil(t, result.GetAllCourseUsecase)
		assert.NotNil(t, result.GetCourseByIdUsecase)
		assert.NotNil(t, result.CreateCourseUsecase)
		assert.NotNil(t, result.UpdateCourseUsecase)
		assert.NotNil(t, result.DeleteCourseUsecase)
		assert.NotNil(t, result.Routes())
	})
}

func TestCourseV1Controller_GetAllCourses(t *testing.T) {
	controller := gomock.NewController(t)
	usecaseMock := usecasesmock.NewMockIGetAllCourseUsecase(controller)
	restController := controllers.CoursesV1Controller{GetAllCourseUsecase: usecaseMock}
	defer controller.Finish()

	const path string = "/public/v1/courses"

	t.Run("Should return StatusInternalServerError when general error returned in GetAllCourseUsecase", func(t *testing.T) {
		mockErr := errors.New("mock error in GetAllCourseUsecase")
		usecaseMock.EXPECT().Execute(gomock.Any()).Return(nil, mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    path,
		}, restController.GetAllCourse)

		var result restserver.Error
		assert.EqualValues(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusOK and all courses", func(t *testing.T) {
		expected := []models.Course{
			{ID: uuid.New(), Name: "Course name 1", Value: 1000, CreatedAt: time.Now().UTC()},
			{ID: uuid.New(), Name: "Course name 2", Value: 2000, CreatedAt: time.Now().UTC()},
			{ID: uuid.New(), Name: "Course name 3", Value: 2500, CreatedAt: time.Now().UTC()},
		}
		usecaseMock.EXPECT().Execute(gomock.Any()).Return(expected, nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    path,
		}, restController.GetAllCourse)

		var result []models.Course
		assert.EqualValues(t, http.StatusOK, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
	})
}

func TestCourseV1Controller_GetCourseById(t *testing.T) {
	controller := gomock.NewController(t)
	mockGetCourseByIdUsecase := usecasesmock.NewMockIGetCourseByIdUsecase(controller)
	restController := controllers.CoursesV1Controller{GetCourseByIdUsecase: mockGetCourseByIdUsecase}
	defer controller.Finish()

	const path string = "/public/v1/courses/{id}"
	const id string = "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"
	const url string = "/public/v1/courses/" + id

	t.Run("Should return error when i try get invalid id path param", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    path,
		}, restController.GetCourseById)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusNotFound when ErrCourseNotFound returned in GetById", func(t *testing.T) {
		mockErr := errors.New(exceptions.ErrCourseNotFound)
		mockGetCourseByIdUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    url,
		}, restController.GetCourseById)

		var result restserver.Error
		assert.EqualValues(t, http.StatusNotFound, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusInternalServerError when general error returned in GetCourseByIdUsecase", func(t *testing.T) {
		mockErr := errors.New("mock error in GetCourseByIdUsecase")
		mockGetCourseByIdUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    url,
		}, restController.GetCourseById)

		var result restserver.Error
		assert.EqualValues(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusOK and course by id", func(t *testing.T) {
		expected := models.Course{
			ID:        uuid.MustParse(id),
			Name:      "Course name 1",
			Value:     1000,
			CreatedAt: time.Now().UTC(),
		}
		mockGetCourseByIdUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(&expected, nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    url,
		}, restController.GetCourseById)

		var result models.Course
		assert.EqualValues(t, http.StatusOK, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
	})
}

func TestCourseV1Controller_CreateCourse(t *testing.T) {
	controller := gomock.NewController(t)
	mockCreateCourseUsecase := usecasesmock.NewMockICreateCourseUsecase(controller)
	restController := controllers.CoursesV1Controller{CreateCourseUsecase: mockCreateCourseUsecase}
	defer controller.Finish()

	const path string = "/public/v1/courses"
	const courseName string = "Test course name"
	const courseValue float64 = 1000.00

	requestBody := fmt.Sprintf(`{ "name": "%s", "value": %.2f }`, courseName, courseValue)
	courseCreate := &models.CourseCreate{
		Name:  courseName,
		Value: courseValue,
	}

	t.Run("Should return StatusUnprocessableEntity when returned error in DecodeBody (nobody)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
		}, restController.CreateCourse)

		assert.EqualValues(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return StatusUnprocessableEntity when returned error in DecodeBody (name is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   `{ "value": 1000 }`,
		}, restController.CreateCourse)

		assert.EqualValues(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return StatusUnprocessableEntity when returned error in DecodeBody (value is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   `{ "name": "Course name" }`,
		}, restController.CreateCourse)

		assert.EqualValues(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return StatusUnprocessableEntity when returned error in DecodeBody (value is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   `{ "name": "Course name", "value": "1000" }`,
		}, restController.CreateCourse)

		assert.EqualValues(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return StatusConflict when returned ErrCourseAlreadyExists in CreateCourseUsecase", func(t *testing.T) {
		mockErr := errors.New(exceptions.ErrCourseAlreadyExists)
		mockCreateCourseUsecase.EXPECT().Execute(gomock.Any(), courseCreate).Return(nil, mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   requestBody,
		}, restController.CreateCourse)

		var result restserver.Error
		assert.EqualValues(t, http.StatusConflict, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusInternalServerError when returned general error in CreateCourseUsecase", func(t *testing.T) {
		mockErr := errors.New("mock error in CreateCourseUsecase")
		mockCreateCourseUsecase.EXPECT().Execute(gomock.Any(), courseCreate).Return(nil, mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   requestBody,
		}, restController.CreateCourse)

		var result restserver.Error
		assert.EqualValues(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should create course and return StatusCreated with created course data", func(t *testing.T) {
		expected := models.Course{
			ID:        uuid.New(),
			Name:      "Course created name",
			Value:     1000,
			CreatedAt: time.Now().UTC(),
		}
		mockCreateCourseUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(&expected, nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   requestBody,
		}, restController.CreateCourse)

		var result models.Course
		assert.EqualValues(t, http.StatusCreated, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
	})
}

func TestCourseV1Controller_UpdateCourse(t *testing.T) {
	controller := gomock.NewController(t)
	mockUpdateCourseUsecase := usecasesmock.NewMockIUpdateCourseUsecase(controller)
	restController := controllers.CoursesV1Controller{UpdateCourseUsecase: mockUpdateCourseUsecase}
	defer controller.Finish()

	const path string = "/public/v1/courses/{id}"
	const id string = "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"
	const url string = "/public/v1/courses/" + id
	const courseName string = "Course name updated"
	const courseValue float64 = 1200.00

	requestBody := fmt.Sprintf(`{ "name": "%s", "value": %.2f }`, courseName, courseValue)
	courseUpdate := &models.CourseUpdate{
		ID:    uuid.MustParse(id),
		Name:  courseName,
		Value: courseValue,
	}

	t.Run("Should return StatusBadRequest when i try get with invalid id path param", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    "/public/v1/courses/abc",
		}, restController.UpdateCourse)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusUnprocessableEntity when returned error in DecodeBody (nobody)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
		}, restController.UpdateCourse)

		assert.EqualValues(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return StatusUnprocessableEntity when returned error in DecodeBody (name is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
			Body:   `{ "value": 1200 }`,
		}, restController.UpdateCourse)

		assert.EqualValues(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return StatusUnprocessableEntity when returned error in DecodeBody (value is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
			Body:   `{ "name": "Course name updated" }`,
		}, restController.UpdateCourse)

		assert.EqualValues(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return StatusUnprocessableEntity when returned error in DecodeBody (value is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
			Body:   `{ "name": "Course name updated", "value": "1200" }`,
		}, restController.UpdateCourse)

		assert.EqualValues(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return StatusConflict when ErrCourseAlreadyExists returned in UpdateCourseUsecase", func(t *testing.T) {
		mockErr := errors.New(exceptions.ErrCourseAlreadyExists)
		mockUpdateCourseUsecase.EXPECT().Execute(gomock.Any(), courseUpdate).Return(mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
			Body:   requestBody,
		}, restController.UpdateCourse)

		var result restserver.Error
		assert.EqualValues(t, http.StatusConflict, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusNotFound when ErrCourseNotFound returned in UpdateCourseUsecase", func(t *testing.T) {
		mockErr := errors.New(exceptions.ErrCourseNotFound)
		mockUpdateCourseUsecase.EXPECT().Execute(gomock.Any(), courseUpdate).Return(mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
			Body:   requestBody,
		}, restController.UpdateCourse)

		var result restserver.Error
		assert.EqualValues(t, http.StatusNotFound, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusInternalServerError when general error returned in UpdateCourseUsecase", func(t *testing.T) {
		mockErr := errors.New("mock error in UpdateCourseUsecase")
		mockUpdateCourseUsecase.EXPECT().Execute(gomock.Any(), courseUpdate).Return(mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
			Body:   requestBody,
		}, restController.UpdateCourse)

		var result restserver.Error
		assert.EqualValues(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should update course and return StatusNoContent", func(t *testing.T) {
		mockUpdateCourseUsecase.EXPECT().Execute(gomock.Any(), courseUpdate).Return(nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
			Body:   requestBody,
		}, restController.UpdateCourse)

		assert.EqualValues(t, http.StatusNoContent, response.StatusCode())
	})
}

func TestCourseV1Controller_DeleteCourse(t *testing.T) {
	controller := gomock.NewController(t)
	mockDeleteCourseUsecase := usecasesmock.NewMockIDeleteCourseUsecase(controller)
	restController := controllers.CoursesV1Controller{DeleteCourseUsecase: mockDeleteCourseUsecase}
	defer controller.Finish()

	const path string = "/public/v1/courses/{id}"
	const id string = "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"
	const url string = "/public/v1/courses/" + id

	t.Run("Should return StatusBadRequest when i try get invalid id path param", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   path,
			Url:    "/public/v1/courses/abc",
		}, restController.DeleteCourse)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusNotFound when ErrCourseNotFound returned in DeleteCourseUsecase", func(t *testing.T) {
		mockErr := errors.New(exceptions.ErrCourseNotFound)
		mockDeleteCourseUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   path,
			Url:    url,
		}, restController.DeleteCourse)

		var result restserver.Error
		assert.EqualValues(t, http.StatusNotFound, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusInternalServerError when general error returned in DeleteCourseUsecase", func(t *testing.T) {
		mockErr := errors.New("mock error in DeleteCourseUsecase")
		mockDeleteCourseUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   path,
			Url:    url,
		}, restController.DeleteCourse)

		var result restserver.Error
		assert.EqualValues(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should delete course and return StatusNoContent", func(t *testing.T) {
		mockDeleteCourseUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   path,
			Url:    url,
		}, restController.DeleteCourse)

		assert.EqualValues(t, http.StatusNoContent, response.StatusCode())
	})
}
