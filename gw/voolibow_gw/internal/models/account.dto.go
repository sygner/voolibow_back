package models

type GetLoginHistoriesDTO struct {
	UserId     int32
	Pagination Pagination `json:"pagination"`
}
