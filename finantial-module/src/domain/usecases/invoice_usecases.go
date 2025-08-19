//go:generate mockgen -source invoice_usecases.go -destination mock/invoice_usecases_mock.go -package usecasesmock
package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/finantial-module/src/domain/enums"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/finantial-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/finantial-module/src/infra/producers"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/finantial-module/src/infra/repositories"
	"github.com/google/uuid"
)

type InvoiceUsecases interface {
	GetAll(ctx context.Context) ([]models.Invoice, error)
	Create(ctx context.Context, model *models.Account) error
	ProcessAllOverdueInvoices(ctx context.Context) error
	UpdatePaymentDate(ctx context.Context, id uuid.UUID, paymentDate time.Time) error
}

type InvoiceUsecase struct {
	InvoiceRepository repositories.InvoiceRepository
	AccountRepository repositories.AccountRepository
	AccountProducer   producers.AccountProducer
}

func NewInvoiceUsecase() *InvoiceUsecase {
	return &InvoiceUsecase{
		InvoiceRepository: repositories.NewInvoiceDBRepository(),
		AccountRepository: repositories.NewAccountDBRepository(),
		AccountProducer:   producers.NewAccountProducer(),
	}
}

func (u *InvoiceUsecase) GetAll(ctx context.Context) ([]models.Invoice, error) {
	return u.InvoiceRepository.FindAll(ctx)
}

func (u *InvoiceUsecase) Create(ctx context.Context, model *models.Account) error {
	invoices := []models.Invoice{}
	for installment := uint8(1); installment <= model.Installments; installment++ {
		invoice := models.Invoice{
			ID:          uuid.New(),
			Account:     *model,
			Installment: installment,
			DueDate:     model.CreatedAt.Add(time.Duration(installment) * (30 * (24 * time.Hour))),
			Value:       model.Value / float64(model.Installments),
			CreatedAt:   time.Now(),
		}
		invoices = append(invoices, invoice)
	}

	return u.InvoiceRepository.BulkInsert(ctx, invoices)
}

func (u *InvoiceUsecase) ProcessAllOverdueInvoices(ctx context.Context) error {
	result, err := u.InvoiceRepository.FindAllOverdueInvoices(ctx)
	if err != nil {
		return err
	}

	for _, invoice := range result {
		account := &models.Account{
			StudentID:    invoice.StudentID,
			CourseID:     invoice.CourseID,
			Installments: invoice.Installments,
			Status:       enums.INADIMPLENTE,
		}

		u.AccountRepository.UpdateStatus(ctx, account)
		u.AccountProducer.StatusUpdated(ctx, account)
	}

	return nil
}

func (u *InvoiceUsecase) UpdatePaymentDate(ctx context.Context, id uuid.UUID, paymentDate time.Time) error {
	invoice, err := u.InvoiceRepository.FindById(ctx, id)
	if err != nil {
		return err
	}

	if invoice.ID == uuid.Nil {
		return errors.New("parcela não encontrada")
	}

	invoice.PaidAt.Time = paymentDate
	if err := u.InvoiceRepository.UpdatePaymentDate(ctx, invoice); err != nil {
		return err
	}

	if invoice.Account.Status == enums.INADIMPLENTE {
		total, _ := u.InvoiceRepository.FindTotalOverdueInvoicesByAccount(ctx, invoice.Account.ID)

		if *total == 0 {
			invoice.Account.Status = enums.ADIMPLENTE
			u.AccountRepository.UpdateStatus(ctx, &invoice.Account)
			u.AccountProducer.StatusUpdated(ctx, &invoice.Account)
		}
	}

	return nil
}
