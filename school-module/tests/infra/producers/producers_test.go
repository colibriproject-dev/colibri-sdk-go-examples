package producers

import (
	"context"
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/test"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/messaging"
)

var (
	ctx = context.Background()
)

func TestMain(m *testing.M) {
	test.InitializeTestLocalstack(test.MountAbsolutPath("../../../development-environment/localstack"))
	messaging.Initialize()

	m.Run()
}
