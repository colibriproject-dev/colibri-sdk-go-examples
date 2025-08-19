package repositories

import (
	"context"
	"os"
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/test"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/database/sqlDB"
)

var (
	ctx      = context.Background()
	pc       *test.PostgresContainer
	basePath string
)

func TestMain(m *testing.M) {
	os.Setenv("MIGRATION_SOURCE_URL", "../../../migrations")
	test.InitializeSqlDBTest()
	basePath = test.MountAbsolutPath("../../../development-environment/database/tests-dataset/")
	pc = test.UsePostgresContainer(context.Background())
	sqlDB.Initialize()

	m.Run()
}
