//go:generate mockgen -source enrollment_created_producer.go -destination mock/enrollment_created_producer_mock.go -package mock
package producers

import (
	"context"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
)

type IEnrollmentCreatedProducer interface {
	Create(ctx context.Context, model *models.Enrollment) error
}

type EnrollmentCreatedProducer struct {
	producer *messaging.Producer
}

func NewEnrollmentCreatedProducer() *EnrollmentCreatedProducer {
	return &EnrollmentCreatedProducer{messaging.NewProducer("SCHOOL_ENROLLMENT_CREATED")}
}

func (p *EnrollmentCreatedProducer) Create(ctx context.Context, model *models.Enrollment) error {
	return p.producer.Publish(ctx, "CREATE_ENROLLMENT", model)
}
