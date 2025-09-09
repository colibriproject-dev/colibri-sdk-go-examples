//go:generate mockgen -source get_course_by_id_usecase.go -destination mock/get_course_by_id_usecase_mock.go -package usecasesmock
package usecases

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/exceptions"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/repositories"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/logging"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/monitoring"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/web/restclient"
	"github.com/google/uuid"
)

const (
	errAnErrorOccurredInGetCourseByIdUsecaseMsg string = "an error occurred in GetCourseByIdUsecase"
)

type IGetCourseByIdUsecase interface {
	Execute(ctx context.Context, id uuid.UUID) (*models.Course, error)
}

type GetCourseByIdUsecase struct {
	CourseRepository      repositories.ICoursesRepository
	financialModuleClient *restclient.RestClient
}

func NewGetCourseByIdUsecase() *GetCourseByIdUsecase {
	client := restclient.NewRestClient(&restclient.RestClientConfig{
		Name:                "financial-module-client",
		BaseURL:             os.Getenv("FINANCIAL_MODULE_BASE_URL"),
		Timeout:             0,
		Retries:             0,
		RetrySleepInSeconds: 0,
	})
	return &GetCourseByIdUsecase{
		CourseRepository:      repositories.NewCoursesDBRepository(),
		financialModuleClient: client,
	}
}

func (u *GetCourseByIdUsecase) Execute(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	result, err := u.CourseRepository.FindById(ctx, id)
	if err != nil {
		logging.Error(ctx).
			Err(err).
			AddParam("step", "CourseRepository.FindById").
			AddParam("id", id).
			Msg(errAnErrorOccurredInGetCourseByIdUsecaseMsg)
		return nil, errors.New(exceptions.ErrOnFindCourseById)
	}

	u.methodA(ctx)

	logging.Info(ctx).
		AddParam("model", result).
		Msg("course found")

	response := restclient.Request[any, any]{
		Ctx:        ctx,
		Client:     u.financialModuleClient,
		HttpMethod: http.MethodGet,
		Path:       "/public/accounts",
	}.Call()

	if response.StatusCode() != http.StatusOK {
		return nil, response.Error()
	}

	if result == nil {
		return nil, errors.New(exceptions.ErrCourseNotFound)
	}

	logging.Info(ctx).
		AddParam("result", result).
		Msg("course found")
	return result, nil
}

func (u *GetCourseByIdUsecase) methodA(ctx context.Context) {
	seg := monitoring.StartTransactionSegment(ctx, "sub segment", map[string]string{"name": "seg attr1", "value": "seg attr2"})
	defer monitoring.EndTransaction(seg)
	time.Sleep(50 * time.Millisecond)
	u.methodB(ctx)
}

func (u *GetCourseByIdUsecase) methodB(ctx context.Context) {
	seg2 := monitoring.StartTransactionSegment(ctx, "seg2", map[string]string{"name": "seg2 attr1", "value": "seg2 attr2"})
	defer monitoring.EndTransactionSegment(seg2)
	time.Sleep(25 * time.Millisecond)
}
