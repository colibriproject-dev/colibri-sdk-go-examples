package controllers

import (
	"net/http"
	"school-module/src/domain/models"
	"school-module/src/domain/usecases"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/webrest"
	"github.com/google/uuid"
)

type IStudentController interface {
	Routes() []webrest.Route
	GetAll(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	PostDocument(w http.ResponseWriter, r *http.Request)
}

type StudentController struct {
	Usecase usecases.IStudentUsecase
}

func NewStudentController() {
	controller := &StudentController{
		Usecase: usecases.NewStudentUsecase(),
	}

	webrest.AddRoutes(controller.Routes())
}

func (c *StudentController) Routes() []webrest.Route {
	return []webrest.Route{
		{
			URI:      "students",
			Method:   http.MethodPost,
			Function: c.Post,
			Prefix:   webrest.PublicApi,
		},
		{
			URI:      "students",
			Method:   http.MethodGet,
			Function: c.GetAll,
			Prefix:   webrest.PublicApi,
		},
		{
			URI:      "students/{id}",
			Method:   http.MethodGet,
			Function: c.GetById,
			Prefix:   webrest.PublicApi,
		},
		{
			URI:      "students/{id}",
			Method:   http.MethodPut,
			Function: c.Put,
			Prefix:   webrest.PublicApi,
		},
		{
			URI:      "students/{id}",
			Method:   http.MethodDelete,
			Function: c.Delete,
			Prefix:   webrest.PublicApi,
		},
		{
			URI:      "students/{id}/upload-document",
			Method:   http.MethodPost,
			Function: c.PostDocument,
			Prefix:   webrest.PublicApi,
		},
	}
}

func (c *StudentController) GetAll(w http.ResponseWriter, r *http.Request) {
	params, err := webrest.DecodeParams[models.StudentParams](r)
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	list, err := c.Usecase.GetAll(r.Context(), params)
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	webrest.JsonResponse(w, http.StatusOK, list)
}

func (c *StudentController) GetById(w http.ResponseWriter, r *http.Request) {
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

func (c *StudentController) Post(w http.ResponseWriter, r *http.Request) {
	body, err := webrest.DecodeBody[models.StudentCreateUpdateDTO](r)
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusUnprocessableEntity, err)
		return
	}

	if err = c.Usecase.Create(r.Context(), body); err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	webrest.JsonResponse(w, http.StatusCreated, body)
}

func (c *StudentController) Put(w http.ResponseWriter, r *http.Request) {
	paramId, err := uuid.Parse(webrest.GetPathParam(r, "id"))
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusBadRequest, err)
		return
	}

	body, err := webrest.DecodeBody[models.StudentCreateUpdateDTO](r)
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

func (c *StudentController) Delete(w http.ResponseWriter, r *http.Request) {
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

func (c *StudentController) PostDocument(w http.ResponseWriter, r *http.Request) {
	paramId, err := uuid.Parse(webrest.GetPathParam(r, "id"))
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusBadRequest, err)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusBadRequest, err)
		return
	}

	url, err := c.Usecase.UploadDocument(r.Context(), paramId, &file)
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	webrest.JsonResponse(w, http.StatusOK, url)
}
