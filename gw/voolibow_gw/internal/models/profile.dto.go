package models

type UpdateUsernameDTO struct {
	Username string `json:"username"`
}

type UpdateProfileDTO struct {
	Avatar          *string `json:"avatar,omitempty"`
	DisplayUsername *string `json:"display_username,omitempty"`
}
