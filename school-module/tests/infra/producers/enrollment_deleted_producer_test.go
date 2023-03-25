package producers

import (
	"math/rand"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/producers"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEnrollmentDeletedProducer(t *testing.T) {
	const testTopic string = "SCHOOL_ENROLLMENT_DELETED_TOPIC_TEST"

	t.Run("should return new enrollment deleted producer", func(t *testing.T) {
		producer := producers.NewEnrollmentDeletedProducer()

		assert.NotNil(t, producer)
	})

	t.Run("should return error when message is invalid", func(t *testing.T) {
		producerFn := func() error {
			return producers.NewEnrollmentDeletedProducer().Delete(ctx, &models.Enrollment{})
		}
		resp, err := messaging.NewTestProducer[models.Enrollment](producerFn, testTopic, 10).Execute()

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should send message", func(t *testing.T) {
		expected := &models.Enrollment{
			Student: models.Student{
				ID:        uuid.New(),
				Name:      "New Student",
				Email:     "new.student@email.com",
				Birthday:  time.Now().UTC().Add(-1 * time.Hour),
				CreatedAt: time.Now().UTC(),
			},
			Course: models.Course{
				ID:        uuid.New(),
				Name:      "New Course",
				Value:     rand.Float64(),
				CreatedAt: time.Now().UTC(),
			},
			Installments: uint8(1),
			Status:       models.ADIMPLENTE,
			CreatedAt:    time.Now().UTC(),
		}

		producerFn := func() error {
			return producers.NewEnrollmentDeletedProducer().Delete(ctx, expected)
		}
		resp, err := messaging.NewTestProducer[models.Enrollment](producerFn, testTopic, 10).Execute()

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expected, resp)
	})
}
