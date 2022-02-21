package repository

import (
	"os"
	"strconv"
	"wallet/model"
)

type IWalletRepo interface {
	FindAllWallet() ([]model.WalletData, error)
	CreateWallet(username string) (model.WalletData, error)
	UpdateWalletByUsername(username string, newBalance int) (model.WalletData, error)
	WalletByUsername(username string) (model.WalletData, error)
	CheckUser(username string) bool
}

type WalletRepo struct {
	db map[string]*model.WalletData
}

func NewWalletRepo(database map[string]*model.WalletData) IWalletRepo {
	return &WalletRepo{
		db: database,
	}
}

func (w *WalletRepo) FindAllWallet() ([]model.WalletData, error) {
	wallets := []model.WalletData{}
	for _, wallet := range w.db {
		wallets = append(wallets, *wallet)
	}
	return wallets, nil
}

func (w *WalletRepo) CreateWallet(username string) (model.WalletData, error) {

	initialBalance := os.Getenv("initialBalanceAmount")
	balance, _ := strconv.Atoi(initialBalance)

	wallet := &model.WalletData{Username: username, Balance: balance}

	w.db[username] = wallet

	return *wallet, nil
}

func (w *WalletRepo) UpdateWalletByUsername(username string, newBalance int) (model.WalletData, error) {
	w.db[username].Balance = newBalance
	return *w.db[username], nil
}

func (w *WalletRepo) WalletByUsername(username string) (model.WalletData, error) {
	if wallet, ok := w.db[username]; ok {
		return *wallet, nil
	}
	return model.WalletData{}, nil
}

func (w *WalletRepo) CheckUser(username string) bool {
	_, ok := w.db[username]
	return ok
}
