package consumers

import (
	"context"
	"school-module/src/domain/models"
	"school-module/src/domain/usecases"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
)

type IFinantialInstallmentConsumer interface {
	Update(ctx context.Context, message *messaging.ProviderMessage) error
}

type FinantialInstallmentConsumer struct {
	QueueName string
	Usecase   usecases.IEnrollmentUsecase
}

func NewFinantialInstallmentConsumer() {
	consumer := FinantialInstallmentConsumer{
		QueueName: "FINANCIAL_INSTALLMENT_SCHOOL",
		Usecase:   usecases.NewEnrollmentUsecase(),
	}

	messaging.AddConsumer(messaging.NewConsumer(consumer.QueueName, consumer.Update))
}

func (c *FinantialInstallmentConsumer) Update(ctx context.Context, message *messaging.ProviderMessage) error {
	var model models.Account
	if err := message.DecodeMessage(&model); err != nil {
		return err
	}

	if err := c.Usecase.UpdateStatus(ctx, model.StudentID, model.CourseID, model.Status); err != nil {
		return err
	}

	return nil
}
