package model

import (
	"time"
)

type Comment struct {
	ID         int
	FromID     int
	Date       time.Time
	Text       string
	Likes      int
	ReplyToUID int
	ReplyToCID int
	TopicID    int
}
