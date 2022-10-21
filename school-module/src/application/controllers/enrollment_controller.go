package controllers

import (
	"net/http"
	"school-module/src/domain/models"
	"school-module/src/domain/usecases"
	"strings"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/webrest"
	"github.com/google/uuid"
)

type IEnrollmentController interface {
	Routes() []webrest.Route
	GetAll(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type EnrollmentController struct {
	Usecase usecases.IEnrollmentUsecase
}

func NewEnrollmentController() {
	controller := &EnrollmentController{
		Usecase: usecases.NewEnrollmentUsecase(),
	}

	webrest.AddRoutes(controller.Routes())
}

func (c *EnrollmentController) Routes() []webrest.Route {
	return []webrest.Route{
		{
			URI:      "enrollments",
			Method:   http.MethodPost,
			Function: c.Post,
			Prefix:   webrest.PublicApi,
		},
		{
			URI:      "enrollments",
			Method:   http.MethodGet,
			Function: c.GetAll,
			Prefix:   webrest.PublicApi,
		},
		{
			URI:      "enrollments",
			Method:   http.MethodDelete,
			Function: c.Delete,
			Prefix:   webrest.PublicApi,
		},
	}
}

func (c *EnrollmentController) GetAll(w http.ResponseWriter, r *http.Request) {
	studentName := strings.ToLower(r.URL.Query().Get("studentName"))
	courseName := strings.ToLower(r.URL.Query().Get("courseName"))

	list, err := c.Usecase.GetAll(r.Context(), studentName, courseName)
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	webrest.JsonResponse(w, http.StatusOK, list)
}

func (c *EnrollmentController) Post(w http.ResponseWriter, r *http.Request) {
	body, err := webrest.DecodeBody[models.EnrollmentCreateDTO](r)
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusUnprocessableEntity, err)
		return
	}

	if err = c.Usecase.Create(r.Context(), body); err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *EnrollmentController) Delete(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	studentId, err := uuid.Parse(params.Get("student_id"))
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusBadRequest, err)
		return
	}

	courseId, err := uuid.Parse(params.Get("course_id"))
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusBadRequest, err)
		return
	}

	if err = c.Usecase.Delete(r.Context(), studentId, courseId); err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
