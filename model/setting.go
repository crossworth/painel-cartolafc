package model

type Setting struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

const (
	MembersRuleSettingName = "members_rule"
	HomePageSettingName    = "home_page"
)
