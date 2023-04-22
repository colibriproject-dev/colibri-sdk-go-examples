package controllers

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/application/controllers"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/usecases/mock"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/types"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/restserver"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEnrollmentController(t *testing.T) {
	t.Run("Should return new enrollment controller", func(t *testing.T) {
		result := controllers.NewEnrollmentController()

		assert.NotNil(t, result)
		assert.NotNil(t, result.Usecase)
		assert.NotNil(t, result.Routes())
	})
}

func TestGetPageEnrollments(t *testing.T) {
	controller := gomock.NewController(t)
	usecaseMock := mock.NewMockIEnrollmentUsecases(controller)
	restController := controllers.EnrollmentController{Usecase: usecaseMock}
	defer controller.Finish()

	t.Run("Should return error when occurred error in DecodeParams", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   "/enrollments",
			Url:    "/enrollments",
		}, restController.GetPage)

		var result restserver.Error
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return error when occurred error in DecodeParams when page is invalid", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   "/enrollments",
			Url:    "/enrollments?page=abc",
		}, restController.GetPage)

		var result restserver.Error
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return error when occurred error in DecodeParams when size is invalid", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   "/enrollments",
			Url:    "/enrollments?pageSize=abc",
		}, restController.GetPage)

		var result restserver.Error
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return error when occurred error in GetPage", func(t *testing.T) {
		errMock := errors.New("mock GetPage")
		usecaseMock.EXPECT().GetPage(gomock.Any(), gomock.Any()).Return(nil, errMock)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   "/enrollments",
			Url:    "/enrollments?page=1&pageSize=10",
		}, restController.GetPage)

		var result restserver.Error
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, errMock.Error(), result.Error)
	})

	t.Run("Should return all enrollments", func(t *testing.T) {
		content := []models.Enrollment{
			{
				Student: models.Student{
					ID:        uuid.MustParse("9f9fa978-7df0-4474-b1d4-6be55e0dbd1d"),
					Name:      "Student name 1",
					Email:     "student1@email.com",
					Birthday:  time.Time{},
					CreatedAt: time.Time{},
				},
				Course: models.Course{
					ID:        uuid.MustParse("8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"),
					Name:      "Course name 1",
					Value:     1000,
					CreatedAt: time.Time{},
				},
				Installments: 3,
				Status:       models.ADIMPLENTE,
				CreatedAt:    time.Time{},
			},
		}
		var expected models.EnrollmentPage = &types.Page[models.Enrollment]{TotalElements: 1, Content: content}

		usecaseMock.EXPECT().GetPage(gomock.Any(), gomock.Any()).Return(expected, nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   "/enrollments",
			Url:    "/enrollments?page=1&pageSize=10",
		}, restController.GetPage)

		var result models.EnrollmentPage
		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, expected, result)
	})
}

func TestPostEnrollment(t *testing.T) {
	bodyDTO := `{ 
		"studentId": "9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
		"courseId": "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
		"installments": 10
	}`
	controller := gomock.NewController(t)
	usecaseMock := mock.NewMockIEnrollmentUsecases(controller)
	restController := controllers.EnrollmentController{Usecase: usecaseMock}
	defer controller.Finish()

	t.Run("Should return error when occurred error in DecodeBody (nobody)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/enrollments",
			Url:    "/enrollments",
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (studentId is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/enrollments",
			Url:    "/enrollments",
			Body:   `{ "courseId": "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d", "value": 10 }`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (courseId is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/enrollments",
			Url:    "/enrollments",
			Body:   `{ "studentId": "9f9fa978-7df0-4474-b1d4-6be55e0dbd1d", "value": 10 }`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (installments is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/enrollments",
			Url:    "/enrollments",
			Body:   `{ "studentId": "9f9fa978-7df0-4474-b1d4-6be55e0dbd1d", "courseId": "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d" }`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (installments is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/enrollments",
			Url:    "/enrollments",
			Body: `{ 
				"studentId": "9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
				"courseId": "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
				"installments": "10"
			}`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in DecodeBody (installments is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/enrollments",
			Url:    "/enrollments",
			Body: `{ 
				"studentId": "9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
				"courseId": "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
				"installments": -1
			}`,
		}, restController.Post)

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return error when occurred error in Create", func(t *testing.T) {
		errMock := errors.New("mock Create")
		usecaseMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errMock)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/enrollments",
			Url:    "/enrollments",
			Body:   bodyDTO,
		}, restController.Post)

		var result restserver.Error
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, errMock.Error(), result.Error)
	})

	t.Run("Should create enrollment", func(t *testing.T) {
		usecaseMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   "/enrollments",
			Url:    "/enrollments",
			Body:   bodyDTO,
		}, restController.Post)

		assert.Equal(t, http.StatusCreated, response.StatusCode())
	})
}

func TestDeleteEnrollment(t *testing.T) {
	controller := gomock.NewController(t)
	usecaseMock := mock.NewMockIEnrollmentUsecases(controller)
	restController := controllers.EnrollmentController{Usecase: usecaseMock}
	defer controller.Finish()

	t.Run("Should return error when occurred error in DecodeParams", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   "/enrollments",
			Url:    "/enrollments",
		}, restController.Delete)

		var result restserver.Error
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return error when occurred error in DecodeParams when studentId is invalid", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   "/enrollments",
			Url:    "/enrollments?studentId=abc",
		}, restController.Delete)

		var result restserver.Error
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return error when occurred error in DecodeParams when courseId is invalid", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   "/enrollments",
			Url:    "/enrollments?courseId=abc",
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
			Path:   "/enrollments",
			Url:    "/enrollments?studentId=8f9fa978-7df0-4474-b1d4-6be55e0dbd0d&courseId=9f0fa978-7df0-4474-b1d4-6be55e0dbd1d",
		}, restController.Delete)

		var result restserver.Error
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.Equal(t, errMock.Error(), result.Error)
	})

	t.Run("Should delete enrollment", func(t *testing.T) {
		usecaseMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   "/enrollments",
			Url:    "/enrollments?studentId=8f9fa978-7df0-4474-b1d4-6be55e0dbd0d&courseId=9f0fa978-7df0-4474-b1d4-6be55e0dbd1d",
		}, restController.Delete)

		assert.Equal(t, http.StatusNoContent, response.StatusCode())
	})
}
