package controllers

import (
	"net/http"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/usecases"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/web/restserver"
	"github.com/google/uuid"
)

type StudentController struct {
	GetAllPaginatedStudentUsecase usecases.IGetAllPaginatedStudentUsecase
	GetStudentByIdUsecase         usecases.IGetStudentByIdUsecase
	CreateStudentUsecase          usecases.ICreateStudentUsecase
	UpdateStudentUsecase          usecases.IUpdateStudentUsecase
	DeleteStudentUsecase          usecases.IDeleteStudentUsecase
	UploadStudentDocumentUsecase  usecases.IUploadStudentDocumentUsecase
}

func NewStudentController() *StudentController {
	return &StudentController{
		GetAllPaginatedStudentUsecase: usecases.NewGetAllPaginatedStudentUsecase(),
		GetStudentByIdUsecase:         usecases.NewGetStudentByIdUsecase(),
		CreateStudentUsecase:          usecases.NewCreateStudentUsecase(),
		UpdateStudentUsecase:          usecases.NewUpdateStudentUsecase(),
		DeleteStudentUsecase:          usecases.NewDeleteStudentUsecase(),
		UploadStudentDocumentUsecase:  usecases.NewUploadStudentDocumentUsecase(),
	}
}

func (c *StudentController) Routes() []restserver.Route {
	const basePath = "v1/students"
	const basePathWithId = basePath + "/{id}"

	return []restserver.Route{
		{
			URI:      basePath,
			Method:   http.MethodPost,
			Function: c.CreateStudent,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      basePath,
			Method:   http.MethodGet,
			Function: c.GetAllPaginatedStudent,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      basePathWithId,
			Method:   http.MethodGet,
			Function: c.GetStudentById,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      basePathWithId,
			Method:   http.MethodPut,
			Function: c.UpdateStudent,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      basePathWithId,
			Method:   http.MethodDelete,
			Function: c.DeleteStudent,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      basePathWithId + "/upload-document",
			Method:   http.MethodPost,
			Function: c.UploadStudentDocument,
			Prefix:   restserver.PublicApi,
		},
	}
}

// @Summary Get students list
// @Tags students
// @Accept json
// @Produce json
// @Success 200 {object} models.StudentPage
// @Failure 400
// @Failure 500
// @Param name query string false "name of student"
// @Router /public/students [get]
func (c *StudentController) GetAllPaginatedStudent(wctx restserver.WebContext) {
	var params models.StudentPageParams
	if err := wctx.DecodeQueryParams(&params); err != nil {
		wctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	result, err := c.GetAllPaginatedStudentUsecase.Execute(wctx.Context(), &params)
	if err != nil {
		wctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	wctx.JsonResponse(http.StatusOK, result)
}

// @Summary Get student by id
// @Tags students
// @Accept json
// @Produce json
// @Success 200 {object} models.Student
// @Failure 400
// @Failure 404
// @Failure 500
// @Param id path string true "Student ID"
// @Router /public/students/{id} [get]
func (c *StudentController) GetStudentById(wctx restserver.WebContext) {
	paramId, err := uuid.Parse(wctx.PathParam("id"))
	if err != nil {
		wctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	result, err := c.GetStudentByIdUsecase.Execute(wctx.Context(), paramId)
	if err != nil {
		if err.Error() == exceptions.ErrStudentNotFound {
			wctx.ErrorResponse(http.StatusNotFound, err)
			return
		}

		wctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	wctx.JsonResponse(http.StatusOK, result)
}

// @Summary Student create
// @Tags students
// @Accept json
// @Produce json
// @Success 201
// @Failure 409
// @Failure 422
// @Failure 500
// @Param request body models.StudentCreate true "request body"
// @Router /public/students [post]
func (c *StudentController) CreateStudent(wctx restserver.WebContext) {
	var body models.StudentCreate
	if err := wctx.DecodeBody(&body); err != nil {
		wctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	if err := c.CreateStudentUsecase.Execute(wctx.Context(), &body); err != nil {
		if err.Error() == exceptions.ErrStudentAlreadyExists {
			wctx.ErrorResponse(http.StatusConflict, err)
			return
		}

		wctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	wctx.EmptyResponse(http.StatusCreated)
}

// @Summary Student update
// @Tags students
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 409
// @Failure 422
// @Failure 500
// @Param id path string true "Student ID"
// @Param request body models.StudentUpdate true "request body"
// @Router /public/students/{id} [put]
func (c *StudentController) UpdateStudent(wctx restserver.WebContext) {
	paramId, err := uuid.Parse(wctx.PathParam("id"))
	if err != nil {
		wctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	var body models.StudentUpdate
	if err := wctx.DecodeBody(&body); err != nil {
		wctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	body.ID = paramId
	if err = c.UpdateStudentUsecase.Execute(wctx.Context(), &body); err != nil {
		switch err.Error() {
		case exceptions.ErrStudentAlreadyExists:
			wctx.ErrorResponse(http.StatusConflict, err)
		case exceptions.ErrStudentNotFound:
			wctx.ErrorResponse(http.StatusNotFound, err)
		default:
			wctx.ErrorResponse(http.StatusInternalServerError, err)
		}
		return
	}

	wctx.EmptyResponse(http.StatusNoContent)
}

// @Summary Student delete
// @Tags students
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
// @Param id path string true "Student ID"
// @Router /public/students/{id} [delete]
func (c *StudentController) DeleteStudent(wctx restserver.WebContext) {
	paramId, err := uuid.Parse(wctx.PathParam("id"))
	if err != nil {
		wctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	if err = c.DeleteStudentUsecase.Execute(wctx.Context(), paramId); err != nil {
		if err.Error() == exceptions.ErrStudentNotFound {
			wctx.ErrorResponse(http.StatusNotFound, err)
			return
		}

		wctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	wctx.EmptyResponse(http.StatusNoContent)
}

// @Summary Upload student document
// @Tags students
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} models.StudentDocumentUrl
// @Failure 400
// @Failure 404
// @Failure 500
// @Param id path string true "Student ID"
// @Param file formData file true "file path"
// @Router /public/students/{id}/upload-document [post]
func (c *StudentController) UploadStudentDocument(wctx restserver.WebContext) {
	paramId, err := uuid.Parse(wctx.PathParam("id"))
	if err != nil {
		wctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	file, _, err := wctx.FormFile("file")
	if err != nil {
		wctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	url, err := c.UploadStudentDocumentUsecase.Execute(wctx.Context(), paramId, &file)
	if err != nil {
		if err.Error() == exceptions.ErrStudentNotFound {
			wctx.ErrorResponse(http.StatusNotFound, err)
			return
		}

		wctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	wctx.JsonResponse(http.StatusOK, url)
}
