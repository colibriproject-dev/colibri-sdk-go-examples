package producers

import (
	"context"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/producers"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCourseProducer(t *testing.T) {
	basePath := test.MountAbsolutPath("../../../development-environment/localstack/")
	test.InitializeTestLocalstack(basePath)
	messaging.Initialize()

	t.Run("should send message ok", func(t *testing.T) {
		event := models.Course{
			ID:        uuid.New(),
			Name:      "New Course",
			Value:     1000,
			CreatedAt: time.Now(),
		}
		fn := func() error {
			producers.NewCourseProducer().Delete(context.Background(), &event)
			return nil
		}

		resp, err := messaging.NewTestProducer[models.Course](fn, "SCHOOL_COURSE_OK_TEST", 5).Execute()

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}
