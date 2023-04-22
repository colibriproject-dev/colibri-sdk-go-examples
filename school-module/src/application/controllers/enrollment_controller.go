package controllers

import (
	"net/http"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/usecases"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/restserver"
)

type EnrollmentController struct {
	Usecase usecases.IEnrollmentUsecases
}

func NewEnrollmentController() *EnrollmentController {
	return &EnrollmentController{
		Usecase: usecases.NewEnrollmentUsecases(),
	}
}

func (c *EnrollmentController) Routes() []restserver.Route {
	return []restserver.Route{
		{
			URI:      "enrollments",
			Method:   http.MethodGet,
			Function: c.GetPage,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      "enrollments",
			Method:   http.MethodPost,
			Function: c.Post,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      "enrollments",
			Method:   http.MethodDelete,
			Function: c.Delete,
			Prefix:   restserver.PublicApi,
		},
	}
}

// @Summary Get enrollments page
// @Tags enrollments
// @Accept json
// @Produce json
// @Success 200 {array} models.Enrollment
// @Failure 400
// @Failure 500
// @Param page query uint16 true "page" minimum(1) default(1)
// @Param pageSize query uint16 true "size of page" minimum(1) default(10)
// @Param studentName query string false "name of student"
// @Param courseName query string false "name of course"
// @Router /public/enrollments [get]
func (c *EnrollmentController) GetPage(ctx restserver.WebContext) {
	var params models.EnrollmentPageParamsDTO
	if err := ctx.DecodeQueryParams(&params); err != nil {
		ctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	page, err := c.Usecase.GetPage(ctx.Context(), &params)
	if err != nil {
		ctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	ctx.JsonResponse(http.StatusOK, page)
}

// @Summary Enrollment create
// @Tags enrollments
// @Accept json
// @Produce json
// @Success 201
// @Failure 422
// @Failure 500
// @Param request body models.EnrollmentCreateDTO true "request body"
// @Router /public/enrollments [post]
func (c *EnrollmentController) Post(ctx restserver.WebContext) {
	var body models.EnrollmentCreateDTO
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

// @Summary Enrollment delete
// @Tags enrollments
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 500
// @Param studentId query string true "ID of student"
// @Param courseId query string true "ID of course"
// @Router /public/enrollments [delete]
func (c *EnrollmentController) Delete(ctx restserver.WebContext) {
	var params models.EnrollmentDeleteParamsDTO
	if err := ctx.DecodeQueryParams(&params); err != nil {
		ctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	if err := c.Usecase.Delete(ctx.Context(), &params); err != nil {
		ctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	ctx.EmptyResponse(http.StatusNoContent)
}
