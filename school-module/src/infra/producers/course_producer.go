//go:generate mockgen -source course_producer.go -destination mock/course_producer_mock.go -package mock
package producers

import (
	"context"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
)

const (
	action_DELETE_COURSE = "DELETE_COURSE"
	topic_COURSE         = "SCHOOL_COURSE"
)

type ICourseProducer interface {
	Delete(ctx context.Context, model *models.Course)
}

type CourseProducer struct {
	producer *messaging.Producer
}

func NewCourseProducer() *CourseProducer {
	return &CourseProducer{messaging.NewProducer(topic_COURSE)}
}

func (p *CourseProducer) Delete(ctx context.Context, model *models.Course) {
	p.producer.Publish(ctx, action_DELETE_COURSE, model)
}
