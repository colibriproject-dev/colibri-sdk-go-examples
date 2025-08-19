package models

import (
	"testing"
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAccount_ToEnrollmentUpdateStatus(t *testing.T) {
	studentID := uuid.New()
	courseID := uuid.New()
	status := enums.ADIMPLENTE

	account := &models.Account{
		ID:           uuid.New(),
		StudentID:    studentID,
		CourseID:     courseID,
		Installments: 12,
		Value:        1000.0,
		Status:       status,
		CreatedAt:    time.Now(),
	}

	result := account.ToEnrollmentUpdateStatus()

	assert.NotNil(t, result)
	assert.Equal(t, studentID, result.StudentID)
	assert.Equal(t, courseID, result.CourseID)
	assert.Equal(t, status, result.Status)
}

func TestAccount_ToEnrollmentUpdateStatus_NilAccount(t *testing.T) {
	var account *models.Account

	assert.Panics(t, func() {
		account.ToEnrollmentUpdateStatus()
	})
}
