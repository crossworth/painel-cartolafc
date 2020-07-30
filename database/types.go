package database

import (
	"github.com/crossworth/cartola-web-admin/model"
)

type PaginationTimestamps struct {
	Next int
	Prev int
}

type OrderByDirection string

const (
	OrderByASC  = OrderByDirection("ASC")
	OrderByDESC = OrderByDirection("DESC")
)

func (o OrderByDirection) Stringer() string {
	return string(o)
}

type Period string

const (
	PeriodAll   = Period("")
	PeriodMonth = Period("month")
	PeriodWeek  = Period("week")
)

func (p Period) Stringer() string {
	return string(p)
}

func (p Period) URLString() string {
	if p == PeriodMonth {
		return "last_month"
	}

	if p == PeriodWeek {
		return "last_week"
	}

	return "all"
}

type ProfileWithStats struct {
	model.Profile
	Topics   int `json:"topics"`
	Comments int `json:"comments"`
	Likes    int `json:"likes"`
	Position int `json:"position,omitempty"`
}
