package controllers

import (
	"net/http"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/finantial-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/finantial-module/src/domain/usecases"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/web/restserver"
	"github.com/google/uuid"
)

type InvoiceController struct {
	Usecase usecases.InvoiceUsecases
}

func NewInvoiceController() *InvoiceController {
	return &InvoiceController{
		Usecase: usecases.NewInvoiceUsecase(),
	}
}

func (p *InvoiceController) Routes() []restserver.Route {
	return []restserver.Route{
		{
			URI:      "invoices",
			Method:   http.MethodGet,
			Function: p.GetAll,
			Prefix:   restserver.PublicApi,
		},
		{
			URI:      "invoices/{id}/patch-payment-date",
			Method:   http.MethodPatch,
			Function: p.PatchPaymentDate,
			Prefix:   restserver.PublicApi,
		},
	}
}

// @Summary Get invoice list
// @Tags invoices
// @Accept json
// @Produce json
// @Success 200 {array} models.Invoice
// @Failure 500
// @Router /public/invoices [get]
func (p *InvoiceController) GetAll(ctx restserver.WebContext) {
	list, err := p.Usecase.GetAll(ctx.Context())
	if err != nil {
		ctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	ctx.JsonResponse(http.StatusOK, list)
}

// @Summary Update payment date
// @Tags invoices
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 422
// @Failure 500
// @Param id path string true "Invoice ID"
// @Param request body models.Invoice true "request body"
// @Router /public/invoices/{id}/patch-payment-date [patch]
func (p *InvoiceController) PatchPaymentDate(ctx restserver.WebContext) {
	paramId, err := uuid.Parse(ctx.PathParam("id"))
	if err != nil {
		ctx.ErrorResponse(http.StatusBadRequest, err)
		return
	}

	var body models.Invoice
	if err := ctx.DecodeBody(&body); err != nil {
		ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
		return
	}

	if err = p.Usecase.UpdatePaymentDate(ctx.Context(), paramId, body.PaidAt.Time); err != nil {
		ctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	ctx.EmptyResponse(http.StatusNoContent)
}
