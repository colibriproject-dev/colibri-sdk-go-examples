package enums

import (
	"slices"

	"github.com/go-playground/validator/v10"
)

type EnrollmentStatus string

const (
	ADIMPLENTE   EnrollmentStatus = "ADIMPLENTE"
	INADIMPLENTE EnrollmentStatus = "INADIMPLENTE"
)

var enrollmentStatusValues = []string{
	ADIMPLENTE.String(),
	INADIMPLENTE.String(),
}

func (obj EnrollmentStatus) String() string {
	return string(obj)
}

func EnrollmentStatusValidator(fl validator.FieldLevel) bool {
	return slices.Contains(enrollmentStatusValues, fl.Field().String())
}
