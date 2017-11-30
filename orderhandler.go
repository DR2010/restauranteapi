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
	orders "restauranteapi/orders"

	_ "github.com/go-sql-driver/mysql"
)

// Hfind finds orders
func Hfind(httpwriter http.ResponseWriter, httprequest *http.Request) {

	objfound := orders.Order{}

	objtofind := httprequest.FormValue("orderid") // This is the key, must be unique

	objfound, _ = orders.Find(redisclient, objtofind)

	json.NewEncoder(httpwriter).Encode(&objfound)
}

// Horderadd add orders
func Horderadd(httpwriter http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)
	// bodystr := string(bodybyte[:])

	type dcOrder struct {
		OrderID         string // random ID for order, yet to define algorithm
		OrderClientID   string // Client Name
		OrderClientName string // Client ID in case they logon
		OrderDate       string // Order Date
		OrderTime       string // Order Time
		Foodeatplace    string // Order Time
		Status          string // Order Time
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
	objtoaction.DeliveryContactPhone = req.FormValue("orderDeliveryContactPhone")
	objtoaction.DeliveryFee = req.FormValue("orderDeliveryFee")

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
