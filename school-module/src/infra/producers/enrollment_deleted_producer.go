//go:generate mockgen -source enrollment_deleted_producer.go -destination mock/enrollment_deleted_producer_mock.go -package producersmock
package producers

import (
	"context"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/messaging"
)

type IEnrollmentDeletedProducer interface {
	Send(ctx context.Context, model *models.EnrollmentDelete) error
}

type EnrollmentDeletedProducer struct {
	producer *messaging.Producer
}

func NewEnrollmentDeletedProducer() *EnrollmentDeletedProducer {
	return &EnrollmentDeletedProducer{messaging.NewProducer("SCHOOL_ENROLLMENT")}
}

func (p *EnrollmentDeletedProducer) Send(ctx context.Context, model *models.EnrollmentDelete) error {
	return p.producer.Publish(ctx, "DELETE_ENROLLMENT", model)
}
