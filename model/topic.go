package model

type Topic struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	IsClosed  bool   `json:"is_closed"`
	IsFixed   bool   `json:"is_fixed"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
	CreatedBy int    `json:"created_by"`
	UpdatedBy int    `json:"updated_by"`
	Deleted   bool   `json:"-"`
}

type TopicWithComments struct {
	Topic
	CommentsCount int `json:"comments_count"`
}

type TopicWithLikes struct {
	Topic
	LikesCount int `json:"likes_count"`
}
