package controllers

import (
	"testing"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
)

func TestMain(m *testing.M) {
	test.InitializeBaseTest()

	m.Run()
}
