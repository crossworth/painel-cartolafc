package database

import (
	"strings"
	"time"

	"github.com/crossworth/painel-cartolafc/model"
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
	PeriodDay   = Period("day")
	PeriodWeek  = Period("week")
	PeriodMonth = Period("month")
)

func (p Period) Stringer() string {
	return string(p)
}

func PeriodFromString(periodStr string) Period {
	period := PeriodAll

	if strings.ToLower(periodStr) == "last_day" {
		period = PeriodDay
	}

	if strings.ToLower(periodStr) == "last_week" {
		period = PeriodWeek
	}

	if strings.ToLower(periodStr) == "last_month" {
		period = PeriodMonth
	}

	return period
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

type OrderBy string

const (
	OrderByCreatedAt = OrderBy("created_at")
	OrderByUpdatedAt = OrderBy("updated_at")
)

func (o OrderBy) Stringer() string {
	return string(o)
}

type ProfileWithStats struct {
	model.Profile
	Topics             int `json:"topics"`
	Comments           int `json:"comments"`
	Likes              int `json:"likes"`
	Position           int `json:"position,omitempty"`
	TopicsPlusComments int `json:"topics_plus_comments"`
}

type CommentWithProfileAndAttachment struct {
	model.Comment
	Profile     model.Profile      `json:"profile"`
	Attachments []model.Attachment `json:"attachments"`
}

type PollWithAnswers struct {
	model.Poll
	Answers []model.PollAnswer `json:"answers"`
}

type TopicWithPollAndCommentsCount struct {
	model.Topic
	Poll          *PollWithAnswers `json:"poll,omitempty"`
	CommentsCount int              `json:"comments_count"`
}

type SearchType string

const (
	SearchTypeTopic   = SearchType("topic")
	SearchTypeComment = SearchType("comment")
)

func (s SearchType) Stringer() string {
	return string(s)
}

type Search struct {
	Term           string     `json:"term"`
	Headline       string     `json:"headline"`
	Type           SearchType `json:"type"`
	Date           int        `json:"date"`
	TopicID        int        `json:"topic_id"`
	CommentID      int        `json:"comment_id"`
	CommentsCount  int        `json:"comments_count"`
	FromID         int        `json:"from_id"`
	FromName       string     `json:"from_name"`
	FromScreenName string     `json:"from_screen_name"`
	FromPhoto      string     `json:"from_photo"`
	LikesCount     int        `json:"likes_count"`
	Rank           float32    `json:"-"`
}

type TopicsWithStats struct {
	model.Topic
	Comments int `json:"comments"`
	Likes    int `json:"likes"`
	Position int `json:"position,omitempty"`
}

type QuotesByBot struct {
	TopicID     int    `json:"topic_id"`
	CommentID   int    `json:"comment_id"`
	TopicTitle  string `json:"topic_title"`
	DateComment int    `json:"date_comment"`
}

type GraphValue struct {
	Day   time.Time `json:"day"`
	Value int       `json:"value"`
}
