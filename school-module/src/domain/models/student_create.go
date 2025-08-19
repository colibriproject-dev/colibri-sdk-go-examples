package models

import (
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/types"
)

type StudentCreate struct {
	Name     string        `json:"name" validate:"required"`
	Email    string        `json:"email" validate:"required"`
	Birthday types.IsoDate `json:"birthday" validate:"required"`
}
