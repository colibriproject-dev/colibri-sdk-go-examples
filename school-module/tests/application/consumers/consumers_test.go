package consumers

import (
	"context"
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/logging"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/test"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/validator"
)

var (
	ctx = context.Background()
)

func TestMain(m *testing.M) {
	logging.Initialize()
	test.InitializeBaseTest()

	validator.RegisterCustomValidation("oneOfEnrollmentStatus", enums.EnrollmentStatusValidator)

	m.Run()
}
