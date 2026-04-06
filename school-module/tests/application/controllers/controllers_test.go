package controllers

import (
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/logging"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/test"
)

func TestMain(m *testing.M) {
	logging.Initialize()
	test.InitializeBaseTest()

	m.Run()
}
