package consumers

import (
	"context"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/usecases"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/messaging"
)

type FinantialInstallmentConsumer struct {
	queueName                     string
	UpdateEnrollmentStatusUsecase usecases.IUpdateEnrollmentStatusUsecase
}

func NewFinantialInstallmentConsumer() messaging.QueueConsumer {
	return &FinantialInstallmentConsumer{
		queueName:                     "FINANCIAL_INSTALLMENT_SCHOOL",
		UpdateEnrollmentStatusUsecase: usecases.NewUpdateEnrollmentStatusUsecase(),
	}
}

func (c *FinantialInstallmentConsumer) Consume(ctx context.Context, providerMessage *messaging.ProviderMessage) error {
	var model models.Account
	if err := providerMessage.DecodeAndValidateMessage(&model); err != nil {
		return err
	}

	if err := c.UpdateEnrollmentStatusUsecase.Execute(ctx, model.ToEnrollmentUpdateStatus()); err != nil {
		return err
	}

	return nil
}

func (c *FinantialInstallmentConsumer) QueueName() string {
	return c.queueName
}
