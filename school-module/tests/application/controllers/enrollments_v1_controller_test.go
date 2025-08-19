package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/application/controllers"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	usecasesmock "github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/usecases/mock"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/types"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/web/restserver"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEnrollmentsV1Controller(t *testing.T) {
	t.Run("Should return new enrollments v1 controller", func(t *testing.T) {
		result := controllers.NewEnrollmentsV1Controller()

		assert.NotNil(t, result)
		assert.NotNil(t, result.GetAllPaginatedEnrollmentUsecase)
		assert.NotNil(t, result.CreateEnrollmentUsecase)
		assert.NotNil(t, result.DeleteEnrollmentUsecase)
		assert.NotNil(t, result.UpdateEnrollmentStatusUsecase)
		assert.NotNil(t, result.Routes())
	})
}

func TestEnrollmentsV1Controller_GetAllPaginatedEnrollment(t *testing.T) {
	controller := gomock.NewController(t)
	mockGetAllPaginatedEnrollmentUsecase := usecasesmock.NewMockIGetAllPaginatedEnrollmentUsecase(controller)
	restController := controllers.EnrollmentsV1Controller{GetAllPaginatedEnrollmentUsecase: mockGetAllPaginatedEnrollmentUsecase}
	defer controller.Finish()

	const path string = "/public/v1/enrollments"
	const page uint16 = 1
	const pageSize uint16 = 10
	const studentName string = "testStudentName"
	const courseName string = "testCourseName"

	urlWithParams := fmt.Sprintf("%s?page=%d&pageSize=%d&studentName=%s&courseName=%s", path, page, pageSize, studentName, courseName)
	queryParams := &models.EnrollmentPageParams{
		Page:        page,
		Size:        pageSize,
		StudentName: studentName,
		CourseName:  courseName,
	}

	t.Run("Should return StatusBadRequest when returned error in DecodeParams", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    path,
		}, restController.GetAllPaginatedEnrollment)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeParams (page is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    path + "?pageSize=10",
		}, restController.GetAllPaginatedEnrollment)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeParams (page is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    path + "?page=abc",
		}, restController.GetAllPaginatedEnrollment)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeParams (pageSize is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    path + "?page=10",
		}, restController.GetAllPaginatedEnrollment)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeParams (pageSize is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    path + "?page=10&pageSize=abc",
		}, restController.GetAllPaginatedEnrollment)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusInternalServerError returned error in GetAllPaginatedEnrollmentUsecase", func(t *testing.T) {
		mockErr := errors.New("mock error in GetAllPaginatedEnrollmentUsecase")
		mockGetAllPaginatedEnrollmentUsecase.EXPECT().Execute(gomock.Any(), queryParams).Return(nil, mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    urlWithParams,
		}, restController.GetAllPaginatedEnrollment)

		var result restserver.Error
		assert.EqualValues(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusOK and all enrollments", func(t *testing.T) {
		var expected models.EnrollmentPage = &types.Page[models.Enrollment]{
			TotalItems: 1,
			Items: []models.Enrollment{
				{
					Student: models.Student{
						ID:        uuid.New(),
						Name:      "Student name 1",
						Email:     "student1@email.com",
						Birthday:  types.IsoDate(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
						CreatedAt: time.Now().UTC(),
					},
					Course: models.Course{
						ID:        uuid.New(),
						Name:      "Course name 1",
						Value:     1000,
						CreatedAt: time.Now().UTC(),
					},
					Installments: 3,
					Status:       enums.ADIMPLENTE,
					CreatedAt:    time.Now().UTC(),
				},
			},
		}

		mockGetAllPaginatedEnrollmentUsecase.EXPECT().Execute(gomock.Any(), queryParams).Return(expected, nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    urlWithParams,
		}, restController.GetAllPaginatedEnrollment)

		var result models.EnrollmentPage
		assert.EqualValues(t, http.StatusOK, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
	})
}

func TestEnrollmentsV1Controller_CreateEnrollment(t *testing.T) {
	controller := gomock.NewController(t)
	mockCreateEnrollmentUsecase := usecasesmock.NewMockICreateEnrollmentUsecase(controller)
	restController := controllers.EnrollmentsV1Controller{CreateEnrollmentUsecase: mockCreateEnrollmentUsecase}
	defer controller.Finish()

	const path string = "/public/v1/enrollments"
	const studentId string = "9f9fa978-7df0-4474-b1d4-6be55e0dbd1d"
	const courseId string = "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"
	const installments uint8 = 10

	requestBody := fmt.Sprintf(`{"studentId":"%s","courseId":"%s","installments":%d}`, studentId, courseId, installments)
	enrollmentCreate := &models.EnrollmentCreate{
		StudentID:    uuid.MustParse(studentId),
		CourseID:     uuid.MustParse(courseId),
		Installments: installments,
	}

	t.Run("Should return StatusUnprocessableEntity when returned error in DecodeBody (nobody)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   `{}`,
		}, restController.CreateEnrollment)

		assert.EqualValues(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return StatusUnprocessableEntity when returned error in DecodeBody (studentId is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   `{ "courseId": "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d", "value": 10 }`,
		}, restController.CreateEnrollment)

		assert.EqualValues(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return StatusUnprocessableEntity when returned error in DecodeBody (courseId is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   `{ "studentId": "9f9fa978-7df0-4474-b1d4-6be55e0dbd1d", "value": 10 }`,
		}, restController.CreateEnrollment)

		assert.EqualValues(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return StatusUnprocessableEntity when returned error in DecodeBody (installments is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   `{ "studentId": "9f9fa978-7df0-4474-b1d4-6be55e0dbd1d", "courseId": "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d" }`,
		}, restController.CreateEnrollment)

		assert.EqualValues(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return StatusUnprocessableEntity when returned error in DecodeBody (installments is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body: `{ 
				"studentId": "9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
				"courseId": "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
				"installments": "10"
			}`,
		}, restController.CreateEnrollment)

		assert.EqualValues(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return StatusUnprocessableEntity when returned error in DecodeBody (installments is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body: `{ 
				"studentId": "9f9fa978-7df0-4474-b1d4-6be55e0dbd1d",
				"courseId": "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d",
				"installments": -1
			}`,
		}, restController.CreateEnrollment)

		assert.EqualValues(t, http.StatusUnprocessableEntity, response.StatusCode())
	})

	t.Run("Should return StatusInternalServerError when returned error in CreateEnrollmentUsecase", func(t *testing.T) {
		mockErr := errors.New("mock error in CreateEnrollmentUsecase")
		mockCreateEnrollmentUsecase.EXPECT().Execute(gomock.Any(), enrollmentCreate).Return(mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   requestBody,
		}, restController.CreateEnrollment)

		var result restserver.Error
		assert.EqualValues(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should create enrollment and return StatusCreated", func(t *testing.T) {
		mockCreateEnrollmentUsecase.EXPECT().Execute(gomock.Any(), enrollmentCreate).Return(nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   requestBody,
		}, restController.CreateEnrollment)

		assert.EqualValues(t, http.StatusCreated, response.StatusCode())
	})
}

func TestEnrollmentsV1Controller_DeleteEnrollment(t *testing.T) {
	controller := gomock.NewController(t)
	mockDeleteEnrollmentUsecase := usecasesmock.NewMockIDeleteEnrollmentUsecase(controller)
	restController := controllers.EnrollmentsV1Controller{DeleteEnrollmentUsecase: mockDeleteEnrollmentUsecase}
	defer controller.Finish()

	const path string = "/public/v1/enrollments"
	const studentId string = "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"
	const courseId string = "9f0fa978-7df0-4474-b1d4-6be55e0dbd1d"

	urlWithParams := fmt.Sprintf("%s?studentId=%s&courseId=%s", path, studentId, courseId)
	deleteParams := &models.EnrollmentDelete{
		StudentID: uuid.MustParse(studentId),
		CourseID:  uuid.MustParse(courseId),
	}

	t.Run("Should return StatusBadRequest when returned error in DecodeParams", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   path,
			Url:    path,
		}, restController.DeleteEnrollment)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusBadRequest when occurred returned in DecodeParams (studentId is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   path,
			Url:    path + "?studentId=abc",
		}, restController.DeleteEnrollment)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeParams (courseId is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   path,
			Url:    path + "?courseId=abc",
		}, restController.DeleteEnrollment)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusInternalServerError when returned error in DeleteEnrollmentUsecase", func(t *testing.T) {
		mockErr := errors.New("mock error in DeleteEnrollmentUsecase")
		mockDeleteEnrollmentUsecase.EXPECT().Execute(gomock.Any(), deleteParams).Return(mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   path,
			Url:    urlWithParams,
		}, restController.DeleteEnrollment)

		var result restserver.Error
		assert.EqualValues(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should delete enrollment and return StatusNoContent", func(t *testing.T) {
		mockDeleteEnrollmentUsecase.EXPECT().Execute(gomock.Any(), deleteParams).Return(nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   path,
			Url:    urlWithParams,
		}, restController.DeleteEnrollment)

		assert.EqualValues(t, http.StatusNoContent, response.StatusCode())
	})
}
