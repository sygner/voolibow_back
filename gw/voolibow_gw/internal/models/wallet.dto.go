package models

type RegisterWalletDTO struct {
	UserId           int32
	Currency         int32 `json:"currency_id"`
	CurrencyFullName string
}

type TransferETHDTO struct {
	UserId    int32
	ToAddress string  `json:"to_address"`
	Amount    float64 `json:"amount"`
	RoomId    string  `json:"room_id,omitempty"`
}

type CurrencyIdDTO struct {
	Currency int32 `json:"currency_id"`
}

type AddressDTO struct {
	Address string `json:"address"`
}
