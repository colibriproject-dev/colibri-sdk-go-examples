package models

import (
	"github.com/google/uuid"
)

type OverdueInvoices struct {
	StudentID    uuid.UUID
	CourseID     uuid.UUID
	Installments uint8
}
