//go:generate mockgen -source invoice_repository.go -destination mock/invoice_repository_mock.go -package repositoriesmock
package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/finantial-module/src/domain/models"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/google/uuid"
)

type InvoiceRepository interface {
	FindAll(ctx context.Context) ([]models.Invoice, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Invoice, error)
	Insert(ctx context.Context, invoice *models.Invoice) error
	BulkInsert(ctx context.Context, invoices []models.Invoice) error
	UpdatePaymentDate(ctx context.Context, invoice *models.Invoice) error
	FindAllOverdueInvoices(ctx context.Context) ([]models.OverdueInvoices, error)
	FindTotalOverdueInvoicesByAccount(ctx context.Context, id uuid.UUID) (*uint64, error)
}

type InvoiceDBRepository struct{}

func NewInvoiceDBRepository() *InvoiceDBRepository {
	return &InvoiceDBRepository{}
}

func (r *InvoiceDBRepository) FindAll(ctx context.Context) ([]models.Invoice, error) {
	const query = `
		SELECT
			i.id,
			a.id, a.student_id, a.course_id, a.installments, a.value, a.status, a.created_at,
			i.installment, i.due_date, i.value, i.created_at, i.paid_at
		FROM invoices i
		INNER JOIN accounts a ON i.account_id = a.id`

	return sqlDB.NewQuery[models.Invoice](ctx, query).Many()
}

func (r *InvoiceDBRepository) FindById(ctx context.Context, id uuid.UUID) (*models.Invoice, error) {
	const query = `
		SELECT
			i.id,
			a.id, a.student_id, a.course_id, a.installments, a.value, a.status, a.created_at,
			i.installment, i.due_date, i.value, i.created_at, i.paid_at
		FROM invoices i
		INNER JOIN accounts a ON i.account_id = a.id
		WHERE i.id = $1`

	return sqlDB.NewQuery[models.Invoice](ctx, query, id).One()
}

func (r *InvoiceDBRepository) Insert(ctx context.Context, invoice *models.Invoice) error {
	const query = `INSERT INTO invoices (id, account_id, installment, due_date, value, created_at) VALUES ($1, $2, $3, $4, $5, $6)`

	return sqlDB.NewStatement(ctx, query, invoice.ID, invoice.Account.ID, invoice.Installment, invoice.DueDate, invoice.Value, invoice.CreatedAt).Execute()
}

func (r *InvoiceDBRepository) BulkInsert(ctx context.Context, invoices []models.Invoice) error {
	const query = `INSERT INTO invoices (id, account_id, installment, due_date, value, created_at) VALUES %s`

	values := []string{}
	for _, invoice := range invoices {
		value := fmt.Sprintf("('%s', '%s', %d, '%v', %v, '%v')", invoice.ID, invoice.Account.ID, invoice.Installment, invoice.DueDate.Format(time.RFC3339Nano), invoice.Value, invoice.CreatedAt.Format(time.RFC3339Nano))
		values = append(values, value)
	}

	return sqlDB.NewStatement(ctx, fmt.Sprintf(query, strings.Join(values, ", "))).Execute()
}

func (r *InvoiceDBRepository) UpdatePaymentDate(ctx context.Context, invoice *models.Invoice) error {
	const query = `UPDATE invoices SET due_date=$2 WHERE id=$1`

	return sqlDB.NewStatement(ctx, query, invoice.ID, invoice.DueDate).Execute()
}

func (r *InvoiceDBRepository) FindAllOverdueInvoices(ctx context.Context) ([]models.OverdueInvoices, error) {
	const query = `
		SELECT
			a.student_id,
			a.course_id,
			COUNT(i.id) AS TOTAL
		FROM invoices i
		INNER JOIN accounts a ON i.account_id = a.id AND a.status = 'ADIMPLENTE'
		WHERE i.due_date < CURRENT_DATE
		GROUP BY a.student_id, a.course_id
		HAVING COUNT(i.id) > 0`

	return sqlDB.NewQuery[models.OverdueInvoices](ctx, query).Many()
}

func (r *InvoiceDBRepository) FindTotalOverdueInvoicesByAccount(ctx context.Context, id uuid.UUID) (*uint64, error) {
	const query = `
		SELECT
			COUNT(i.id) AS TOTAL
		FROM invoices i
		INNER JOIN accounts a ON i.account_id = a.id AND a.id = $1
		WHERE i.due_date < CURRENT_DATE`

	return sqlDB.NewQuery[uint64](ctx, query).One()
}
