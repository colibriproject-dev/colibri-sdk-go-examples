package consumers

import (
	"context"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/usecases"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
)

type FinantialInstallmentConsumer struct {
	queueName string
	Usecase   usecases.IEnrollmentUsecases
}

func NewFinantialInstallmentConsumer() messaging.QueueConsumer {
	return &FinantialInstallmentConsumer{
		queueName: "FINANCIAL_INSTALLMENT_SCHOOL",
		Usecase:   usecases.NewEnrollmentUsecases(),
	}
}

func (c *FinantialInstallmentConsumer) Consume(ctx context.Context, providerMessage *messaging.ProviderMessage) error {
	var model models.Account
	if err := providerMessage.DecodeMessage(&model); err != nil {
		return err
	}

	if err := c.Usecase.UpdateStatus(ctx, model.StudentID, model.CourseID, model.Status); err != nil {
		return err
	}

	return nil
}

func (c *FinantialInstallmentConsumer) QueueName() string {
	return c.queueName
}
