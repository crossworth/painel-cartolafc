package model

import (
	"time"
)

type ProfileNames struct {
	ProfileID  int
	FirstName  string
	LastName   string
	ScreenName string
	Photo      string
	Date       time.Time
}
