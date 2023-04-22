package consumers

import (
	"context"
	"finantial-module/src/domain/models"
	"finantial-module/src/domain/usecases"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
)

type SchoolEnrollmentConsumer struct {
	queueName string
	Usecase   usecases.AccountUsecases
}

func NewSchoolEnrollmentConsumer() messaging.QueueConsumer {
	return &SchoolEnrollmentConsumer{
		queueName: "SCHOOL_ENROLLMENT_FINANCIAL",
		Usecase:   usecases.NewAccountUsecase(),
	}
}

func (p *SchoolEnrollmentConsumer) Consume(ctx context.Context, providerMessage *messaging.ProviderMessage) error {
	var model models.Enrollment
	if err := providerMessage.DecodeMessage(&model); err != nil {
		return err
	}

	if providerMessage.Action == "CREATE_ENROLLMENT" {
		if err := p.Usecase.Create(ctx, model.ToAccount()); err != nil {
			return err
		}
	} else if providerMessage.Action == "DELETE_ENROLLMENT" {
		if err := p.Usecase.DeleteByStudentAndCourse(ctx, model.Student.ID, model.Course.ID); err != nil {
			return err
		}
	}

	return nil
}

func (c *SchoolEnrollmentConsumer) QueueName() string {
	return c.queueName
}
