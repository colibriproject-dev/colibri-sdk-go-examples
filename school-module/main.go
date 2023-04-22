package main

import (
	"github.com/colibri-project-io/colibri-sdk-go"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/application/consumers"
	"github.com/colibri-project-io/colibri-sdk-go-examples/school-module/src/application/controllers"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/cacheDB"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/sqlDB"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/messaging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/storage"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/restserver"
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
	messaging.NewConsumer(consumers.NewFinantialInstallmentConsumer())

	restserver.AddRoutes(controllers.NewCourseController().Routes())
	restserver.AddRoutes(controllers.NewStudentController().Routes())
	restserver.AddRoutes(controllers.NewEnrollmentController().Routes())
	restserver.ListenAndServe()
}
