package controllers

import (
	"net/http"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/usecases"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/restserver"
	"github.com/google/uuid"
)

type CourseController struct {
	Usecase usecases.ICourseUsecases
}

func NewCourseController() *CourseController {
	return &CourseController{
		Usecase: usecases.NewCourseUsecases(),
	}
}

func (c *CourseController) Routes() []restserver.Route {
	return []restserver.Route{
		{
			URI:      "courses",
			Method:   http.MethodGet,
			Function: c.GetAll,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      "courses/{id}",
			Method:   http.MethodGet,
			Function: c.GetById,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      "courses",
			Method:   http.MethodPost,
			Function: c.Post,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      "courses/{id}",
			Method:   http.MethodPut,
			Function: c.Put,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      "courses/{id}",
			Method:   http.MethodDelete,
			Function: c.Delete,
			Prefix:   restserver.PublicApi,
		},
	}
}

// @Summary Get courses list
// @Tags courses
// @Accept json
// @Produce json
// @Success 200 {array} models.Course
// @Failure 500
// @Router /public/courses [get]
func (c *CourseController) GetAll(ctx restserver.WebContext) {
	result, err := c.Usecase.GetAll(ctx.Context())
	if err != nil {
		ctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	ctx.JsonResponse(http.StatusOK, result)
}

// @Summary Get course by id
// @Tags courses
// @Accept json
// @Produce json
// @Success 200 {object} models.Course
// @Failure 400
// @Failure 500
// @Param id path string true "Course ID"
// @Router /public/courses/{id} [get]
func (c *CourseController) GetById(ctx restserver.WebContext) {
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

// @Summary Course create
// @Tags courses
// @Accept json
// @Produce json
// @Success 201 {object} models.Course
// @Failure 422
// @Failure 500
// @Param request body models.CourseCreateUpdateDTO true "request body"
// @Router /public/courses [post]
func (c *CourseController) Post(ctx restserver.WebContext) {
	var body models.CourseCreateUpdateDTO
	if err := ctx.DecodeBody(&body); err != nil {
		ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
		return
	}

	id, err := c.Usecase.Create(ctx.Context(), &body)
	if err != nil {
		ctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	ctx.JsonResponse(http.StatusCreated, id)
}

// @Summary Course update
// @Tags courses
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 422
// @Failure 500
// @Param id path string true "Course ID"
// @Param request body models.CourseCreateUpdateDTO true "request body"
// @Router /public/courses/{id} [put]
func (c *CourseController) Put(ctx restserver.WebContext) {
	paramId, err := uuid.Parse(ctx.PathParam("id"))
	if err != nil {
		ctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	var body models.CourseCreateUpdateDTO
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

// @Summary Course delete
// @Tags courses
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 500
// @Param id path string true "Course ID"
// @Router /public/courses/{id} [delete]
func (c *CourseController) Delete(ctx restserver.WebContext) {
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
