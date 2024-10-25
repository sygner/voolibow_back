package models

type Profile struct {
	UserId          int32  `json:"user_id,omitempty"`
	DisplaySid      string `json:"display_sid"`
	DisplayUsername string `json:"display_username"`
	Username        string `json:"username"`
	Avatar          string `json:"avatar"`
}
