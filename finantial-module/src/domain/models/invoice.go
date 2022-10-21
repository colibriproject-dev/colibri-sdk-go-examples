package models

import (
	"fmt"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/types"
	"github.com/google/uuid"
)

type Invoice struct {
	ID          uuid.UUID      `json:"id"`
	Account     Account        `json:"account"`
	Installment uint8          `json:"installment"`
	DueDate     time.Time      `json:"dueDate"`
	Value       float64        `json:"value"`
	CreatedAt   time.Time      `json:"createdAt"`
	PaidAt      types.NullTime `json:"paidAt"`
}

func (i *Invoice) Prepare() error {
	if err := i.validate(); err != nil {
		return err
	}

	i.format()
	return nil
}

func (i *Invoice) validate() error {
	if i.Account.ID == uuid.Nil {
		return fmt.Errorf("campo %s é requerido", "conta")
	}

	if i.Installment == 0 {
		return fmt.Errorf("campo %s é requerido", "Parcela")
	}

	if i.DueDate.IsZero() {
		return fmt.Errorf("campo %s é requerido", "Data de vencimento")
	}

	if i.Value == 0.0 {
		return fmt.Errorf("campo %s é requerido", "Valor")
	}

	return nil
}

func (i *Invoice) format() {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}

	if i.CreatedAt.IsZero() {
		i.CreatedAt = time.Now()
	}
}
