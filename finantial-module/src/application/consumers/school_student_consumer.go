package consumers

import (
	"context"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/finantial-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/finantial-module/src/domain/usecases"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/logging"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/messaging"
)

type SchoolStudentConsumer struct {
	queueName string
	Usecase   usecases.AccountUsecases
}

func NewSchoolStudentConsumer() messaging.QueueConsumer {
	return &SchoolStudentConsumer{
		queueName: "SCHOOL_STUDENT_FINANCIAL",
		Usecase:   usecases.NewAccountUsecase(),
	}
}

func (p *SchoolStudentConsumer) Consume(ctx context.Context, providerMessage *messaging.ProviderMessage) error {
	var model models.Student
	if err := providerMessage.DecodeMessage(&model); err != nil {
		return err
	}

	logging.Info(ctx).
		AddParam("studentID", model.ID).
		Msg("Student received")

	if providerMessage.Action != "DELETE_STUDENT" {
		if err := p.Usecase.DeleteByStudent(ctx, model.ID); err != nil {
			return err
		}
	}

	return nil
}

func (c *SchoolStudentConsumer) QueueName() string {
	return c.queueName
}
