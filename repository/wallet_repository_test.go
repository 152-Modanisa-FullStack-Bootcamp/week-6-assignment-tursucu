package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"wallet/model"
)

func TestWalletRepo_FindAllWallet(t *testing.T) {
	wallets := map[string]*model.WalletData{
		"niyazi": {
			Username: "niyazi",
			Balance:  500,
		},
		"nuri": {
			Username: "nuri",
			Balance:  200,
		},
		"kazım": {
			Username: "kazım",
			Balance:  -99,
		},
	}
	repo := NewWalletRepo(wallets)

	var list []model.WalletData
	for _, wallet := range wallets {
		list = append(list, *wallet)
	}
	result, err := repo.FindAllWallet()

	assert.ElementsMatch(t, list, result)
	assert.Nil(t, err)
}

func TestWalletRepo_WalletByUsername(t *testing.T) {
	wallets := map[string]*model.WalletData{
		"niyazi": {
			Username: "niyazi",
			Balance:  500,
		},
		"nuri": {
			Username: "nuri",
			Balance:  200,
		},
		"kazım": {
			Username: "kazım",
			Balance:  -99,
		},
	}
	repo := NewWalletRepo(wallets)

	cases := []struct {
		name           string
		username       string
		expectedWallet model.WalletData
	}{
		{
			name:           "niyazi wallet wait",
			username:       "niyazi",
			expectedWallet: model.WalletData{Username: "niyazi", Balance: 500},
		},
		{
			name:           "nuri wallet wait",
			username:       "nuri",
			expectedWallet: model.WalletData{Username: "nuri", Balance: 200},
		},
		{
			name:           "kazım wallet wait",
			username:       "kazım",
			expectedWallet: model.WalletData{Username: "kazım", Balance: -99},
		},
		{
			name:           "mehmet none wallet",
			username:       "mehmet",
			expectedWallet: model.WalletData{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := repo.WalletByUsername(c.username)

			assert.Equal(t, c.expectedWallet, result)
			assert.Nil(t, err)
		})
	}
}

func TestWalletRepo_CreateWallet(t *testing.T) {
	wallets := map[string]*model.WalletData{
		"mehmet": {
			Username: "mehmet",
			Balance:  350,
		},
		"nur": {
			Username: "nur",
			Balance:  90,
		},
	}
	repo := NewWalletRepo(wallets)

	newWallet := &model.WalletData{Username: "hasan", Balance: 0}
	repo.CreateWallet(newWallet.Username)

	assert.Contains(t, wallets, newWallet.Username)
	assert.Equal(t, "hasan", newWallet.Username)
	assert.Equal(t, 0, newWallet.Balance)
}

func TestUpdate(t *testing.T) {
	wallets := map[string]*model.WalletData{
		"harun": {
			Username: "harun",
			Balance:  550,
		},
		"naz": {
			Username: "naz",
			Balance:  -50,
		},
	}
	repo := NewWalletRepo(wallets)

	cases := []struct {
		name           string
		username       string
		balance        int
		expectedWallet model.WalletData
	}{
		{
			name:           "harun balance 50, expect wallet balance to 50",
			username:       "harun",
			balance:        50,
			expectedWallet: model.WalletData{Username: "harun", Balance: 50},
		},
		{
			name:           "naz balance -99, expect wallet balance to -99",
			username:       "naz",
			balance:        -99,
			expectedWallet: model.WalletData{Username: "naz", Balance: -99},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := repo.UpdateWalletByUsername(c.username, c.balance)

			assert.Equal(t, c.expectedWallet, result)
			assert.Nil(t, err)
		})
	}
}

func TestWalletRepo_CheckUser(t *testing.T) {
	wallets := map[string]*model.WalletData{
		"jack": {
			Username: "jack",
			Balance:  2500,
		},
		"jeff": {
			Username: "jeff",
			Balance:  500,
		},
	}
	repo := NewWalletRepo(wallets)

	cases := []struct {
		name          string
		username      string
		expectedCheck bool
	}{
		{
			name:          "check true",
			username:      "jack",
			expectedCheck: true,
		},
		{
			name:          "does not check false",
			username:      "bora",
			expectedCheck: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			exists := repo.CheckUser(c.username)
			assert.Equal(t, c.expectedCheck, exists)
		})
	}
}
