package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"wallet/model"
	"wallet/routes"
	"wallet/service"
)

type IWalletHandler interface {
	Wallets(w http.ResponseWriter, r *http.Request)
	CreateWallet(w http.ResponseWriter, r *http.Request)
	UpdateWalletByUsername(w http.ResponseWriter, r *http.Request)
	WalletByUsername(w http.ResponseWriter, r *http.Request)
}

type WalletHandler struct {
	walsvc service.IWalletService
}

type HTTPError struct {
	StatusCode int    `json:"status-code"`
	Message    string `json:"message"`
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	httpErr := HTTPError{
		StatusCode: code,
		Message:    msg,
	}
	response, _ := json.Marshal(httpErr)
	w.WriteHeader(code)
	w.Write(response)
}
func NewWalletHandler(walsvc service.IWalletService) IWalletHandler {
	return &WalletHandler{
		walsvc: walsvc,
	}
}

func (wa *WalletHandler) Wallets(w http.ResponseWriter, r *http.Request) {

	wallets, err := wa.walsvc.Wallets()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	response, err := json.Marshal(wallets)

	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("content-type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (wa *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	slug := routes.GetField(r, 0)

	wallet, err := wa.walsvc.CreateWallet(slug)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	response, err := json.Marshal(wallet)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("content-type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (wa *WalletHandler) WalletByUsername(w http.ResponseWriter, r *http.Request) {
	slug := routes.GetField(r, 0)

	wallet, err := wa.walsvc.WalletByUsername(slug)

	if err != nil && errors.Is(err, service.ErrorNotFound) {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	response, err := json.Marshal(wallet)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (wa *WalletHandler) UpdateWalletByUsername(w http.ResponseWriter, r *http.Request) {
	var wallet model.WalletData
	slug := routes.GetField(r, 0)

	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := wa.walsvc.UpdateWalletByUsername(slug, wallet.Balance)

	if err != nil && errors.Is(err, service.ErrorNotFound) {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil && errors.Is(err, service.ErrorBalanceLimit) {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	} else if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	response, err := json.Marshal(updated)

	w.Header().Add("content-type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
