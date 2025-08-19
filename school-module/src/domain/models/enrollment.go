package models

import (
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
)

type Enrollment struct {
	Student      Student                `json:"student"`
	Course       Course                 `json:"course"`
	Installments uint8                  `json:"installments"`
	Status       enums.EnrollmentStatus `json:"status"`
	CreatedAt    time.Time              `json:"createdAt"`
}
