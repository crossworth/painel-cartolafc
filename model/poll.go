package model

import (
	"time"
)

type Poll struct {
	ID       int       `json:"id"`
	Question string    `json:"question"`
	Votes    int       `json:"votes"`
	Multiple bool      `json:"multiple"`
	EndDate  time.Time `json:"end_date"`
	Closed   bool      `json:"closed"`
	TopicID  int       `json:"topic_id"`
}
