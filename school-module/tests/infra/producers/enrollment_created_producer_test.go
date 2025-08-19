package producers

import (
	"testing"
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/producers"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/messaging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEnrollmentCreatedProducer_Send(t *testing.T) {
	const testQueue string = "SCHOOL_ENROLLMENT_CREATED_TOPIC_TEST"

	t.Run("should send message", func(t *testing.T) {
		expected := &models.EnrollmentCreated{
			Student: models.EnrollmentCreatedStudent{
				ID: uuid.New(),
			},
			Course: models.EnrollmentCreatedCourse{
				ID: uuid.New(),
			},
			Installments: uint8(1),
			Status:       enums.ADIMPLENTE,
			CreatedAt:    time.Now().UTC(),
		}

		producerFn := func() error {
			return producers.NewEnrollmentCreatedProducer().Send(ctx, expected)
		}
		resp, err := messaging.NewTestProducer[models.EnrollmentCreated](producerFn, testQueue, 10).Execute()

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.EqualValues(t, expected, resp)
	})
}
