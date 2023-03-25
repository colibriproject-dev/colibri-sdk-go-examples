package consumers

import (
	"context"
	"testing"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
)

var (
	ctx = context.Background()
)

func TestMain(m *testing.M) {
	test.InitializeBaseTest()

	m.Run()
}
