package controllers

import (
	"net/http"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/usecases"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/restserver"
	"github.com/google/uuid"
)

type StudentController struct {
	Usecase usecases.IStudentUsecases
}

func NewStudentController() *StudentController {
	return &StudentController{
		Usecase: usecases.NewStudentUsecases(),
	}
}

func (c *StudentController) Routes() []restserver.Route {
	return []restserver.Route{
		{
			URI:      "students",
			Method:   http.MethodPost,
			Function: c.Post,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      "students",
			Method:   http.MethodGet,
			Function: c.GetAll,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      "students/{id}",
			Method:   http.MethodGet,
			Function: c.GetById,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      "students/{id}",
			Method:   http.MethodPut,
			Function: c.Put,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      "students/{id}",
			Method:   http.MethodDelete,
			Function: c.Delete,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      "students/{id}/upload-document",
			Method:   http.MethodPost,
			Function: c.PostDocument,
			Prefix:   restserver.PublicApi,
		},
	}
}

// @Summary Get students list
// @Tags students
// @Accept json
// @Produce json
// @Success 200 {array} models.Student
// @Failure 500
// @Param name query string false "name of student"
// @Router /public/students [get]
func (c *StudentController) GetAll(ctx restserver.WebContext) {
	var params models.StudentParams
	ctx.DecodeQueryParams(&params)

	list, err := c.Usecase.GetAll(ctx.Context(), &params)
	if err != nil {
		ctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	ctx.JsonResponse(http.StatusOK, list)
}

// @Summary Get student by id
// @Tags students
// @Accept json
// @Produce json
// @Success 200 {object} models.Student
// @Failure 400
// @Failure 500
// @Param id path string true "Student ID"
// @Router /public/students/{id} [get]
func (c *StudentController) GetById(ctx restserver.WebContext) {
	paramId, err := uuid.Parse(ctx.PathParam("id"))
	if err != nil {
		ctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	result, err := c.Usecase.GetById(ctx.Context(), paramId)
	if err != nil {
		ctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	ctx.JsonResponse(http.StatusOK, result)
}

// @Summary Student create
// @Tags students
// @Accept json
// @Produce json
// @Success 201
// @Failure 422
// @Failure 500
// @Param request body models.StudentCreateUpdateDTO true "request body"
// @Router /public/students [post]
func (c *StudentController) Post(ctx restserver.WebContext) {
	var body models.StudentCreateUpdateDTO
	if err := ctx.DecodeBody(&body); err != nil {
		ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
		return
	}

	if err := c.Usecase.Create(ctx.Context(), &body); err != nil {
		ctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	ctx.EmptyResponse(http.StatusCreated)
}

// @Summary Student update
// @Tags students
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 422
// @Failure 500
// @Param id path string true "Student ID"
// @Param request body models.StudentCreateUpdateDTO true "request body"
// @Router /public/students/{id} [put]
func (c *StudentController) Put(ctx restserver.WebContext) {
	paramId, err := uuid.Parse(ctx.PathParam("id"))
	if err != nil {
		ctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	var body models.StudentCreateUpdateDTO
	if err := ctx.DecodeBody(&body); err != nil {
		ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
		return
	}

	if err = c.Usecase.Update(ctx.Context(), paramId, &body); err != nil {
		ctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	ctx.EmptyResponse(http.StatusNoContent)
}

// @Summary Student delete
// @Tags students
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 500
// @Param id path string true "Student ID"
// @Router /public/students/{id} [delete]
func (c *StudentController) Delete(ctx restserver.WebContext) {
	paramId, err := uuid.Parse(ctx.PathParam("id"))
	if err != nil {
		ctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	if err = c.Usecase.Delete(ctx.Context(), paramId); err != nil {
		ctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	ctx.EmptyResponse(http.StatusNoContent)
}

// PostDocument
// @Summary Upload student document
// @Tags students
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200
// @Failure 400
// @Failure 500
// @Param id path string true "Student ID"
// @Param file formData file true "file path"
// @Router /public/students/{id}/upload-document [post]
func (c *StudentController) PostDocument(ctx restserver.WebContext) {
	paramId, err := uuid.Parse(ctx.PathParam("id"))
	if err != nil {
		ctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	file, _, err := ctx.FormFile("file")
	if err != nil {
		ctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	url, err := c.Usecase.UploadDocument(ctx.Context(), paramId, &file)
	if err != nil {
		ctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	ctx.JsonResponse(http.StatusOK, url)
}
