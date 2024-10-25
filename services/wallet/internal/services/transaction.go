package services

import (
	"voolibow_wallet/internal/models"
	"voolibow_wallet/internal/repository"
	"voolibow_wallet/internal/types"
)

type (
	TransactionService interface{}
	transactionService struct {
		repository repository.WalletRepository
	}
)

func NewTransactionService(repository repository.WalletRepository) TransactionService {
	return &transactionService{
		repository: repository,
	}
}

func (c *transactionService) GetTransactionByRoomId(roomId string) ([]models.Transaction, *types.Error) {
	return c.repository.GetTransactionsByRoomId(roomId)
}
