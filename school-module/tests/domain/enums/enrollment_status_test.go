package enums

import (
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestEnrollmentStatusValidator(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("oneOfEnrollmentStatus", enums.EnrollmentStatusValidator)

	tests := []struct {
		input    string
		expected bool
	}{
		{string(enums.ADIMPLENTE), true},
		{string(enums.INADIMPLENTE), true},
		{"INVALID_VALUE", false},
	}

	for _, test := range tests {
		err := validate.Var(test.input, "oneOfEnrollmentStatus")
		if test.expected {
			assert.NoError(t, err, "Expected no error for input: %s", test.input)
		} else {
			assert.Error(t, err, "Expected error for input: %s", test.input)
		}
	}
}
