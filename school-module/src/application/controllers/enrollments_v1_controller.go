package controllers

import (
	"net/http"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/usecases"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/web/restserver"
)

type EnrollmentsV1Controller struct {
	GetAllPaginatedEnrollmentUsecase usecases.IGetAllPaginatedEnrollmentUsecase
	CreateEnrollmentUsecase          usecases.ICreateEnrollmentUsecase
	DeleteEnrollmentUsecase          usecases.IDeleteEnrollmentUsecase
	UpdateEnrollmentStatusUsecase    usecases.IUpdateEnrollmentStatusUsecase
}

func NewEnrollmentsV1Controller() *EnrollmentsV1Controller {
	return &EnrollmentsV1Controller{
		GetAllPaginatedEnrollmentUsecase: usecases.NewGetAllPaginatedEnrollmentUsecase(),
		CreateEnrollmentUsecase:          usecases.NewCreateEnrollmentUsecase(),
		DeleteEnrollmentUsecase:          usecases.NewDeleteEnrollmentUsecase(),
		UpdateEnrollmentStatusUsecase:    usecases.NewUpdateEnrollmentStatusUsecase(),
	}
}

func (c *EnrollmentsV1Controller) Routes() []restserver.Route {
	const basePath = "v1/enrollments"

	return []restserver.Route{
		{
			URI:      basePath,
			Method:   http.MethodGet,
			Function: c.GetAllPaginatedEnrollment,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      basePath,
			Method:   http.MethodPost,
			Function: c.CreateEnrollment,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      basePath,
			Method:   http.MethodDelete,
			Function: c.DeleteEnrollment,
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
// @Router /public/v1/enrollments [get]
func (c *EnrollmentsV1Controller) GetAllPaginatedEnrollment(wctx restserver.WebContext) {
	var params models.EnrollmentPageParams
	if err := wctx.DecodeQueryParams(&params); err != nil {
		wctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	result, err := c.GetAllPaginatedEnrollmentUsecase.Execute(wctx.Context(), &params)
	if err != nil {
		wctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	wctx.JsonResponse(http.StatusOK, result)
}

// @Summary Enrollment create
// @Tags enrollments
// @Accept json
// @Produce json
// @Success 201
// @Failure 409
// @Failure 422
// @Failure 500
// @Param request body models.EnrollmentCreate true "request body"
// @Router /public/v1/enrollments [post]
func (c *EnrollmentsV1Controller) CreateEnrollment(wctx restserver.WebContext) {
	var body models.EnrollmentCreate
	if err := wctx.DecodeBody(&body); err != nil {
		wctx.ErrorResponse(http.StatusUnprocessableEntity, err)
		return
	}

	if err := c.CreateEnrollmentUsecase.Execute(wctx.Context(), &body); err != nil {
		if err.Error() == exceptions.ErrEnrollmentAlreadyExists {
			wctx.ErrorResponse(http.StatusConflict, err)
			return
		}

		wctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	wctx.EmptyResponse(http.StatusCreated)
}

// @Summary Enrollment delete
// @Tags enrollments
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
// @Param studentId query string true "ID of student"
// @Param courseId query string true "ID of course"
// @Router /public/v1/enrollments [delete]
func (c *EnrollmentsV1Controller) DeleteEnrollment(wctx restserver.WebContext) {
	var params models.EnrollmentDelete
	if err := wctx.DecodeQueryParams(&params); err != nil {
		wctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	if err := c.DeleteEnrollmentUsecase.Execute(wctx.Context(), &params); err != nil {
		if err.Error() == exceptions.ErrEnrollmentNotFound {
			wctx.ErrorResponse(http.StatusNotFound, err)
			return
		}

		wctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	wctx.EmptyResponse(http.StatusNoContent)
}
