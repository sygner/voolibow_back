package models

type TransferETHDTO struct {
	UserId    int32
	ToAddress string
	Amount    float64
	RoomId    string
}
