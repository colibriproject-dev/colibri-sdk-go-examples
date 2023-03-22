//go:generate mockgen -source student_producer.go -destination mock/student_producer_mock.go -package mock
package producers

import (
	"context"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
)

const (
	action_DELETE_STUDENT = "DELETE_STUDENT"
	topic_STUDENT         = "SCHOOL_STUDENT"
)

type IStudentProducer interface {
	Delete(ctx context.Context, model *models.Student)
}

type StudentProducer struct {
	producer *messaging.Producer
}

func NewStudentProducer() *StudentProducer {
	return &StudentProducer{messaging.NewProducer(topic_STUDENT)}
}

func (p *StudentProducer) Delete(ctx context.Context, model *models.Student) {
	p.producer.Publish(ctx, action_DELETE_STUDENT, model)
}
