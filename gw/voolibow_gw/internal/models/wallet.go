package models

type Wallet struct {
	UserId           int32
	CurrencyId       string
	CurrencyFullName string
	Address          string
	PublicKey        string
	PrivateKey       string
	LastBalance      float64
	LastBalanceAt    string
	CreatedAt        string
}

type Transaction struct {
	TxId             string
	RoomId           string
	CurrencyId       string
	CurrencyFullName string
	FromAddress      string
	ToAddress        string
	TransactionAt    string
}
