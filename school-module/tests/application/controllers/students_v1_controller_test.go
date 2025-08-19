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
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/types"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/web/restserver"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStudentController(t *testing.T) {
	t.Run("Should return new students v1 controller", func(t *testing.T) {
		result := controllers.NewStudentController()

		assert.NotNil(t, result)
		assert.NotNil(t, result.GetAllPaginatedStudentUsecase)
		assert.NotNil(t, result.GetStudentByIdUsecase)
		assert.NotNil(t, result.CreateStudentUsecase)
		assert.NotNil(t, result.UpdateStudentUsecase)
		assert.NotNil(t, result.DeleteStudentUsecase)
		assert.NotNil(t, result.UploadStudentDocumentUsecase)
		assert.NotNil(t, result.Routes())
	})
}

func TestStudentController_GetAllPaginatedStudent(t *testing.T) {
	controller := gomock.NewController(t)
	mockGetAllPaginatedStudentUsecase := usecasesmock.NewMockIGetAllPaginatedStudentUsecase(controller)
	restController := controllers.StudentController{GetAllPaginatedStudentUsecase: mockGetAllPaginatedStudentUsecase}
	defer controller.Finish()

	const path string = "/public/v1/students"
	const page uint16 = 1
	const pageSize uint16 = 10
	const studentName string = "testStudentName"

	urlWithParams := fmt.Sprintf("%s?page=%d&pageSize=%d&name=%s", path, page, pageSize, studentName)
	queryParams := &models.StudentPageParams{
		Page: page,
		Size: pageSize,
		Name: studentName,
	}

	t.Run("Should return StatusBadRequest when returned error in DecodeQueryParams", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    path,
		}, restController.GetAllPaginatedStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeQueryParams (page is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    path + "?pageSize=10",
		}, restController.GetAllPaginatedStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeQueryParams (page is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    path + "?page=abc",
		}, restController.GetAllPaginatedStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeQueryParams (pageSize is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    path + "?page=1",
		}, restController.GetAllPaginatedStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeQueryParams (pageSize is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    path + "?page=1&pageSize=abc",
		}, restController.GetAllPaginatedStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusInternalServerError when general error returned in GetAllPaginatedStudentUsecase", func(t *testing.T) {
		mockErr := errors.New("mock error in GetAllPaginatedStudentUsecase")
		mockGetAllPaginatedStudentUsecase.EXPECT().Execute(gomock.Any(), queryParams).Return(nil, mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    urlWithParams,
		}, restController.GetAllPaginatedStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusOK and paginated students", func(t *testing.T) {
		expectedStudents := []models.Student{
			{ID: uuid.New(), Name: "Student 1", Email: "student1@test.com", Birthday: types.IsoDate(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)), CreatedAt: time.Date(2023, 8, 24, 9, 0, 0, 0, time.UTC)},
			{ID: uuid.New(), Name: "Student 2", Email: "student2@test.com", Birthday: types.IsoDate(time.Date(1991, 2, 1, 0, 0, 0, 0, time.UTC)), CreatedAt: time.Date(2023, 8, 25, 15, 0, 0, 0, time.UTC)},
		}
		expectedPage := &types.Page[models.Student]{
			TotalItems: uint64(len(expectedStudents)),
			Items:      expectedStudents,
		}
		mockGetAllPaginatedStudentUsecase.EXPECT().Execute(gomock.Any(), queryParams).Return(expectedPage, nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    urlWithParams,
		}, restController.GetAllPaginatedStudent)

		var result types.Page[models.Student]
		assert.EqualValues(t, http.StatusOK, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, expectedPage.Items, result.Items)
		assert.EqualValues(t, expectedPage.TotalItems, result.TotalItems)
	})
}

func TestStudentController_GetStudentById(t *testing.T) {
	controller := gomock.NewController(t)
	mockGetStudentByIdUsecase := usecasesmock.NewMockIGetStudentByIdUsecase(controller)
	restController := controllers.StudentController{GetStudentByIdUsecase: mockGetStudentByIdUsecase}
	defer controller.Finish()

	const path string = "/public/v1/students/{id}"
	const id string = "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"
	const url string = "/public/v1/students/" + id

	t.Run("Should return error when i try get invalid id path param", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    path,
		}, restController.GetStudentById)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusNotFound when ErrStudentNotFound returned in GetById", func(t *testing.T) {
		mockErr := errors.New(exceptions.ErrStudentNotFound)
		mockGetStudentByIdUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    url,
		}, restController.GetStudentById)

		var result restserver.Error
		assert.EqualValues(t, http.StatusNotFound, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusInternalServerError when general error returned in GetStudentByIdUsecase", func(t *testing.T) {
		mockErr := errors.New("mock error in GetStudentByIdUsecase")
		mockGetStudentByIdUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    url,
		}, restController.GetStudentById)

		var result restserver.Error
		assert.EqualValues(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusOK and student by id", func(t *testing.T) {
		expected := models.Student{
			ID:        uuid.MustParse(id),
			Name:      "Student Test Name",
			Email:     "test@student.com",
			Birthday:  types.IsoDate(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)),
			CreatedAt: time.Now().UTC(),
		}
		mockGetStudentByIdUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(&expected, nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodGet,
			Path:   path,
			Url:    url,
		}, restController.GetStudentById)

		var result models.Student
		assert.EqualValues(t, http.StatusOK, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
	})
}

func TestStudentController_CreateStudent(t *testing.T) {
	controller := gomock.NewController(t)
	mockCreateStudentUsecase := usecasesmock.NewMockICreateStudentUsecase(controller)
	restController := controllers.StudentController{CreateStudentUsecase: mockCreateStudentUsecase}
	defer controller.Finish()

	const path string = "/public/v1/students"
	const studentName string = "Test Student Name"
	const studentEmail string = "test@student.com"
	const studentBirthday string = "1990-01-01"

	requestBody := fmt.Sprintf(`{ "name": "%s", "email": "%s", "birthday": "%s" }`, studentName, studentEmail, studentBirthday)

	t.Run("Should return StatusBadRequest when returned error in DecodeBody (nobody)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
		}, restController.CreateStudent)

		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeBody (name is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   `{ "email": "test@student.com", "birthday": "1990-01-01" }`,
		}, restController.CreateStudent)

		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeBody (email is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   `{ "name": "Student Test", "birthday": "1990-01-01" }`,
		}, restController.CreateStudent)

		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeBody (birthday is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   `{ "name": "Student Test", "email": "test@student.com" }`,
		}, restController.CreateStudent)

		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeBody (birthday is invalid)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   `{ "name": "Student Test", "email": "test@student.com", "birthday": "invalid-date" }`,
		}, restController.CreateStudent)

		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
	})

	t.Run("Should return StatusConflict when returned ErrStudentAlreadyExists in CreateStudentUsecase", func(t *testing.T) {
		mockErr := errors.New(exceptions.ErrStudentAlreadyExists)
		mockCreateStudentUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   requestBody,
		}, restController.CreateStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusConflict, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusInternalServerError when returned general error in CreateStudentUsecase", func(t *testing.T) {
		mockErr := errors.New("mock error in CreateStudentUsecase")
		mockCreateStudentUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   requestBody,
		}, restController.CreateStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should create student and return StatusCreated", func(t *testing.T) {
		mockCreateStudentUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
			Body:   requestBody,
		}, restController.CreateStudent)

		assert.EqualValues(t, http.StatusCreated, response.StatusCode())
	})
}

func TestStudentController_UpdateStudent(t *testing.T) {
	controller := gomock.NewController(t)
	mockUpdateStudentUsecase := usecasesmock.NewMockIUpdateStudentUsecase(controller)
	restController := controllers.StudentController{UpdateStudentUsecase: mockUpdateStudentUsecase}
	defer controller.Finish()

	const path string = "/public/v1/students/{id}"
	const id string = "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"
	const url string = "/public/v1/students/" + id
	const studentName string = "Student name updated"
	const studentEmail string = "updated@student.com"
	const studentBirthday string = "1990-01-01"

	requestBody := fmt.Sprintf(`{ "name": "%s", "email": "%s", "birthday": "%s" }`, studentName, studentEmail, studentBirthday)
	studentUpdate := &models.StudentUpdate{
		ID:       uuid.MustParse(id),
		Name:     studentName,
		Email:    studentEmail,
		Birthday: types.IsoDate(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	t.Run("Should return StatusBadRequest when i try get with invalid id path param", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    path,
			Body:   requestBody,
		}, restController.UpdateStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeBody (nobody)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
		}, restController.UpdateStudent)

		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeBody (name is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
			Body:   `{ "email": "test@student.com", "birthday": "1990-01-01" }`,
		}, restController.UpdateStudent)

		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeBody (email is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
			Body:   `{ "name": "Student Test", "birthday": "1990-01-01" }`,
		}, restController.UpdateStudent)

		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
	})

	t.Run("Should return StatusBadRequest when returned error in DecodeBody (birthday is empty)", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
			Body:   `{ "name": "Student Test", "email": "test@student.com" }`,
		}, restController.UpdateStudent)

		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
	})

	t.Run("Should return StatusConflict when returned ErrStudentAlreadyExists in UpdateStudentUsecase", func(t *testing.T) {
		mockErr := errors.New(exceptions.ErrStudentAlreadyExists)
		mockUpdateStudentUsecase.EXPECT().Execute(gomock.Any(), studentUpdate).Return(mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
			Body:   requestBody,
		}, restController.UpdateStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusConflict, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusNotFound when returned ErrStudentNotFound in UpdateStudentUsecase", func(t *testing.T) {
		mockErr := errors.New(exceptions.ErrStudentNotFound)
		mockUpdateStudentUsecase.EXPECT().Execute(gomock.Any(), studentUpdate).Return(mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
			Body:   requestBody,
		}, restController.UpdateStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusNotFound, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusInternalServerError when returned general error in UpdateStudentUsecase", func(t *testing.T) {
		mockErr := errors.New("mock error in UpdateStudentUsecase")
		mockUpdateStudentUsecase.EXPECT().Execute(gomock.Any(), studentUpdate).Return(mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
			Body:   requestBody,
		}, restController.UpdateStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should update student and return StatusNoContent", func(t *testing.T) {
		mockUpdateStudentUsecase.EXPECT().Execute(gomock.Any(), studentUpdate).Return(nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPut,
			Path:   path,
			Url:    url,
			Body:   requestBody,
		}, restController.UpdateStudent)

		assert.EqualValues(t, http.StatusNoContent, response.StatusCode())
	})
}

func TestStudentController_DeleteStudent(t *testing.T) {
	controller := gomock.NewController(t)
	mockDeleteStudentUsecase := usecasesmock.NewMockIDeleteStudentUsecase(controller)
	restController := controllers.StudentController{DeleteStudentUsecase: mockDeleteStudentUsecase}
	defer controller.Finish()

	const path string = "/public/v1/students/{id}"
	const id string = "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"
	const url string = "/public/v1/students/" + id

	t.Run("Should return StatusBadRequest when i try delete with invalid id path param", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   path,
			Url:    path,
		}, restController.DeleteStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusNotFound when returned ErrStudentNotFound in DeleteStudentUsecase", func(t *testing.T) {
		mockErr := errors.New(exceptions.ErrStudentNotFound)
		mockDeleteStudentUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   path,
			Url:    url,
		}, restController.DeleteStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusNotFound, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should return StatusInternalServerError when returned general error in DeleteStudentUsecase", func(t *testing.T) {
		mockErr := errors.New("mock error in DeleteStudentUsecase")
		mockDeleteStudentUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(mockErr)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   path,
			Url:    url,
		}, restController.DeleteStudent)

		var result restserver.Error
		assert.EqualValues(t, http.StatusInternalServerError, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
		assert.EqualValues(t, mockErr.Error(), result.Error)
	})

	t.Run("Should delete student and return StatusNoContent", func(t *testing.T) {
		mockDeleteStudentUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)

		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodDelete,
			Path:   path,
			Url:    url,
		}, restController.DeleteStudent)

		assert.EqualValues(t, http.StatusNoContent, response.StatusCode())
	})
}

func TestStudentController_UploadStudentDocument(t *testing.T) {
	controller := gomock.NewController(t)
	mockUploadStudentDocumentUsecase := usecasesmock.NewMockIUploadStudentDocumentUsecase(controller)
	restController := controllers.StudentController{UploadStudentDocumentUsecase: mockUploadStudentDocumentUsecase}
	defer controller.Finish()

	const path string = "/public/v1/students/{id}/upload-document"
	const id string = "8f9fa978-7df0-4474-b1d4-6be55e0dbd0d"
	const url string = "/public/v1/students/" + id + "/upload-document"

	t.Run("Should return StatusBadRequest when i try upload with invalid id path param", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    path,
		}, restController.UploadStudentDocument)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	t.Run("Should return StatusBadRequest when no file is provided", func(t *testing.T) {
		response := restserver.NewRequestTest(&restserver.RequestTest{
			Method: http.MethodPost,
			Path:   path,
			Url:    url,
		}, restController.UploadStudentDocument)

		var result restserver.Error
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode())
		assert.NoError(t, response.DecodeBody(&result))
		assert.NotNil(t, result)
	})

	// Note: Testing file upload endpoints with restserver.RequestTest is complex
	// so we focus on testing the business logic through the usecase layer
	// The controller logic for successful uploads would be tested at integration level
}
