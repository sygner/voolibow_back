package models

import "time"

type Transaction struct {
	TxId             string
	RoomId           string
	CurrencyId       string
	CurrencyFullName string
	FromAddress      string
	ToAddress        string
	TransactionAt    time.Time
}
