package models

import (
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/finantial-module/src/domain/enums"
	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID           `json:"id"`
	StudentID    uuid.UUID           `json:"studentId"`
	CourseID     uuid.UUID           `json:"courseId"`
	Installments uint8               `json:"installments"`
	Value        float64             `json:"value"`
	Status       enums.AccountStatus `json:"status"`
	CreatedAt    time.Time           `json:"createdAt"`
}
