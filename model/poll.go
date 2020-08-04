package model

type Poll struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Votes    int    `json:"votes"`
	Multiple bool   `json:"multiple"`
	EndDate  int    `json:"end_date"`
	Closed   bool   `json:"closed"`
	TopicID  int    `json:"topic_id"`
}
