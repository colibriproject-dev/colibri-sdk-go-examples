//go:generate mockgen -source enrollment_deleted_producer.go -destination mock/enrollment_deleted_producer_mock.go -package mock
package producers

import (
	"context"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
)

type IEnrollmentDeletedProducer interface {
	Delete(ctx context.Context, model *models.Enrollment) error
}

type EnrollmentDeletedProducer struct {
	producer *messaging.Producer
}

func NewEnrollmentDeletedProducer() *EnrollmentDeletedProducer {
	return &EnrollmentDeletedProducer{messaging.NewProducer("SCHOOL_ENROLLMENT_DELETED")}
}

func (p *EnrollmentDeletedProducer) Delete(ctx context.Context, model *models.Enrollment) error {
	return p.producer.Publish(ctx, "DELETE_ENROLLMENT", model)
}
