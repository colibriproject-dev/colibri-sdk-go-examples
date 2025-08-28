//go:generate mockgen -source enrollment_created_producer.go -destination mock/enrollment_created_producer_mock.go -package producersmock
package producers

import (
	"context"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/logging"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/messaging"
)

type IEnrollmentCreatedProducer interface {
	Send(ctx context.Context, model *models.EnrollmentCreated) error
}

type EnrollmentCreatedProducer struct {
	producer *messaging.Producer
}

func NewEnrollmentCreatedProducer() *EnrollmentCreatedProducer {
	return &EnrollmentCreatedProducer{messaging.NewProducer("SCHOOL_ENROLLMENT")}
}

func (p *EnrollmentCreatedProducer) Send(ctx context.Context, model *models.EnrollmentCreated) error {
	logging.Info(ctx).Msg("Sending enrollment created message")
	return p.producer.Publish(ctx, "CREATE_ENROLLMENT", model)
}
