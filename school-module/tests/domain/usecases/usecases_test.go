package usecases

import (
	"context"
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/logging"
)

var (
	ctx       = context.Background()
	exists    = true
	notExists = false
)

func TestMain(m *testing.M) {
	logging.Initialize()
}
