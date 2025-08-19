package controllers

import (
	"net/http"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/usecases"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/web/restserver"
	"github.com/google/uuid"
)

type CoursesV1Controller struct {
	GetAllCourseUsecase  usecases.IGetAllCourseUsecase
	GetCourseByIdUsecase usecases.IGetCourseByIdUsecase
	CreateCourseUsecase  usecases.ICreateCourseUsecase
	UpdateCourseUsecase  usecases.IUpdateCourseUsecase
	DeleteCourseUsecase  usecases.IDeleteCourseUsecase
}

func NewCoursesV1Controller() *CoursesV1Controller {
	return &CoursesV1Controller{
		GetAllCourseUsecase:  usecases.NewGetAllCourseUsecase(),
		GetCourseByIdUsecase: usecases.NewGetCourseByIdUsecase(),
		CreateCourseUsecase:  usecases.NewCreateCourseUsecase(),
		UpdateCourseUsecase:  usecases.NewUpdateCourseUsecase(),
		DeleteCourseUsecase:  usecases.NewDeleteCourseUsecase(),
	}
}

func (c *CoursesV1Controller) Routes() []restserver.Route {
	const basePath = "v1/courses"
	const basePathWithId = basePath + "/{id}"

	return []restserver.Route{
		{
			URI:      basePath,
			Method:   http.MethodGet,
			Function: c.GetAllCourse,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      basePathWithId,
			Method:   http.MethodGet,
			Function: c.GetCourseById,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      basePath,
			Method:   http.MethodPost,
			Function: c.CreateCourse,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      basePathWithId,
			Method:   http.MethodPut,
			Function: c.UpdateCourse,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      basePathWithId,
			Method:   http.MethodDelete,
			Function: c.DeleteCourse,
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
// @Router /public/v1/courses [get]
func (c *CoursesV1Controller) GetAllCourse(wctx restserver.WebContext) {
	result, err := c.GetAllCourseUsecase.Execute(wctx.Context())
	if err != nil {
		wctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	wctx.JsonResponse(http.StatusOK, result)
}

// @Summary Get course by id
// @Tags courses
// @Accept json
// @Produce json
// @Success 200 {object} models.Course
// @Failure 400
// @Failure 404
// @Failure 500
// @Param id path string true "Course ID"
// @Router /public/v1/courses/{id} [get]
func (c *CoursesV1Controller) GetCourseById(wctx restserver.WebContext) {
	paramId, err := uuid.Parse(wctx.PathParam("id"))
	if err != nil {
		wctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	result, err := c.GetCourseByIdUsecase.Execute(wctx.Context(), paramId)
	if err != nil {
		if err.Error() == exceptions.ErrCourseNotFound {
			wctx.ErrorResponse(http.StatusNotFound, err)
			return
		}

		wctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	wctx.JsonResponse(http.StatusOK, result)
}

// @Summary Course create
// @Tags courses
// @Accept json
// @Produce json
// @Success 201 {object} models.Course
// @Failure 409
// @Failure 422
// @Failure 500
// @Param request body models.CourseCreate true "request body"
// @Router /public/v1/courses [post]
func (c *CoursesV1Controller) CreateCourse(wctx restserver.WebContext) {
	var body models.CourseCreate
	if err := wctx.DecodeBody(&body); err != nil {
		wctx.ErrorResponse(http.StatusUnprocessableEntity, err)
		return
	}

	result, err := c.CreateCourseUsecase.Execute(wctx.Context(), &body)
	if err != nil {
		if err.Error() == exceptions.ErrCourseAlreadyExists {
			wctx.ErrorResponse(http.StatusConflict, err)
			return
		}

		wctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	wctx.JsonResponse(http.StatusCreated, result)
}

// @Summary Course update
// @Tags courses
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 409
// @Failure 422
// @Failure 500
// @Param id path string true "Course ID"
// @Param request body models.CourseUpdate true "request body"
// @Router /public/v1/courses/{id} [put]
func (c *CoursesV1Controller) UpdateCourse(wctx restserver.WebContext) {
	paramId, err := uuid.Parse(wctx.PathParam("id"))
	if err != nil {
		wctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	var body models.CourseUpdate
	if err := wctx.DecodeBody(&body); err != nil {
		wctx.ErrorResponse(http.StatusUnprocessableEntity, err)
		return
	}

	body.ID = paramId
	if err = c.UpdateCourseUsecase.Execute(wctx.Context(), &body); err != nil {
		switch err.Error() {
		case exceptions.ErrCourseAlreadyExists:
			wctx.ErrorResponse(http.StatusConflict, err)
		case exceptions.ErrCourseNotFound:
			wctx.ErrorResponse(http.StatusNotFound, err)
		default:
			wctx.ErrorResponse(http.StatusInternalServerError, err)
		}
		return
	}

	wctx.EmptyResponse(http.StatusNoContent)
}

// @Summary Course delete
// @Tags courses
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
// @Param id path string true "Course ID"
// @Router /public/v1/courses/{id} [delete]
func (c *CoursesV1Controller) DeleteCourse(wctx restserver.WebContext) {
	paramId, err := uuid.Parse(wctx.PathParam("id"))
	if err != nil {
		wctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	if err = c.DeleteCourseUsecase.Execute(wctx.Context(), paramId); err != nil {
		if err.Error() == exceptions.ErrCourseNotFound {
			wctx.ErrorResponse(http.StatusNotFound, err)
			return
		}

		wctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	wctx.EmptyResponse(http.StatusNoContent)
}
