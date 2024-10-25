package services

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"
	"voolibow_wallet/internal/models"
	"voolibow_wallet/internal/repository"
	"voolibow_wallet/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type (
	WalletService interface {
		CreateWallet(*models.Wallet) (string, *types.Error)
		GetWalletBalance(int32, string) (float64, *types.Error)
		TransferETH(*models.TransferETHDTO) *types.Error
		GetWalletsByUserId(int32) ([]models.Wallet, *types.Error)
		GetTransactionsByAddress(string) ([]models.Transaction, *types.Error)
	}
	walletService struct {
		repository   repository.WalletRepository
		infuraClient *ethclient.Client
	}
)

func NewWalletService(repository repository.WalletRepository, infuraClient *ethclient.Client) WalletService {
	return &walletService{
		repository:   repository,
		infuraClient: infuraClient,
	}
}

func (c *walletService) CreateWallet(data *models.Wallet) (string, *types.Error) {
	exists, err := c.repository.CheckWalletExists(data.UserId, data.CurrencyId)
	if err != nil {
		return "", err
	}
	if exists {
		return "", types.NewBadRequestError("you already have the wallet")
	}

	privateKey, rerr := crypto.GenerateKey()
	if rerr != nil {
		return "", types.NewInternalError("failed to create key")
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("SAVE BUT DO NOT SHARE THIS (Private Key):", hexutil.Encode(privateKeyBytes))

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		// log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return "", types.NewInternalError("cannot assert type")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("Public Key:", hexutil.Encode(publicKeyBytes))

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("Address:", address)

	data.Address = address
	data.PrivateKey = hexutil.Encode(privateKeyBytes)
	data.PublicKey = hexutil.Encode(publicKeyBytes)

	account := common.HexToAddress(address)
	balance, rerr := c.infuraClient.BalanceAt(context.Background(), account, nil)
	if rerr != nil {
		return "", types.NewInternalError("failed to get the balance")
	}
	balanceInEth := weiToEth(balance)
	data.LastBalance = balanceInEth
	data.LastBalanceAt = time.Now().UTC()

	err = c.repository.CreateWallet(data)
	if err != nil {
		return "", err
	}

	return address, nil
}

func (c *walletService) GetWalletBalance(userId int32, currencyId string) (float64, *types.Error) {
	address, err := c.repository.GetWalletAddressByUserIdAndCurrencyId(userId, currencyId)
	if err != nil {
		return 0, err
	}

	account := common.HexToAddress(address)
	balance, rerr := c.infuraClient.BalanceAt(context.Background(), account, nil)
	if rerr != nil {
		return 0, types.NewInternalError("failed to get the balance")
	}

	// Convert balance from wei to ETH (float64) for display
	balanceInEth := weiToEth(balance)

	// Update last balance in the database
	err = c.repository.UpdateLastBalance(&models.UpdateLastBalanceDTO{
		UserId:        userId,
		CurrencyId:    currencyId,
		LastBalance:   balanceInEth,
		LastBalanceAt: time.Now().UTC(),
	})
	if err != nil {
		return 0, err
	}

	return balanceInEth, nil
}

func (c *walletService) GetWalletBalanceByAddress(address string) (float64, *types.Error) {
	wallet, err := c.repository.GetWalletByAddress(address)
	if err != nil {
		return 0, err
	}

	account := common.HexToAddress(address)
	balance, rerr := c.infuraClient.BalanceAt(context.Background(), account, nil)
	if rerr != nil {
		return 0, types.NewInternalError("failed to get the balance")
	}

	// Convert balance from wei to ETH (float64) for display
	balanceInEth := weiToEth(balance)

	// Update last balance in the database
	err = c.repository.UpdateLastBalance(&models.UpdateLastBalanceDTO{
		UserId:        wallet.UserId,
		CurrencyId:    wallet.CurrencyId,
		LastBalance:   balanceInEth,
		LastBalanceAt: time.Now().UTC(),
	})
	if err != nil {
		return 0, err
	}

	return balanceInEth, nil
}

func (c *walletService) TransferETH(data *models.TransferETHDTO) *types.Error {
	wallet, err := c.repository.GetWalletByUserIdAndCurrencyId(data.UserId, "ETH")
	if err != nil {
		return err
	}
	if strings.HasPrefix(wallet.PrivateKey, "0x") {
		wallet.PrivateKey = wallet.PrivateKey[2:]
	}
	fmt.Println(wallet)
	privateKey, rerr := crypto.HexToECDSA(wallet.PrivateKey)
	if rerr != nil {
		fmt.Println("#1")
		return types.NewInternalError(rerr.Error())
	}
	fromAddress := common.HexToAddress(wallet.Address)

	nonce, rerr := c.infuraClient.PendingNonceAt(context.Background(), fromAddress)
	if rerr != nil {
		fmt.Println("#2")
		return types.NewInternalError(rerr.Error())
	}

	// value := big.NewInt(1000000000000000000) // in wei (1 eth)
	value := ethToWei(data.Amount)
	gasLimit := uint64(21000) // in units
	gasPrice, rerr := c.infuraClient.SuggestGasPrice(context.Background())
	if rerr != nil {
		fmt.Println("#3")
		return types.NewInternalError(rerr.Error())
	}

	toAddress := common.HexToAddress(data.ToAddress)
	var transactionData []byte
	tx := ethTypes.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, transactionData)

	chainID, rerr := c.infuraClient.NetworkID(context.Background())
	if rerr != nil {
		fmt.Println("#4")
		return types.NewInternalError(rerr.Error())

	}

	signedTx, rerr := ethTypes.SignTx(tx, ethTypes.NewEIP155Signer(chainID), privateKey)
	if rerr != nil {
		fmt.Println("#5")
		return types.NewInternalError(rerr.Error())
	}

	rerr = c.infuraClient.SendTransaction(context.Background(), signedTx)
	if rerr != nil {
		return types.NewInternalError(rerr.Error())
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
	time.Sleep(time.Second * 2)
	res, err := c.GetWalletBalance(data.UserId, "ETH")
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(wallet.Address, res)
	res, err = c.GetWalletBalanceByAddress(data.ToAddress)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(data.ToAddress, res)
	err = c.repository.CreateTransaction(&models.Transaction{
		TxId:             signedTx.Hash().Hex(),
		CurrencyId:       "ETH",
		CurrencyFullName: "Ethereum",
		FromAddress:      fromAddress.String(),
		ToAddress:        data.ToAddress,
		RoomId:           data.RoomId,
		TransactionAt:    time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *walletService) GetTransactionsByAddress(address string) ([]models.Transaction, *types.Error) {
	return c.repository.GetTransactionsByAddress(address)
}
func (c *walletService) GetWalletsByUserId(userId int32) ([]models.Wallet, *types.Error) {
	return c.repository.GetWalletsByUserId(userId)
}
func ethToWei(eth float64) *big.Int {
	ethInWei := new(big.Float).SetFloat64(eth)
	ethInWei.Mul(ethInWei, big.NewFloat(math.Pow10(18))) // Convert ETH to Wei
	wei := new(big.Int)
	wei, _ = ethInWei.Int(wei)
	return wei
}

// Convert balance from wei to ETH (float64)
func weiToEth(balance *big.Int) float64 {
	ethInWei := new(big.Float).SetInt(balance)
	ethInWei.Quo(ethInWei, big.NewFloat(math.Pow10(18))) // Convert wei to ETH
	eth, _ := ethInWei.Float64()
	return eth
}

// 1	"ETH"	"Ethereum"	"0x9B6CFDEaA8D95525D96F6D15B266033F1F9761DA"	"0x04eb88dd60e605ea161f989a5eb884b8643c17eb2396f9ce9d5b1b9dd791916e73030bf740f520fb5427cb3d46f040b9a8872f9e8cc266cd7b124e7d2ad4bf2c90"	"0x43e3dc92720deb91a3a29fb8d9ec584fea68c5c0043f8b5551c41522ed3c73ea"	0	"2024-04-05 11:42:01.94888+00"	"2024-04-05 11:42:01.949646+00"
