package model

type ProfileNames struct {
	ProfileID  int    `json:"profile_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ScreenName string `json:"screen_name"`
	Photo      string `json:"photo"`
	Date       int    `json:"date"`
}
