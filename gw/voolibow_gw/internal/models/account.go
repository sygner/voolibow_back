package models

type Session struct {
	Agent     string `json:"agent"`
	Ip        string `json:"ip"`
	Status    string `json:"status"`
	SessionId int32  `json:"session_id"`
	CreatedAt string `json:"created_at"`
}

type LoginHistory struct {
	Id           string `json:"id"`
	UserId       int32  `json:"user_id"`
	UserRole     string `json:"user_role"`
	Section      string `json:"section"`
	Ip           string `json:"ip"`
	Agent        string `json:"agent"`
	Logged_in_at string `json:"logged_in_at"`
}
