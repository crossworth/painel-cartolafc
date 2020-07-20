package model

type Comment struct {
	ID         int    `json:"id"`
	FromID     int    `json:"from_id"`
	Date       int    `json:"date"`
	Text       string `json:"text"`
	Likes      int    `json:"likes"`
	ReplyToUID int    `json:"reply_to_uid"`
	ReplyToCID int    `json:"reply_to_cid"`
	TopicID    int    `json:"topic_id"`
	ProfileID  int    `json:"-"`
}
