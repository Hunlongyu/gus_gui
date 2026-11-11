package state

type Preferences struct {
	Language   string `json:"language,omitempty"`
	Appearance string `json:"appearance,omitempty"`
}

func (s Preferences) Normalized() Preferences {
	if s.Language != "en" {
		s.Language = "zh-CN"
	}
	if s.Appearance != "light" && s.Appearance != "dark" {
		s.Appearance = "system"
	}
	return s
}

type WindowSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}
