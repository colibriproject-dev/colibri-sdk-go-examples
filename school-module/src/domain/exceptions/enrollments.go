package exceptions

const (
	// Business exceptions
	ErrEnrollmentNotFound      string = "errEnrollmentNotFound"
	ErrEnrollmentAlreadyExists string = "errEnrollmentAlreadyExists"

	// Infra exceptions
	ErrOnFindAllPaginatedEnrollment             string = "ErrOnFindAllPaginatedEnrollment"
	ErrOnFindEnrollmentById                     string = "errOnFindEnrollmentById"
	ErrOnExistsEnrollmentById                   string = "errOnExistsEnrollmentById"
	ErrOnExistsEnrollmentByStudentIdAndCourseId string = "errOnExistsEnrollmentByStudentIdAndCourseId"
	ErrOnInsertEnrollment                       string = "errOnInsertEnrollment"
	ErrOnUpdateEnrollment                       string = "errOnUpdateEnrollment"
	ErrOnDeleteEnrollment                       string = "errOnDeleteEnrollment"
	ErrOnUpdateEnrollmentStatus                 string = "errOnUpdateEnrollmentStatus"
)
