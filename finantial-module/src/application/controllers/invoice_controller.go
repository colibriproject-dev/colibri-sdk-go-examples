package controllers

import (
	"finantial-module/src/domain/models"
	"finantial-module/src/domain/usecases"
	"net/http"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/webrest"
	"github.com/google/uuid"
)

type InvoiceController interface {
	Routes() []webrest.Route
	GetAll(w http.ResponseWriter, r *http.Request)
	PatchPaymentDate(w http.ResponseWriter, r *http.Request)
}

type InvoiceRestController struct {
	Usecase usecases.InvoiceUsecases
}

func NewInvoiceRestController() {
	controller := &InvoiceRestController{
		Usecase: usecases.NewInvoiceUsecase(),
	}

	webrest.AddRoutes(controller.Routes())
}

func (p *InvoiceRestController) Routes() []webrest.Route {
	return []webrest.Route{
		{
			URI:      "invoices",
			Method:   http.MethodGet,
			Function: p.GetAll,
			Prefix:   webrest.PublicApi,
		},
		{
			URI:      "invoices/{id}/patch-payment-date",
			Method:   http.MethodPatch,
			Function: p.PatchPaymentDate,
			Prefix:   webrest.PublicApi,
		},
	}
}

func (p *InvoiceRestController) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := p.Usecase.GetAll(r.Context())
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	webrest.JsonResponse(w, http.StatusOK, list)
}

func (p *InvoiceRestController) PatchPaymentDate(w http.ResponseWriter, r *http.Request) {
	paramId, err := uuid.Parse(webrest.GetPathParam(r, "id"))
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusBadRequest, err)
		return
	}

	body, err := webrest.DecodeBody[models.Invoice](r)
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusUnprocessableEntity, err)
		return
	}

	if err = p.Usecase.UpdatePaymentDate(r.Context(), paramId, body.PaidAt.Time); err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
