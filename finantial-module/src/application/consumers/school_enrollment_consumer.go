package consumers

import (
	"context"
	"finantial-module/src/domain/models"
	"finantial-module/src/domain/usecases"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
)

type SchoolEnrollmentConsumer interface {
	CreateOrDelete(ctx context.Context, message *messaging.ProviderMessage) error
}

type SchoolEnrollmentQueueConsumer struct {
	QueueName string
	Usecase   usecases.AccountUsecases
}

func NewSchoolEnrollmentQueueConsumer() {
	consumer := SchoolEnrollmentQueueConsumer{
		QueueName: "SCHOOL_ENROLLMENT_FINANCIAL",
		Usecase:   usecases.NewAccountUsecase(),
	}

	messaging.AddConsumer(messaging.NewConsumer(consumer.QueueName, consumer.CreateOrDelete))
}

func (p *SchoolEnrollmentQueueConsumer) CreateOrDelete(ctx context.Context, message *messaging.ProviderMessage) error {
	var model models.Enrollment
	if err := message.DecodeMessage(&model); err != nil {
		return err
	}

	if message.Action == "CREATE_ENROLLMENT" {
		if err := p.Usecase.Create(ctx, model.ToAccount()); err != nil {
			return err
		}
	} else if message.Action == "DELETE_ENROLLMENT" {
		if err := p.Usecase.DeleteByStudentAndCourse(ctx, model.Student.ID, model.Course.ID); err != nil {
			return err
		}
	}

	return nil
}
