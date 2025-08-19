package controllers

import (
	"net/http"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/finantial-module/src/domain/usecases"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/web/restserver"
)

type ScheduledController struct {
	Usecase usecases.InvoiceUsecases
}

func NewScheduledController() *ScheduledController {
	return &ScheduledController{
		Usecase: usecases.NewInvoiceUsecase(),
	}
}

func (p *ScheduledController) Routes() []restserver.Route {
	return []restserver.Route{
		{
			URI:      "scheduled",
			Method:   http.MethodPost,
			Function: p.ProcessAllOverdueInvoices,
			Prefix:   restserver.PublicApi,
		},
	}
}

// @Summary Run scheduled routine
// @Tags scheduled
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /public/scheduled [post]
func (p *ScheduledController) ProcessAllOverdueInvoices(ctx restserver.WebContext) {
	if err := p.Usecase.ProcessAllOverdueInvoices(ctx.Context()); err != nil {
		ctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	ctx.EmptyResponse(http.StatusOK)
}
