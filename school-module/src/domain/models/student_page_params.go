package models

import (
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/types"
)

type StudentPage *types.Page[Student]

type StudentPageParams struct {
	Page uint16 `form:"page" validate:"required"`
	Size uint16 `form:"pageSize" validate:"required"`
	Name string `form:"name"`
}
