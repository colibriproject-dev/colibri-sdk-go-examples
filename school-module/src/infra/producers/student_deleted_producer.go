//go:generate mockgen -source student_deleted_producer.go -destination mock/student_deleted_producer_mock.go -package mock
package producers

import (
	"context"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
)

type IStudentDeletedProducer interface {
	Delete(ctx context.Context, model *models.Student) error
}

type StudentDeletedProducer struct {
	producer *messaging.Producer
}

func NewStudentDeletedProducer() *StudentDeletedProducer {
	return &StudentDeletedProducer{messaging.NewProducer("SCHOOL_STUDENT_DELETED")}
}

func (p *StudentDeletedProducer) Delete(ctx context.Context, model *models.Student) error {
	return p.producer.Publish(ctx, "DELETE_STUDENT", model)
}
