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

	var dishlist = dishes.GetAll(mongodbvar)

	json.NewEncoder(httpwriter).Encode(&dishlist)
}
