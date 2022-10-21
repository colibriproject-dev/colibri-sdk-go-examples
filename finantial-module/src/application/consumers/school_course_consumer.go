package consumers

import (
	"context"
	"finantial-module/src/domain/models"
	"finantial-module/src/domain/usecases"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
)

type SchoolCourseConsumer interface {
	Delete(ctx context.Context, message *messaging.ProviderMessage) error
}

type SchoolCourseQueueConsumer struct {
	QueueName string
	Usecase   usecases.AccountUsecases
}

func NewSchoolCourseQueueConsumer() {
	consumer := SchoolCourseQueueConsumer{
		QueueName: "SCHOOL_COURSE_FINANCIAL",
		Usecase:   usecases.NewAccountUsecase(),
	}

	messaging.AddConsumer(messaging.NewConsumer(consumer.QueueName, consumer.Delete))
}

func (p *SchoolCourseQueueConsumer) Delete(ctx context.Context, message *messaging.ProviderMessage) error {
	var model models.Course
	if err := message.DecodeMessage(&model); err != nil {
		return err
	}

	if message.Action != "DELETE_COURSE" {
		if err := p.Usecase.DeleteByCourse(ctx, model.ID); err != nil {
			return err
		}
	}

	return nil
}
