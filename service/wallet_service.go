package service

import (
	"fmt"
	"os"
	"strconv"
	"wallet/model"
	"wallet/repository"
)

var (
	ErrorNotFound     = fmt.Errorf("The wallet not found.")
	ErrorBalanceLimit = fmt.Errorf("The wallet cannot under below the minimum value.")
)

type IWalletService interface {
	Wallets() ([]model.WalletData, error)
	CreateWallet(username string) (model.WalletData, error)
	UpdateWalletByUsername(username string, balanceAdd int) (model.WalletData, error)
	WalletByUsername(username string) (model.WalletData, error)
}

type WalletService struct {
	repository repository.IWalletRepo
}

func NewWalletService(repository repository.IWalletRepo) IWalletService {
	return &WalletService{
		repository: repository,
	}
}

func (w *WalletService) Wallets() ([]model.WalletData, error) {
	return w.repository.FindAllWallet()
}

func (w *WalletService) WalletByUsername(username string) (model.WalletData, error) {
	if !w.repository.CheckUser(username) {
		return model.WalletData{}, ErrorNotFound
	}
	return w.repository.WalletByUsername(username)
}

func (w *WalletService) CreateWallet(username string) (model.WalletData, error) {
	if w.repository.CheckUser(username) {
		return w.repository.WalletByUsername(username)
	}
	return w.repository.CreateWallet(username)
}

func (w *WalletService) UpdateWalletByUsername(username string, balanceAdd int) (model.WalletData, error) {
	wallet, err := w.WalletByUsername(username)
	if err != nil {
		return model.WalletData{}, err
	}
	minimumBalance := os.Getenv("minimumBalanceAmount")
	minBalance, _ := strconv.Atoi(minimumBalance)

	newBalance := wallet.Balance + balanceAdd
	if newBalance < minBalance {
		return model.WalletData{}, ErrorBalanceLimit
	}

	return w.repository.UpdateWalletByUsername(username, newBalance)
}
