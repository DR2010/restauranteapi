// Package main is the main package
// -------------------------------------
// .../restauranteapi/orderhandler.go
// -------------------------------------
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"restauranteapi/helper"
	orders "restauranteapi/orders"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// Hfind finds orders
func Hfind(httpwriter http.ResponseWriter, httprequest *http.Request) {

	objfound := orders.Order{}

	objtofind := httprequest.FormValue("orderid") // This is the key, must be unique

	objfound, _ = orders.Find(redisclient, objtofind)

	json.NewEncoder(httpwriter).Encode(&objfound)
}

// Hdishfind is
func Horderfind(httpwriter http.ResponseWriter, httprequest *http.Request) {

	redisclient := helper.GetRedisPointer()

	orderfound := orders.Order{}

	ordertofind := httprequest.FormValue("orderid") // This is the key, must be unique

	params := httprequest.URL.Query()
	parmorderid := params.Get("orderid")

	fmt.Println("params.Get parmorderid")
	fmt.Println(parmorderid)

	orderfound, _ = orders.Find(redisclient, ordertofind)

	json.NewEncoder(httpwriter).Encode(&orderfound)
}

// Horderadd add orders
func Horderadd(httpwriter http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)
	// bodystr := string(bodybyte[:])

	type dcOrderItem struct {
		Pratoname  string // random ID for order, yet to define algorithm
		Quantidade string // Client Name
		Preco      string // Client ID in case they logon
	}

	type dcOrder struct {
		OrderID         string // random ID for order, yet to define algorithm
		OrderClientID   string // Client Name
		OrderClientName string // Client ID in case they logon
		OrderDate       string // Order Date
		OrderTime       string // Order Time
		Foodeatplace    string // Order Time
		Status          string // Order Time
		Pratos          []dcOrderItem
	}

	var objtoaction dcOrder
	err = json.Unmarshal(bodybyte, &objtoaction)

	objtoaction.OrderID = objtoaction.OrderClientName + objtoaction.OrderDate + "01"

	_, recordstatus := orders.Find(redisclient, objtoaction.OrderID)

	if recordstatus == "200 OK" {
		fmt.Println("recordstatus")
		fmt.Println(recordstatus)
		http.Error(httpwriter, "Record already exists.", 422)

		return
	}

	objtoactionMAP := orders.Order{}
	objtoactionMAP.ID = objtoaction.OrderID
	objtoactionMAP.ClientID = objtoaction.OrderClientID
	objtoactionMAP.ClientName = objtoaction.OrderClientName
	objtoactionMAP.Date = objtoaction.OrderDate
	objtoactionMAP.Time = objtoaction.OrderTime

	var slen = len(objtoaction.Pratos)
	objtoactionMAP.Items = make([]orders.Item, slen)

	var totalgeral = 0

	// I have to remove the header coming from the caller.
	// Perhaps the caller should suppress the header somehow

	var destindex = 0

	for index, element := range objtoaction.Pratos {
		// index is the index where we are
		// element is the element from someSlice for where we are

		if index == 0 {
			continue
		}

		destindex = index - 1

		objtoactionMAP.Items[destindex].PratoName = element.Pratoname
		objtoactionMAP.Items[destindex].Price = element.Preco
		objtoactionMAP.Items[destindex].Quantidade = element.Quantidade

		prc, _ := strconv.Atoi(element.Preco)
		qty, _ := strconv.Atoi(element.Quantidade)
		tot := prc * qty
		totalgeral = totalgeral + tot

		objtoactionMAP.Items[destindex].Total = strconv.Itoa(tot)
	}
	objtoactionMAP.TotalGeral = strconv.Itoa(totalgeral)

	ret := orders.Add(redisclient, objtoactionMAP)

	if ret.IsSuccessful == "Y" {
		// do something

		fmt.Println("Order added successfully:" + objtoaction.OrderClientName)

		type RespAddOrder struct {
			ID string
		}

		// return value
		obj := &RespAddOrder{ID: objtoaction.OrderID}
		bresp, _ := json.Marshal(obj)

		fmt.Fprintf(httpwriter, string(bresp)) // write data to response
	}

	return
}

// Hupdate updates orders
func Hupdate(httpwriter http.ResponseWriter, req *http.Request) {

	objtoaction := orders.Order{}

	objtoaction.ClientID = req.FormValue("orderClientID") // This is the key, must be unique
	objtoaction.ClientName = req.FormValue("orderClientName")
	objtoaction.Date = req.FormValue("orderDate")

	ret := orders.Update(redisclient, objtoaction)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

// Hdelete delete orders
func Hdelete(httpwriter http.ResponseWriter, req *http.Request) {

	objtoupdate := orders.Order{}

	objtoupdate.ClientID = req.FormValue("orderID") // This is the key, must be unique

	ret := orders.Delete(redisclient, objtoupdate.ClientID)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

// Halsolist list orders
func Halsolist(httpwriter http.ResponseWriter, req *http.Request) {

	var orderlist = orders.Getall(redisclient)

	json.NewEncoder(httpwriter).Encode(&orderlist)
}

// OrderList also list orders
func OrderList(httpwriter http.ResponseWriter, req *http.Request) {

	var orderlist = orders.Getall(redisclient)

	json.NewEncoder(httpwriter).Encode(&orderlist)
}
