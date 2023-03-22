//go:generate mockgen -source enrollment_producer.go -destination mock/enrollment_producer_mock.go -package mock
package producers

import (
	"context"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
)

const (
	action_CREATE_ENROLLMENT = "CREATE_ENROLLMENT"
	action_DELETE_ENROLLMENT = "DELETE_ENROLLMENT"
	topic_ENROLLMENT         = "SCHOOL_ENROLLMENT"
)

type IEnrollmentProducer interface {
	Create(ctx context.Context, model *models.Enrollment)
	Delete(ctx context.Context, model *models.Enrollment)
}

type EnrollmentProducer struct {
	producer *messaging.Producer
}

func NewEnrollmentProducer() *EnrollmentProducer {
	return &EnrollmentProducer{messaging.NewProducer(topic_ENROLLMENT)}
}

func (p *EnrollmentProducer) Create(ctx context.Context, model *models.Enrollment) {
	p.producer.Publish(ctx, action_CREATE_ENROLLMENT, model)
}

func (p *EnrollmentProducer) Delete(ctx context.Context, model *models.Enrollment) {
	p.producer.Publish(ctx, action_DELETE_ENROLLMENT, model)
}
