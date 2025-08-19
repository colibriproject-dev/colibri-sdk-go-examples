package producers

import (
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/producers"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/messaging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStudentDeletedProducer_Send(t *testing.T) {
	const testQueue string = "SCHOOL_STUDENT_DELETED_TOPIC_TEST"

	t.Run("should send message", func(t *testing.T) {
		expected := &models.StudentDelete{ID: uuid.New()}

		producerFn := func() error {
			return producers.NewStudentDeletedProducer().Send(ctx, expected)
		}
		resp, err := messaging.NewTestProducer[models.StudentDelete](producerFn, testQueue, 10).Execute()

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.EqualValues(t, expected, resp)
	})
}
