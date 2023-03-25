package repositories

import (
	"context"
	"os"
	"testing"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/storage"
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
	pc = test.UsePostgresContainer()
	sqlDB.Initialize()

	test.InitializeTestLocalstack(test.MountAbsolutPath("../../../development-environment/localstack"))
	storage.Initialize()

	m.Run()
}
