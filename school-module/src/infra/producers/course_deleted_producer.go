//go:generate mockgen -source course_deleted_producer.go -destination mock/course_deleted_producer_mock.go -package mock
package producers

import (
	"context"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
)

type ICourseDeletedProducer interface {
	Delete(ctx context.Context, model *models.Course) error
}

type CourseDeletedProducer struct {
	producer *messaging.Producer
}

func NewCourseDeletedProducer() *CourseDeletedProducer {
	return &CourseDeletedProducer{messaging.NewProducer("SCHOOL_COURSE_DELETED")}
}

func (p *CourseDeletedProducer) Delete(ctx context.Context, model *models.Course) error {
	return p.producer.Publish(ctx, "DELETE_COURSE", model)
}
