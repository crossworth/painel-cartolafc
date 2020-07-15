package model

import (
	"time"
)

type Topic struct {
	ID        int
	Title     string
	IsClosed  bool
	IsFixed   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy int
	UpdatedBy int
	Deleted   bool
}
