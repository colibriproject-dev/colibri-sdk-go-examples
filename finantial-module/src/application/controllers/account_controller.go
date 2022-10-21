package controllers

import (
	"finantial-module/src/domain/usecases"
	"net/http"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/webrest"
)

type AccountController interface {
	Routes() []webrest.Route
	GetAll(w http.ResponseWriter, r *http.Request)
}

type AccountRestController struct {
	Usecase usecases.AccountUsecases
}

func NewAccountRestController() {
	controller := &AccountRestController{
		Usecase: usecases.NewAccountUsecase(),
	}

	webrest.AddRoutes(controller.Routes())
}

func (p *AccountRestController) Routes() []webrest.Route {
	return []webrest.Route{
		{
			URI:      "accounts",
			Method:   http.MethodGet,
			Function: p.GetAll,
			Prefix:   webrest.PublicApi,
		},
	}
}

func (p *AccountRestController) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := p.Usecase.GetAll(r.Context())
	if err != nil {
		webrest.ErrorResponse(r, w, http.StatusInternalServerError, err)
		return
	}

	webrest.JsonResponse(w, http.StatusOK, list)
}
