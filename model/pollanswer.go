package model

type PollAnswer struct {
	ID     int
	Text   string
	Votes  int
	Rate   float32
	PollID int
}
