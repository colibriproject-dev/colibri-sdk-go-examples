package producers

import (
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/producers"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/messaging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEnrollmentDeletedProducer_Send(t *testing.T) {
	const testQueue string = "SCHOOL_ENROLLMENT_DELETED_TOPIC_TEST"

	t.Run("should send message", func(t *testing.T) {
		expected := &models.EnrollmentDelete{StudentID: uuid.New(), CourseID: uuid.New()}

		producerFn := func() error {
			return producers.NewEnrollmentDeletedProducer().Send(ctx, expected)
		}
		resp, err := messaging.NewTestProducer[models.EnrollmentDelete](producerFn, testQueue, 10).Execute()

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.EqualValues(t, expected, resp)
	})
}
