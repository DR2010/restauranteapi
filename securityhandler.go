// Package main is the main package
// -------------------------------------
// .../restauranteapi/btccotacaohandler.go
// -------------------------------------
package main

import (
	"encoding/json"
	"net/http"
	"restauranteapi/security"
)

// Hsecuritylogin is
func Hsecuritylogin(httpwriter http.ResponseWriter, req *http.Request) {

	var userid = req.FormValue("userid")
	var password = req.FormValue("password")

	// params := req.URL.Query()
	// cotacaotoadd.Currency = params.Get("Currency")
	// cotacaotoadd.Balance = params.Get("Balance")

	token, _ := security.ValidateUserCredentials(redisclient, userid, password)

	if token == "Error" {
		httpwriter.WriteHeader(http.StatusInternalServerError)
		httpwriter.Write([]byte("500 - Something bad happened!"))
	}
	json.NewEncoder(httpwriter).Encode(&token)

}

// Hsecuritysignup is
func Hsecuritysignup(httpwriter http.ResponseWriter, req *http.Request) {

	var userInsert security.Credentials

	userInsert.UserID = req.FormValue("userid")
	userInsert.Password = req.FormValue("password")
	userInsert.PasswordValidate = req.FormValue("passwordvalidate")

	token := ""
	_, resfind := security.Find(redisclient, userInsert.UserID)
	if resfind == "200 OK" {
		token = "User already exists"
	}

	// Add user
	results := security.Useradd(redisclient, userInsert)
	if results.ErrorCode == "200 OK" {
		token = results.ReturnedValue
	}

	json.NewEncoder(httpwriter).Encode(&token)

}
