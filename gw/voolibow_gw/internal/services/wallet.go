package services

import (
	"context"
	"voolibow_gw/internal/models"
	wallet_pb "voolibow_gw/proto/api/wallet"
	"voolibow_gw/types"

	"google.golang.org/grpc"
)

type (
	WalletService interface {
		RegitserWallet(*models.RegisterWalletDTO) (string, *types.Error)
		GetBalance(int32, int32) (float64, *types.Error)
		TransferETH(*models.TransferETHDTO) *types.Error
		GetWalletsByUserId(int32) ([]models.Wallet, *types.Error)
		GetTransactionsByAddress(string) ([]models.Transaction, *types.Error)
	}
	walletService struct {
		walletClient wallet_pb.WalletServiceClient
	}
)

func NewWalletService(walletServiceConnection *grpc.ClientConn) WalletService {
	return &walletService{
		walletClient: wallet_pb.NewWalletServiceClient(walletServiceConnection),
	}
}

func (c *walletService) RegitserWallet(data *models.RegisterWalletDTO) (string, *types.Error) {
	res, err := c.walletClient.RegisterWallet(context.Background(), &wallet_pb.RegisterWalletRequest{
		UserId:           data.UserId,
		Currency:         wallet_pb.Currencies(data.Currency),
		CurrencyFullName: data.CurrencyFullName,
	})

	if err != nil {
		return "", types.ExtractGrpcError(err)
	}
	return res.Address, nil
}

func (c *walletService) GetBalance(userId int32, currencyId int32) (float64, *types.Error) {
	res, err := c.walletClient.GetBalance(context.Background(), &wallet_pb.GetBalanceRequest{
		UserId:   userId,
		Currency: wallet_pb.Currencies(currencyId),
	})
	if err != nil {
		return 0, types.ExtractGrpcError(err)
	}
	return res.Balance, nil
}

func (c *walletService) TransferETH(data *models.TransferETHDTO) *types.Error {
	_, err := c.walletClient.TransferETH(context.Background(), &wallet_pb.TransferETHRequest{
		UserId:    data.UserId,
		ToAddress: data.ToAddress,
		Amount:    data.Amount,
		RoomId:    data.RoomId,
	})
	if err != nil {
		return types.ExtractGrpcError(err)
	}
	return nil
}

func (c *walletService) GetWalletsByUserId(userId int32) ([]models.Wallet, *types.Error) {
	res, err := c.walletClient.GetWalletsByUserId(context.Background(), &wallet_pb.UserIdRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, types.ExtractGrpcError(err)
	}
	wallets := make([]models.Wallet, 0)
	for _, w := range res.Wallets {
		wallets = append(wallets, models.Wallet{
			UserId:           userId,
			CurrencyId:       w.CurrencyId,
			CurrencyFullName: w.CurrencyFullName,
			Address:          w.Address,
			PublicKey:        w.PublicKey,
			PrivateKey:       w.PrivateKey,
			LastBalance:      w.LastBalance,
			LastBalanceAt:    w.LastBalanceAt,
			CreatedAt:        w.CreatedAt,
		})
	}
	return wallets, nil
}

func (c *walletService) GetTransactionsByAddress(address string) ([]models.Transaction, *types.Error) {
	res, err := c.walletClient.GetTransactionsByAddress(context.Background(), &wallet_pb.AddressRequest{
		Address: address,
	})
	if err != nil {
		return nil, types.ExtractGrpcError(err)
	}
	transactions := make([]models.Transaction, 0)
	for _, t := range res.Transactions {
		transactions = append(transactions, models.Transaction{
			TxId:             t.TxId,
			RoomId:           t.RoomId,
			CurrencyId:       t.CurrencyId,
			CurrencyFullName: t.CurrencyFullName,
			FromAddress:      t.FromAddress,
			ToAddress:        t.ToAddress,
			TransactionAt:    t.TranstionAt,
		})
	}
	return transactions, nil
}
