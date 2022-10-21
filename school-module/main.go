package main

import (
	"school-module/src/application/consumers"
	"school-module/src/application/controllers"

	"github.com/colibri-project-io/colibri-sdk-go"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/cacheDB"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/storage"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/webrest"
)

func init() {
	colibri.InitializeApp()
	storage.Initialize()
	cacheDB.Initialize()
	sqlDB.Initialize()
}

func main() {
	consumers.NewFinantialInstallmentConsumer()
	messaging.Initialize()

	controllers.NewCourseController()
	controllers.NewStudentController()
	controllers.NewEnrollmentController()
	webrest.ListenAndServe()
}
