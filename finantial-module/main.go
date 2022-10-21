package main

import (
	"finantial-module/src/application/consumers"
	"finantial-module/src/application/controllers"

	"github.com/colibri-project-io/colibri-sdk-go"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/webrest"
)

func init() {
	colibri.InitializeApp()
	messaging.Initialize()
	sqlDB.Initialize()
}

func main() {
	consumers.NewSchoolCourseQueueConsumer()
	consumers.NewSchoolEnrollmentQueueConsumer()
	consumers.NewSchoolStudentConsumer()
	messaging.Initialize()

	controllers.NewAccountRestController()
	controllers.NewInvoiceRestController()
	controllers.NewScheduledRestController()
	webrest.ListenAndServe()
}
