package producers

import (
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/infra/producers"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStudentDeletedProducer(t *testing.T) {
	const testTopic string = "SCHOOL_STUDENT_DELETED_TOPIC_TEST"

	t.Run("should return new student deleted producer", func(t *testing.T) {
		producer := producers.NewStudentDeletedProducer()

		assert.NotNil(t, producer)
	})

	t.Run("should send message", func(t *testing.T) {
		expected := &models.Student{
			ID:        uuid.New(),
			Name:      "New Student",
			Email:     "new.student@email.com",
			Birthday:  time.Now().UTC().Add(-1 * time.Hour),
			CreatedAt: time.Now().UTC(),
		}

		producerFn := func() error {
			return producers.NewStudentDeletedProducer().Delete(ctx, expected)
		}
		resp, err := messaging.NewTestProducer[models.Student](producerFn, testTopic, 10).Execute()

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expected, resp)
	})
}
