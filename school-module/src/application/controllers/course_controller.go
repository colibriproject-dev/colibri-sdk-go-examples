package controllers

import (
	"net/http"
	"school-module/src/domain/models"
	"school-module/src/domain/usecases"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/webrest"
	"github.com/google/uuid"
)

type ICourseController interface {
	Routes() []webrest.Route
	GetAll(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type CourseController struct {
	Usecase usecases.ICourseUsecase
}

func NewCourseController() {
	controller := &CourseController{
		Usecase: usecases.NewCourseUsecase(),
	}

	webrest.AddRoutes(controller.Routes())
}

func (c *CourseController) Routes() []webrest.Route {
	return []webrest.Route{
		{
			URI:      "courses",
			Method:   http.MethodPost,
			Function: c.Post,
			Prefix:   webrest.PublicApi,
		},
		{
			URI:      "courses",
			Method:   http.MethodGet,
			Function: c.GetAll,
			Prefix:   webrest.PublicApi,
		},
		{
			URI:      "courses/{id}",
			Method:   http.MethodGet,
			Function: c.GetById,
			Prefix:   webrest.PublicApi,
		},
		{
			URI:      "courses/{id}",
			Method:   http.MethodPut,
			Function: c.Put,
			Prefix:   webrest.PublicApi,
		},
		{
			URI:      "courses/{id}",
			Method:   http.MethodDelete,
			Function: c.Delete,
			Prefix:   webrest.PublicApi,
		},
	}
}

func (c *CourseController) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := c.Usecase.GetAll(r.Context())
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	webrest.JsonResponse(w, http.StatusOK, result)
}

func (c *CourseController) GetById(w http.ResponseWriter, r *http.Request) {
	paramId, err := uuid.Parse(webrest.GetPathParam(r, "id"))
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusBadRequest, err)
		return
	}

	result, err := c.Usecase.GetById(r.Context(), paramId)
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	if result == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	webrest.JsonResponse(w, http.StatusOK, result)
}

func (c *CourseController) Post(w http.ResponseWriter, r *http.Request) {
	body, err := webrest.DecodeBody[models.CourseCreateUpdateDTO](r)
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusUnprocessableEntity, err)
		return
	}

	id, err := c.Usecase.Create(r.Context(), body)
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	webrest.JsonResponse(w, http.StatusCreated, id)
}

func (c *CourseController) Put(w http.ResponseWriter, r *http.Request) {
	paramId, err := uuid.Parse(webrest.GetPathParam(r, "id"))
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusBadRequest, err)
		return
	}

	body, err := webrest.DecodeBody[models.CourseCreateUpdateDTO](r)
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusUnprocessableEntity, err)
		return
	}

	if err = c.Usecase.Update(r.Context(), paramId, body); err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *CourseController) Delete(w http.ResponseWriter, r *http.Request) {
	paramId, err := uuid.Parse(webrest.GetPathParam(r, "id"))
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusBadRequest, err)
		return
	}

	if err = c.Usecase.Delete(r.Context(), paramId); err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
