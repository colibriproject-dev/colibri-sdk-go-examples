package consumers

import (
	"context"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/finantial-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/finantial-module/src/domain/usecases"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/messaging"
)

type SchoolCourseConsumer struct {
	queueName string
	Usecase   usecases.AccountUsecases
}

func NewSchoolCourseConsumer() messaging.QueueConsumer {
	return &SchoolCourseConsumer{
		queueName: "SCHOOL_COURSE_FINANCIAL",
		Usecase:   usecases.NewAccountUsecase(),
	}
}

func (p *SchoolCourseConsumer) Consume(ctx context.Context, providerMessage *messaging.ProviderMessage) error {
	var model models.Course
	if err := providerMessage.DecodeMessage(&model); err != nil {
		return err
	}

	if providerMessage.Action != "DELETE_COURSE" {
		if err := p.Usecase.DeleteByCourse(ctx, model.ID); err != nil {
			return err
		}
	}

	return nil
}

func (c *SchoolCourseConsumer) QueueName() string {
	return c.queueName
}
