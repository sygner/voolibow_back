package repository

import (
	"database/sql"
	"fmt"
	"voolibow_wallet/internal/models"
	"voolibow_wallet/internal/types"
)

type (
	WalletRepository interface {
		CheckWalletExists(int32, string) (bool, *types.Error)
		GetWalletAddressByUserIdAndCurrencyId(int32, string) (string, *types.Error)
		GetWalletByUserIdAndCurrencyId(int32, string) (*models.Wallet, *types.Error)
		CreateWallet(*models.Wallet) *types.Error
		UpdateLastBalance(*models.UpdateLastBalanceDTO) *types.Error
		GetWalletByAddress(string) (*models.Wallet, *types.Error)
		CreateTransaction(*models.Transaction) *types.Error
		GetTransactionsByRoomId(string) ([]models.Transaction, *types.Error)
		GetWalletsByUserId(int32) ([]models.Wallet, *types.Error)
		GetTransactionsByAddress(string) ([]models.Transaction, *types.Error)
	}
	walletRepository struct {
		db *sql.DB
	}
)

func NewWalletRepository(db *sql.DB) WalletRepository {
	return &walletRepository{
		db: db,
	}
}

func (c *walletRepository) CheckWalletExists(userId int32, currencyId string) (bool, *types.Error) {
	var exists bool

	if err := c.db.QueryRow("SELECT EXISTS (SELECT user_id FROM wallets WHERE user_id = $1 AND currency_id = $2)", userId, currencyId).Scan(&exists); err != nil {
		if err != sql.ErrNoRows {
			return false, types.NewInternalError("failed to fetch")
		}
	}

	return exists, nil
}

func (c *walletRepository) GetWalletAddressByUserIdAndCurrencyId(userId int32, currencyId string) (string, *types.Error) {
	var address string

	if err := c.db.QueryRow("SELECT address FROM wallets WHERE user_id = $1 AND currency_id = $2", userId, currencyId).Scan(&address); err != nil {
		if err == sql.ErrNoRows {
			return "", types.NewNotFoundError("no wallet found")
		} else {
			return "", types.NewInternalError("failed to fetch")
		}
	}

	return address, nil
}

func (c *walletRepository) GetWalletByUserIdAndCurrencyId(userId int32, currencyId string) (*models.Wallet, *types.Error) {
	var data models.Wallet
	if err := c.db.QueryRow("SELECT user_id, currency_id, currency_full_name, address, public_key, private_key, last_balance, last_balance_at, created_at FROM wallets WHERE user_id = $1 AND currency_id = $2", userId, currencyId).Scan(&data.UserId, &data.CurrencyId, &data.CurrencyFullName, &data.Address, &data.PublicKey, &data.PrivateKey, &data.LastBalance, &data.LastBalanceAt, &data.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("no wallet found")
		} else {
			return nil, types.NewInternalError("failed to fetch")
		}
	}

	return &data, nil
}

func (c *walletRepository) GetWalletByAddress(address string) (*models.Wallet, *types.Error) {
	var data models.Wallet
	if err := c.db.QueryRow("SELECT user_id, currency_id, currency_full_name, address, public_key, private_key, last_balance, last_balance_at, created_at FROM wallets WHERE address = $1", address).Scan(&data.UserId, &data.CurrencyId, &data.CurrencyFullName, &data.Address, &data.PublicKey, &data.PrivateKey, &data.LastBalance, &data.LastBalanceAt, &data.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("no wallet found")
		} else {
			return nil, types.NewInternalError("failed to fetch")
		}
	}

	return &data, nil
}

func (c *walletRepository) CreateWallet(data *models.Wallet) *types.Error {
	tx, err := c.db.Begin()
	if err != nil {
		return types.NewInternalError("failed to add transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var exists bool
	if err := tx.QueryRow("SELECT EXISTS (SELECT user_id FROM wallets WHERE user_id = $1 AND currency_id = $2)", data.UserId, data.CurrencyId).Scan(&exists); err != nil {
		if err != sql.ErrNoRows {
			return types.NewInternalError("failed to fetch")
		}
	}
	if exists {
		return types.NewBadRequestError("you already have the wallet")
	}
	// fmt.Sprintln("INSERT INTO wallets (user_id, currency_id, currency_full_name, address, public_key, private_key, last_balance, last_balance_at, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW())", data.UserId, data.CurrencyId, data.CurrencyFullName, data.Address, data.PublicKey, data.PrivateKey, data.LastBalance, data.LastBalanceAt)
	row := tx.QueryRow("INSERT INTO wallets (user_id, currency_id, currency_full_name, address, public_key, private_key, last_balance, last_balance_at, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW())", data.UserId, data.CurrencyId, data.CurrencyFullName, data.Address, data.PublicKey, data.PrivateKey, data.LastBalance, data.LastBalanceAt)

	if row.Err() != nil {
		return types.NewInternalError("failed to add")
	}

	err = tx.Commit()
	if err != nil {
		return types.NewInternalError("failed to commit")
	}

	return nil
}

func (c *walletRepository) UpdateLastBalance(data *models.UpdateLastBalanceDTO) *types.Error {
	_, err := c.db.Exec("UPDATE wallets SET last_balance = $1, last_balance_at = $2 WHERE user_id = $3 AND currency_id = $4", data.LastBalance, data.LastBalanceAt, data.UserId, data.CurrencyId)
	if err != nil {
		fmt.Println(err)
		return types.NewInternalError("failed to update")
	}
	return nil
}

func (c *walletRepository) CreateTransaction(data *models.Transaction) *types.Error {
	_, err := c.db.Exec("INSERT INTO transactions(tx_id, room_id, currency_id, currency_full_name, from_address, to_address, transaction_at) VALUES ($1,$2,$3,$4,$5,$6,NOW())", data.TxId, data.RoomId, data.CurrencyId, data.CurrencyFullName, data.FromAddress, data.ToAddress)
	if err != nil {
		fmt.Println(err)
		return types.NewInternalError("failed to add transaction")
	}

	return nil
}

func (c *walletRepository) GetTransactionsByRoomId(roomId string) ([]models.Transaction, *types.Error) {
	rows, err := c.db.Query("SELECT (tx_id, room_id,currency_id, currency_full_name, from_address, to_address, transaction_at) FROM transactions WHERE room_id = $1", roomId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("no transaction found")
		}
		return nil, types.NewInternalError("failed to get")
	}

	defer rows.Close()
	transactions := make([]models.Transaction, 0)
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.TxId, &transaction.RoomId, &transaction.CurrencyId, &transaction.CurrencyFullName, &transaction.FromAddress, &transaction.ToAddress, &transaction.TransactionAt)
		if err != nil {
			return nil, types.NewInternalError("failed to fetch")
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (c *walletRepository) GetTransactionsByAddress(address string) ([]models.Transaction, *types.Error) {
	fmt.Println(address)
	rows, err := c.db.Query("SELECT tx_id, room_id, currency_id, currency_full_name, from_address, to_address, transaction_at FROM transactions WHERE from_address = $1", address)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("no transaction found")
		}
		return nil, types.NewInternalError("failed to get")
	}

	defer rows.Close()
	transactions := make([]models.Transaction, 0)
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.TxId, &transaction.RoomId, &transaction.CurrencyId, &transaction.CurrencyFullName, &transaction.FromAddress, &transaction.ToAddress, &transaction.TransactionAt)
		if err != nil {
			fmt.Println(err)
			return nil, types.NewInternalError("failed to fetch")
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (c *walletRepository) GetWalletsByUserId(userId int32) ([]models.Wallet, *types.Error) {
	rows, err := c.db.Query("SELECT user_id, currency_id, currency_full_name, address, public_key, private_key, last_balance, last_balance_at, created_at FROM wallets WHERE user_id = $1", userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("no wallet found")
		}
		return nil, types.NewInternalError("failed to get")
	}

	defer rows.Close()
	wallets := make([]models.Wallet, 0)
	for rows.Next() {
		var wallet models.Wallet
		err := rows.Scan(&wallet.UserId, &wallet.CurrencyId, &wallet.CurrencyFullName, &wallet.Address, &wallet.PublicKey, &wallet.PrivateKey, &wallet.LastBalance, &wallet.LastBalanceAt, &wallet.CreatedAt)
		if err != nil {
			return nil, types.NewInternalError("failed to fetch")
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}
