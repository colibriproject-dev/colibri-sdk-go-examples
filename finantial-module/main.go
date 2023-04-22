package main

import (
	"finantial-module/src/application/consumers"
	"finantial-module/src/application/controllers"

	"github.com/colibri-project-io/colibri-sdk-go"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/restserver"
)

func init() {
	colibri.InitializeApp()
	messaging.Initialize()
	sqlDB.Initialize()
}

// @title colibri-sdk-go-examples/finantial-module
// @version 1.0
// @description Microservice responsible for finantial management
func main() {
	messaging.NewConsumer(consumers.NewSchoolCourseConsumer())
	messaging.NewConsumer(consumers.NewSchoolEnrollmentConsumer())
	messaging.NewConsumer(consumers.NewSchoolStudentConsumer())

	restserver.AddRoutes(controllers.NewAccountController().Routes())
	restserver.AddRoutes(controllers.NewInvoiceController().Routes())
	restserver.AddRoutes(controllers.NewScheduledController().Routes())
	restserver.ListenAndServe()
}
