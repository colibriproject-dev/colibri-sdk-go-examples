package exceptions

const (
	// Business exceptions
	ErrCourseNotFound      string = "errCourseNotFound"
	ErrCourseAlreadyExists string = "errCourseAlreadyExists"

	// Infra exceptions
	ErrOnFindAllCourses     string = "errOnFindAllCourses"
	ErrOnFindCourseById     string = "errOnFindCourseById"
	ErrOnFindCourseByName   string = "errOnFindCourseByName"
	ErrOnExistsCourseById   string = "errOnExistsCourseById"
	ErrOnExistsCourseByName string = "errOnExistsCourseByName"
	ErrOnInsertCourse       string = "errOnInsertCourse"
	ErrOnUpdateCourse       string = "errOnUpdateCourse"
	ErrOnDeleteCourse       string = "errOnDeleteCourse"
)
