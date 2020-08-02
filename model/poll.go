package model

import (
	"time"
)

type Poll struct {
	ID       int
	Question string
	Votes    int
	Multiple bool
	EndDate  time.Time
	Closed   bool
	TopicID  int
}
