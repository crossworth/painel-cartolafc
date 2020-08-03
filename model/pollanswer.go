package model

type PollAnswer struct {
	ID     int     `json:"id"`
	Text   string  `json:"text"`
	Votes  int     `json:"votes"`
	Rate   float32 `json:"rate"`
	PollID int     `json:"poll_id"`
}
