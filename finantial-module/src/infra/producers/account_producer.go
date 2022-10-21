package producers

import (
	"context"
	"finantial-module/src/domain/models"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
)

const (
	action_UPDATE_ACCOUNT_STATUS = "UPDATE_ACCOUNT_STATUS"
	topic_FINANCIAL_INSTALLMENT  = "FINANCIAL_INSTALLMENT"
)

type AccountProducer interface {
	StatusUpdated(ctx context.Context, model *models.Account)
}

type AccountTopicProducer struct {
	producer *messaging.Producer
}

func NewAccountProducer() *AccountTopicProducer {
	return &AccountTopicProducer{messaging.NewProducer(topic_FINANCIAL_INSTALLMENT)}
}

func (p *AccountTopicProducer) StatusUpdated(ctx context.Context, model *models.Account) {
	p.producer.Publish(ctx, action_UPDATE_ACCOUNT_STATUS, model)
}
