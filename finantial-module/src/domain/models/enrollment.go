package models

type Enrollment struct {
	Student      Student `json:"student"`
	Course       Course  `json:"course"`
	Installments uint8   `json:"installments"`
}

func (e *Enrollment) ToAccount() *Account {
	return &Account{
		StudentID:    e.Student.ID,
		CourseID:     e.Course.ID,
		Installments: e.Installments,
		Value:        e.Course.Value,
	}
}
