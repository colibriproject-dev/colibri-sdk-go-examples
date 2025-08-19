package controllers

import (
	"net/http"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/finantial-module/src/domain/usecases"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/web/restserver"
)

type AccountController struct {
	Usecase usecases.AccountUsecases
}

func NewAccountController() *AccountController {
	return &AccountController{
		Usecase: usecases.NewAccountUsecase(),
	}
}

func (p *AccountController) Routes() []restserver.Route {
	return []restserver.Route{
		{
			URI:      "accounts",
			Method:   http.MethodGet,
			Function: p.GetAll,
			Prefix:   restserver.PublicApi,
		},
	}
}

// @Summary Get account list
// @Tags accounts
// @Accept json
// @Produce json
// @Success 200 {array} models.Course
// @Failure 500
// @Router /public/accounts [get]
func (p *AccountController) GetAll(ctx restserver.WebContext) {
	list, err := p.Usecase.GetAll(ctx.Context())
	if err != nil {
		ctx.ErrorResponse(http.StatusInternalServerError, err)
		return
	}

	ctx.JsonResponse(http.StatusOK, list)
}
