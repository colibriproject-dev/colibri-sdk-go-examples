package consumers

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/application/consumers"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/domain/usecases/mock"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewFinantialInstallmentConsumer(t *testing.T) {
	t.Run("Should return new accountancy created consumer", func(t *testing.T) {
		result := consumers.NewFinantialInstallmentConsumer()
		assert.NotNil(t, result)
		assert.NotNil(t, result.QueueName())
	})
}

func TestAccountancyCreatedConsumer(t *testing.T) {
	ctx := context.Background()
	providerMessageMock := &messaging.ProviderMessage{
		Message: models.Account{
			ID:           uuid.New(),
			StudentID:    uuid.New(),
			CourseID:     uuid.New(),
			Installments: uint8(rand.Int()),
			Value:        rand.Float64(),
			Status:       models.ADIMPLENTE,
			CreatedAt:    time.Now(),
		},
	}

	controller := gomock.NewController(t)
	mockUsecase := mock.NewMockIEnrollmentUsecases(controller)
	consumer := consumers.FinantialInstallmentConsumer{Usecase: mockUsecase}
	defer controller.Finish()

	t.Run("Should return error when occurred error in DecodeMessage", func(t *testing.T) {
		err := consumer.Consume(ctx, &messaging.ProviderMessage{Message: ""})
		assert.Error(t, err)
	})

	t.Run("Should return error when occurred error in UpdateStatus", func(t *testing.T) {
		expected := errors.New("mock error in UpdateStatus")
		mockUsecase.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(expected)

		err := consumer.Consume(ctx, providerMessageMock)
		assert.Error(t, expected, err)
	})

	t.Run("Should consume message and update enrollment status", func(t *testing.T) {
		mockUsecase.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := consumer.Consume(ctx, providerMessageMock)
		assert.NoError(t, err)
	})
}
