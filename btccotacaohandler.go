// Package main is the main package
// -------------------------------------
// .../restauranteapi/btccotacaohandler.go
// -------------------------------------
package main

import (
	"encoding/json"
	"net/http"
	"restauranteapi/btcmarkets"
	"time"
)

// Hbtccotacaoadd is
func Hbtccotacaoadd(httpwriter http.ResponseWriter, req *http.Request) {

	cotacaotoadd := btcmarkets.BalanceCrypto{}

	cotacaotoadd.Currency = req.FormValue("cryptoCurrency")
	cotacaotoadd.Balance = req.FormValue("cryptoBalance")
	cotacaotoadd.CotacaoAtual = req.FormValue("cryptoCotacaoAtual")
	cotacaotoadd.ValueInCashAUD = req.FormValue("cryptoValueInCashAUD")
	cotacaotoadd.DateTime = time.Now().String()

	// params := req.URL.Query()
	// cotacaotoadd.Currency = params.Get("Currency")
	// cotacaotoadd.Balance = params.Get("Balance")
	// cotacaotoadd.CotacaoAtual = params.Get("CotacaoAtual")
	// cotacaotoadd.DateTime = params.Get("DateTime")
	// cotacaotoadd.ValueInCashAUD = params.Get("ValueInCashAUD")

	ret := btcmarkets.CryptoCotacaoAdd(redisclient, cotacaotoadd)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

// Hbtccotacaolist is a function to return a list of dishes
func Hbtccotacaolist(httpwriter http.ResponseWriter, req *http.Request) {

	params := req.URL.Query()
	var currency = params.Get("currency")
	var rows = params.Get("rows")

	var cotacaolist = btcmarkets.GetAll(redisclient, currency, rows)

	json.NewEncoder(httpwriter).Encode(&cotacaolist)
}

// Hbtccotacaolistdate is a function to return a list of dishes
func Hbtccotacaolistdate(httpwriter http.ResponseWriter, req *http.Request) {

	params := req.URL.Query()
	var currency = params.Get("currency")
	var yeardaymonth = params.Get("yeardaymonth")
	var yeardaymonthend = params.Get("yeardaymonthend")

	var cotacaolist = btcmarkets.GetDayStats(redisclient, currency, yeardaymonth, yeardaymonthend)

	json.NewEncoder(httpwriter).Encode(&cotacaolist)
}
