package models

type CourseCreate struct {
	Name  string  `json:"name" validate:"required"`
	Value float64 `json:"value" validate:"required"`
}
