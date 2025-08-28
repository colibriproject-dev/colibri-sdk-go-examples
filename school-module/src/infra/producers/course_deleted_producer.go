//go:generate mockgen -source course_deleted_producer.go -destination mock/course_deleted_producer_mock.go -package producersmock
package producers

import (
	"context"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/messaging"
)

type ICourseDeletedProducer interface {
	Send(ctx context.Context, model *models.CourseDelete) error
}

type CourseDeletedProducer struct {
	producer *messaging.Producer
}

func NewCourseDeletedProducer() *CourseDeletedProducer {
	return &CourseDeletedProducer{messaging.NewProducer("SCHOOL_COURSE")}
}

func (p *CourseDeletedProducer) Send(ctx context.Context, model *models.CourseDelete) error {
	return p.producer.Publish(ctx, "DELETE_COURSE", model)
}
