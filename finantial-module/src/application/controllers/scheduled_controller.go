package controllers

import (
	"finantial-module/src/domain/usecases"
	"net/http"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/webrest"
)

type ScheduledController interface {
	Routes() []webrest.Route
	ProcessAllOverdueInvoices(w http.ResponseWriter, r *http.Request)
}

type ScheduledRestController struct {
	Usecase usecases.InvoiceUsecases
}

func NewScheduledRestController() {
	controller := &ScheduledRestController{
		Usecase: usecases.NewInvoiceUsecase(),
	}

	webrest.AddRoutes(controller.Routes())
}

func (p *ScheduledRestController) Routes() []webrest.Route {
	return []webrest.Route{
		{
			URI:      "scheduled",
			Method:   http.MethodPost,
			Function: p.ProcessAllOverdueInvoices,
			Prefix:   webrest.PublicApi,
		},
	}
}

func (p *ScheduledRestController) ProcessAllOverdueInvoices(w http.ResponseWriter, r *http.Request) {
	if err := p.Usecase.ProcessAllOverdueInvoices(r.Context()); err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
