package models

import "time"

type Wallet struct {
	UserId           int32
	CurrencyId       string
	CurrencyFullName string
	Address          string
	PublicKey        string
	PrivateKey       string
	LastBalance      float64
	LastBalanceAt    time.Time
	CreatedAt        time.Time
}
