package controllers

import (
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/test"
)

func TestMain(m *testing.M) {
	test.InitializeBaseTest()

	m.Run()
}
