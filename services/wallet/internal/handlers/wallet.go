package handlers

import (
	"context"
	"voolibow_wallet/internal/models"
	"voolibow_wallet/internal/services"
	pb "voolibow_wallet/proto/api"
)

type WalletHandler struct {
	pb.UnimplementedWalletServiceServer
	walletService services.WalletService
}

func NewWalletHandler(walletService services.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

func (c *WalletHandler) RegisterWallet(ctx context.Context, request *pb.RegisterWalletRequest) (*pb.RegisterWalletResponse, error) {
	data := models.Wallet{
		UserId:           request.UserId,
		CurrencyId:       request.Currency.String(),
		CurrencyFullName: request.CurrencyFullName,
	}
	res, err := c.walletService.CreateWallet(&data)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}

	return &pb.RegisterWalletResponse{Address: res}, nil
}

func (c *WalletHandler) GetBalance(ctx context.Context, request *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	res, err := c.walletService.GetWalletBalance(request.UserId, request.Currency.String())
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}

	return &pb.GetBalanceResponse{Balance: res}, nil
}

func (c *WalletHandler) TransferETH(ctx context.Context, request *pb.TransferETHRequest) (*pb.Empty, error) {
	err := c.walletService.TransferETH(&models.TransferETHDTO{
		UserId:    request.UserId,
		ToAddress: request.ToAddress,
		Amount:    request.Amount,
		RoomId:    request.RoomId,
	})
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}

	return &pb.Empty{}, nil
}

func (c *WalletHandler) GetWalletsByUserId(ctx context.Context, request *pb.UserIdRequest) (*pb.Wallets, error) {
	res, err := c.walletService.GetWalletsByUserId(request.UserId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	wallets := make([]*pb.Wallet, 0)
	for _, w := range res {
		wallets = append(wallets, &pb.Wallet{
			UserId:           w.UserId,
			CurrencyId:       w.CurrencyId,
			CurrencyFullName: w.CurrencyFullName,
			Address:          w.Address,
			PublicKey:        w.PublicKey,
			PrivateKey:       w.PrivateKey,
			LastBalance:      w.LastBalance,
			LastBalanceAt:    w.LastBalanceAt.String(),
			CreatedAt:        w.CreatedAt.String(),
		})
	}
	return &pb.Wallets{Wallets: wallets}, nil
}

func (c *WalletHandler) GetTransactionsByAddress(ctx context.Context, request *pb.AddressRequest) (*pb.Transactions, error) {
	res, err := c.walletService.GetTransactionsByAddress(request.Address)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	transactions := make([]*pb.Transaction, 0)
	for _, t := range res {
		transactions = append(transactions, &pb.Transaction{
			TxId:             t.TxId,
			RoomId:           t.RoomId,
			CurrencyId:       t.CurrencyId,
			CurrencyFullName: t.CurrencyFullName,
			FromAddress:      t.FromAddress,
			ToAddress:        t.ToAddress,
			TranstionAt:      t.TransactionAt.String(),
		})
	}

	return &pb.Transactions{Transactions: transactions}, nil

}
