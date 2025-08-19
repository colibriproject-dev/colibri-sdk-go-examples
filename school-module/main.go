package main

import (
	"github.com/colibriproject-dev/colibri-sdk-go"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/application/consumers"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/application/controllers"
	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/domain/enums"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/validator"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/database/cacheDB"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/messaging"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/storage"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/web/restserver"
)

func init() {
	colibri.InitializeApp()
	storage.Initialize()
	cacheDB.Initialize()
	sqlDB.Initialize()
	messaging.Initialize()
}

// @title colibri-sdk-go-examples/school-module
// @version 1.0
// @description Microservice responsible for school management
func main() {
	registerCustomValidators()
	registerConsumers()
	registerRoutes()

	restserver.ListenAndServe()
}

func registerCustomValidators() {
	validator.RegisterCustomValidation("oneOfEnrollmentStatus", enums.EnrollmentStatusValidator)
}

func registerConsumers() {
	messaging.NewConsumer(consumers.NewFinantialInstallmentConsumer())
}

func registerRoutes() {
	restserver.AddRoutes(controllers.NewCoursesV1Controller().Routes())
	restserver.AddRoutes(controllers.NewStudentController().Routes())
	restserver.AddRoutes(controllers.NewEnrollmentsV1Controller().Routes())
}
