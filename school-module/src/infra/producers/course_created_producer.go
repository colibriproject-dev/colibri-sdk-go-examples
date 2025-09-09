//go:generate mockgen -source course_created_producer.go -destination mock/course_created_producer_mock.go -package producersmock
package producers

import (
	"context"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/messaging"
)

type ICourseCreatedProducer interface {
	Send(ctx context.Context, model *models.Course) error
}

type CourseCreatedProducer struct {
	producer *messaging.Producer
}

func NewCourseCreatedProducer() ICourseCreatedProducer {
	return &CourseCreatedProducer{messaging.NewProducer("SCHOOL_COURSE")}
}

func (p *CourseCreatedProducer) Send(ctx context.Context, model *models.Course) error {
	return p.producer.Publish(ctx, "CREATE_COURSE", model)
}
