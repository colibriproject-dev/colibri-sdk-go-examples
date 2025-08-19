package models

import (
	"time"

	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/types"
	"github.com/google/uuid"
)

type Student struct {
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Birthday  types.IsoDate `json:"birthday"`
	CreatedAt time.Time     `json:"createdAt"`
}
