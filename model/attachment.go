package model

type Attachment struct {
	Content   string `json:"content"`
	CommentID int    `json:"comment_id"`
}
