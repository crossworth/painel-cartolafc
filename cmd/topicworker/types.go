package main

type Profile struct {
	ID         int
	FirstName  string
	LastName   string
	ScreenName string
	Photo      string
}

type Comment struct {
	ID          int
	FromID      int
	Date        int64
	Text        string
	Likes       int
	ReplyToUID  int
	ReplyToCID  int
	Attachments []string
}

type Poll struct {
	ID       int
	Question string
	Votes    int
	Answers  []PollAnswer
	Multiple bool
	EndDate  int64
	Closed   bool
}

type PollAnswer struct {
	ID    int
	Text  string
	Votes int
	Rate  float64
}

type Topic struct {
	ID        int
	Title     string
	IsClosed  bool
	IsFixed   bool
	CreatedAt int64
	UpdatedAt int64
	CreatedBy Profile
	UpdatedBy Profile
	Profiles  map[int]Profile
}
