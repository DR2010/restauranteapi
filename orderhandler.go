// Package main is the main package
// -------------------------------------
// .../restauranteapi/orderhandler.go
// -------------------------------------
package main

import (
	"encoding/json"
	"fmt"
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

	objtoaction := orders.Order{}

	objtoaction.ID = req.FormValue("orderID")             // This is the key, must be unique
	objtoaction.ClientID = req.FormValue("orderClientID") // This is the key, must be unique
	objtoaction.ClientName = req.FormValue("orderClientName")
	objtoaction.Date = req.FormValue("orderDate")
	objtoaction.Time = req.FormValue("orderTime")
	objtoaction.Foodeatplace = req.FormValue("foodeatplace")

	_, recordstatus := orders.Find(redisclient, objtoaction.ID)

	if recordstatus == "200 OK" {
		fmt.Println("recordstatus")
		fmt.Println(recordstatus)
		http.Error(httpwriter, "Record already exists.", 422)
		return
	}

	ret := orders.Add(redisclient, objtoaction)

	if ret.IsSuccessful == "Y" {
		// do something

		fmt.Println("Order added successfully:" + objtoaction.ClientName)
	}
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
