package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mongodb/dishes"
	"mongodb/helper"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var mongodbvar helper.DatabaseX

var db *sql.DB
var err error

// Looks after the main routing
//
func main() {

	mongodbvar.Location = "localhost"
	mongodbvar.Database = "restaurante"

	var porta string
	porta = ":1520"

	fmt.Println("Running... Listening to " + porta)

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

	type ControllerInfo struct {
		Name string
	}
	type Row struct {
		Description []string
	}
	type DisplayTemplate struct {
		Info       ControllerInfo
		FieldNames []string
		Rows       []Row
	}

	var dishlist = dishes.GetAll(mongodbvar)

	items := DisplayTemplate{}
	items.Info.Name = "Dish List"

	var numberoffields = 6

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Name"
	items.FieldNames[1] = "Type"
	items.FieldNames[2] = "Price"
	items.FieldNames[3] = "GlutenFree"
	items.FieldNames[4] = "DairyFree"
	items.FieldNames[5] = "Vegetarian"

	// Set rows to be displayed
	items.Rows = make([]Row, len(dishlist))
	// items.RowID = make([]int, len(dishlist))

	for i := 0; i < len(dishlist); i++ {
		items.Rows[i] = Row{}
		items.Rows[i].Description = make([]string, numberoffields)
		items.Rows[i].Description[0] = dishlist[i].Name
		items.Rows[i].Description[1] = dishlist[i].Type
		items.Rows[i].Description[2] = dishlist[i].Price
		items.Rows[i].Description[3] = dishlist[i].GlutenFree
		items.Rows[i].Description[4] = dishlist[i].DairyFree
		items.Rows[i].Description[5] = dishlist[i].Vegetarian
	}

	json.NewEncoder(httpwriter).Encode(&items)
}
