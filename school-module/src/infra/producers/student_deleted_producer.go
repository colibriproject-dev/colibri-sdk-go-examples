//go:generate mockgen -source student_deleted_producer.go -destination mock/student_deleted_producer_mock.go -package producersmock
package producers

import (
	"context"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/messaging"
)

type IStudentDeletedProducer interface {
	Send(ctx context.Context, model *models.StudentDelete) error
}

type StudentDeletedProducer struct {
	producer *messaging.Producer
}

func NewStudentDeletedProducer() *StudentDeletedProducer {
	return &StudentDeletedProducer{messaging.NewProducer("SCHOOL_STUDENT")}
}

func (p *StudentDeletedProducer) Send(ctx context.Context, model *models.StudentDelete) error {
	return p.producer.Publish(ctx, "DELETE_STUDENT", model)
}
