package models

import "time"

type UpdateLastBalanceDTO struct {
	UserId        int32
	CurrencyId    string
	LastBalance   float64
	LastBalanceAt time.Time
}
