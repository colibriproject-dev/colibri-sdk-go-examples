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

func TestCourseDeletedProducer(t *testing.T) {
	const testTopic string = "SCHOOL_COURSE_DELETED_TOPIC_TEST"

	t.Run("should return new course deleted producer", func(t *testing.T) {
		producer := producers.NewCourseDeletedProducer()

		assert.NotNil(t, producer)
	})

	t.Run("should return error when message is invalid", func(t *testing.T) {
		producerFn := func() error {
			return producers.NewCourseDeletedProducer().Delete(ctx, &models.Course{})
		}
		resp, err := messaging.NewTestProducer[models.Course](producerFn, testTopic, 10).Execute()

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should send message", func(t *testing.T) {
		expected := &models.Course{
			ID:        uuid.New(),
			Name:      "course test mock",
			Value:     rand.Float64(),
			CreatedAt: time.Now().UTC(),
		}

		producerFn := func() error {
			return producers.NewCourseDeletedProducer().Delete(ctx, expected)
		}
		resp, err := messaging.NewTestProducer[models.Course](producerFn, testTopic, 10).Execute()

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expected, resp)
	})
}
