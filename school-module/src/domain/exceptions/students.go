package exceptions

const (
	// Business exceptions
	ErrStudentNotFound      string = "errStudentNotFound"
	ErrStudentAlreadyExists string = "errStudentAlreadyExists"

	// Infra exceptions
	ErrOnFindAllPaginatedStudents string = "ErrOnFindAllPaginatedStudents"
	ErrOnFindStudentById          string = "errOnFindStudentById"
	ErrOnFindStudentByEmail       string = "errOnFindStudentByEmail"
	ErrOnExistsStudentById        string = "errOnExistsStudentById"
	ErrOnExistsStudentByEmail     string = "errOnExistsStudentByEmail"
	ErrOnInsertStudent            string = "errOnInsertStudent"
	ErrOnUpdateStudent            string = "errOnUpdateStudent"
	ErrOnDeleteStudent            string = "errOnDeleteStudent"
	ErrOnUploadStudentDocument    string = "errOnUploadStudentDocument"
)
