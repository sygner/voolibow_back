package models

type PhoneNumberDTO struct {
	PhoneNumber string `json:"phone_number"`
}

type VerificationDTO struct {
	VerificationMethod uint8  `json:"verification_method"`
	Code               uint16 `json:"code"`
	Agent              string
	Ip                 string
}

type RenewTokenDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Ip           string
	Agent        string
}
