package service_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"wallet/mock"
	"wallet/model"
	"wallet/service"
)

func mockRepository(t *testing.T) *mock.MockIWalletRepo {
	return mock.NewMockIWalletRepo(gomock.NewController(t))
}

func TestWalletsService_ReturnWalletsAll(t *testing.T) {
	returnData := []model.WalletData{
		{Username: "david", Balance: 0},
		{Username: "melih", Balance: 1500},
	}

	mockWalletRepository := mockRepository(t)

	mockWalletRepository.EXPECT().FindAllWallet().Return(returnData, nil).Times(1)

	s := service.NewWalletService(mockWalletRepository)
	wallets, err := s.Wallets()

	assert.Equal(t, returnData, wallets)
	assert.Nil(t, err)
}

func TestWalletByUsername_ReturnError(t *testing.T) {
	mockWalletRepository := mockRepository(t)
	mockWalletRepository.EXPECT().CheckUser("hamza").Return(false)

	s := service.NewWalletService(mockWalletRepository)
	wallet, err := s.WalletByUsername("hamza")

	assert.Empty(t, wallet)
	assert.EqualError(t, err, service.ErrorNotFound.Error())
}

func TestWalletByUsername_ReturnWallet(t *testing.T) {
	mockWalletRepository := mockRepository(t)
	mockWalletRepository.EXPECT().CheckUser("hamza").Return(true).Times(1)

	mockWallet := model.WalletData{Username: "hamza", Balance: 500}
	mockWalletRepository.EXPECT().WalletByUsername("hamza").Return(mockWallet, nil).Times(1)

	s := service.NewWalletService(mockWalletRepository)
	wallet, err := s.WalletByUsername("hamza")

	assert.Equal(t, mockWallet, wallet)
	assert.Nil(t, err)
}

func TestCreateWallet_ReturnCheckUser(t *testing.T) {
	mockWalletRepository := mockRepository(t)
	mockWalletRepository.EXPECT().CheckUser("hamza").Return(true).Times(1)

	mockWallet := model.WalletData{Username: "hamza", Balance: 200}
	mockWalletRepository.EXPECT().WalletByUsername("hamza").Return(mockWallet, nil).Times(1)

	s := service.NewWalletService(mockWalletRepository)
	wallet, err := s.CreateWallet("hamza")

	assert.Equal(t, mockWallet, wallet)
	assert.Nil(t, err)
}

func TestCreateWallet_ReturnNewWallet(t *testing.T) {
	mockWalletRepository := mockRepository(t)
	mockWalletRepository.EXPECT().CheckUser("mahmut").Return(false)

	mockWallet := &model.WalletData{Username: "mahmut", Balance: 0}
	mockWalletRepository.EXPECT().CreateWallet("mahmut").Return(*mockWallet, nil).Times(1)

	s := service.NewWalletService(mockWalletRepository)
	wallet, err := s.CreateWallet("mahmut")

	assert.Equal(t, *mockWallet, wallet)
	assert.Nil(t, err)
}

func TestUpdateWallet(t *testing.T) {
	cases := []struct {
		name           string
		username       string
		newBalance     int
		checkUser      bool
		currentWallet  model.WalletData
		expectedWallet model.WalletData
		expectedError  error
	}{
		{
			name:           "murtaza add 750, wallet updated",
			username:       "murtaza",
			newBalance:     750,
			checkUser:      true,
			currentWallet:  model.WalletData{Username: "murtaza", Balance: 0},
			expectedWallet: model.WalletData{Username: "murtaza", Balance: 750},
			expectedError:  nil,
		},
		{
			name:           "murtaza add 5000 but nil error and check user exit",
			username:       "murtaza",
			newBalance:     5000,
			checkUser:      false,
			currentWallet:  model.WalletData{},
			expectedWallet: model.WalletData{},
			expectedError:  service.ErrorNotFound,
		},
		{
			name:           "murtaza want to attract 500, but error below limit",
			username:       "murtaza",
			newBalance:     -500,
			checkUser:      true,
			currentWallet:  model.WalletData{Username: "murtaza", Balance: 0},
			expectedWallet: model.WalletData{},
			expectedError:  service.ErrorBalanceLimit,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockWalletRepository := mockRepository(t)
			mockWalletRepository.EXPECT().CheckUser(c.username).Return(c.checkUser)

			if c.checkUser {
				mockWalletRepository.EXPECT().WalletByUsername(c.username).Return(c.currentWallet, nil)
			}

			if c.expectedWallet != (model.WalletData{}) {
				mockWalletRepository.EXPECT().UpdateWalletByUsername(c.username, c.newBalance).Return(c.expectedWallet, nil)
			}

			s := service.NewWalletService(mockWalletRepository)
			updatedWallet, err := s.UpdateWalletByUsername(c.username, c.newBalance)

			assert.Equal(t, c.expectedWallet, updatedWallet)
			assert.Equal(t, c.expectedError, err)
		})
	}
}
