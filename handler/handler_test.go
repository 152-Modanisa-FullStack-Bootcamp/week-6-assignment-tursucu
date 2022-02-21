package handler_test

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet/handler"
	"wallet/mock"
	"wallet/model"
)

func TestHandlerServiceReturnWallets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock.NewMockIWalletService(ctrl)
	serviceReturn := []model.WalletData{
		{Username: "mehmet", Balance: 50},
		{Username: "onur", Balance: -50},
		{Username: "hamza", Balance: 150},
	}

	service.EXPECT().Wallets().Return(serviceReturn, nil).Times(1)

	h := handler.NewWalletHandler(service)
	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	w := httptest.NewRecorder()
	h.Wallets(w, r)

	response, _ := json.Marshal(serviceReturn)

	assert.Equal(t, string(response), w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
