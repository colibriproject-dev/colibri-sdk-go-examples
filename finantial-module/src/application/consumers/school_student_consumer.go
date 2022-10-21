package consumers

import (
	"context"
	"finantial-module/src/domain/models"
	"finantial-module/src/domain/usecases"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
)

type SchoolStudentConsumer interface {
	Delete(ctx context.Context, message *messaging.ProviderMessage) error
}

type SchoolStudentQueueConsumer struct {
	QueueName string
	Usecase   usecases.AccountUsecases
}

func NewSchoolStudentConsumer() {
	consumer := SchoolStudentQueueConsumer{
		QueueName: "SCHOOL_STUDENT_FINANCIAL",
		Usecase:   usecases.NewAccountUsecase(),
	}

	messaging.AddConsumer(messaging.NewConsumer(consumer.QueueName, consumer.Delete))
}

func (p *SchoolStudentQueueConsumer) Delete(ctx context.Context, message *messaging.ProviderMessage) error {
	var model models.Student
	if err := message.DecodeMessage(&model); err != nil {
		return err
	}

	if message.Action != "DELETE_STUDENT" {
		if err := p.Usecase.DeleteByStudent(ctx, model.ID); err != nil {
			return err
		}
	}

	return nil
}
