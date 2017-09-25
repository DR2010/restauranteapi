package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"restauranteapi/dishes"
	"restauranteapi/helper"

	_ "github.com/go-sql-driver/mysql"
)

var mongodbvar helper.DatabaseX

var db *sql.DB
var err error

// Looks after the main routing
//
func main() {

	mongodbvar.Location = "192.168.2.180"
	mongodbvar.Database = "restaurante"

	var porta string
	porta = ":1520"

	fmt.Println("Running... Listening to " + porta)
	fmt.Println("MongoDB location: " + mongodbvar.Location)
	fmt.Println("MongoDB database: " + mongodbvar.Database)

	router := XNewRouter()

	// handle using the router mux
	//
	http.Handle("/", router) // setting router rule

	err := http.ListenAndServe(porta, nil) // setting listening port
	if err != nil {
		//using the mux router
		log.Fatal("ListenAndServe: ", err)
	}
}

func dishlist(httpwriter http.ResponseWriter, req *http.Request) {

	var dishlist = dishes.GetAll(mongodbvar)

	json.NewEncoder(httpwriter).Encode(&dishlist)
}

func dishadd(httpwriter http.ResponseWriter, req *http.Request) {

	dishtoadd := dishes.Dish{}
	fmt.Println("Aqui 001")

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")
	fmt.Println("dishtoadd.Name")
	fmt.Println(dishtoadd.Name)

	// Keeping one as an example of retrieving query string data
	// Nao sei passar por GoLang os valores no Body do URL request
	// So' sei passar pelo query string, entao este e' o codigo
	// Quando conseguir passar o Body, mudo o codigo

	// params := req.URL.Query()
	// dishtoadd.Name = params.Get("dishname")
	// dishtoadd.Type = params.Get("dishtype")
	// dishtoadd.Price = params.Get("dishprice")
	// dishtoadd.GlutenFree = params.Get("dishglutenfree")
	// dishtoadd.DairyFree = params.Get("dishdairyfree")
	// dishtoadd.Vegetarian = params.Get("dishvegetarian")

	ret := dishes.Dishadd(mongodbvar, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

func dishalsolist(httpwriter http.ResponseWriter, req *http.Request) {

	var dishlist = dishes.GetAll(mongodbvar)

	json.NewEncoder(httpwriter).Encode(&dishlist)
}
