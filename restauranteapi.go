package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"restauranteapi/dishes"
	"restauranteapi/helper"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

var mongodbvar helper.DatabaseX
var redisclient *redis.Client

var db *sql.DB
var err error

// Looks after the main routing
//
func main() {

	redisclient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	fmt.Println(">>> Web Server: restauranteAPI.exe running.")
	fmt.Println("Loading reference data in cache - Redis")
	loadreferencedatainredis()
	APIServerPort, _ := redisclient.Get("API.APIServer.Port").Result()
	MongoDBLocation, _ := redisclient.Get("API.MongoDB.Location").Result()
	MongoDBDatabase, _ := redisclient.Get("API.MongoDB.Database").Result()

	mongodbvar.Location = MongoDBLocation
	mongodbvar.Database = MongoDBDatabase

	fmt.Println("Running... Listening to " + APIServerPort)
	fmt.Println("MongoDB location: " + MongoDBLocation)
	fmt.Println("MongoDB database: " + MongoDBDatabase)

	router := XNewRouter()

	// handle using the router mux
	//
	http.Handle("/", router) // setting router rule

	err := http.ListenAndServe(APIServerPort, nil) // setting listening port
	if err != nil {
		//using the mux router
		log.Fatal("ListenAndServe: ", err)
	}
}

func loadreferencedatainredis() {

	// err = client.Set("MongoDB.Location", "{\"MongoDB.Location\":\"192.168.2.180\"}", 0).Err()
	err = redisclient.Set("API.MongoDB.Location", "192.168.2.180", 0).Err()
	err = redisclient.Set("API.MongoDB.Database", "restaurante", 0).Err()
	err = redisclient.Set("API.APIServer.IPAddress", "192.168.2.170", 0).Err()
	err = redisclient.Set("API.APIServer.Port", ":1520", 0).Err()
}

func dishlist(httpwriter http.ResponseWriter, req *http.Request) {

	var dishlist = dishes.GetAll(redisclient)

	json.NewEncoder(httpwriter).Encode(&dishlist)
}

func dishfind(httpwriter http.ResponseWriter, httprequest *http.Request) {

	dishfound := dishes.Dish{}

	dishtofind := httprequest.FormValue("dishname") // This is the key, must be unique

	params := httprequest.URL.Query()
	parmdishname := params.Get("dishname")

	fmt.Println("params.Get parmdishname")
	fmt.Println(parmdishname)

	fmt.Println("httprequest.FormValue dishname")
	fmt.Println(dishtofind)

	dishfound = dishes.Find(redisclient, dishtofind)

	json.NewEncoder(httpwriter).Encode(&dishfound)
}

func dishadd(httpwriter http.ResponseWriter, req *http.Request) {

	dishtoadd := dishes.Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")
	fmt.Println("dishtoadd.Name")
	fmt.Println(dishtoadd.Name)

	// params := req.URL.Query()
	// dishtoadd.Name = params.Get("dishname")
	// dishtoadd.Type = params.Get("dishtype")
	// dishtoadd.Price = params.Get("dishprice")
	// dishtoadd.GlutenFree = params.Get("dishglutenfree")
	// dishtoadd.DairyFree = params.Get("dishdairyfree")
	// dishtoadd.Vegetarian = params.Get("dishvegetarian")

	ret := dishes.Dishadd(redisclient, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

func dishupdate(httpwriter http.ResponseWriter, req *http.Request) {

	dishtoupdate := dishes.Dish{}

	dishtoupdate.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoupdate.Type = req.FormValue("dishtype")
	dishtoupdate.Price = req.FormValue("dishprice")
	dishtoupdate.GlutenFree = req.FormValue("dishglutenfree")
	dishtoupdate.DairyFree = req.FormValue("dishdairyfree")
	dishtoupdate.Vegetarian = req.FormValue("dishvegetarian")
	fmt.Println("dishtoupdate.Name")
	fmt.Println(dishtoupdate.Name)

	// params := req.URL.Query()
	// dishtoadd.Name = params.Get("dishname")
	// dishtoadd.Type = params.Get("dishtype")
	// dishtoadd.Price = params.Get("dishprice")
	// dishtoadd.GlutenFree = params.Get("dishglutenfree")
	// dishtoadd.DairyFree = params.Get("dishdairyfree")
	// dishtoadd.Vegetarian = params.Get("dishvegetarian")

	ret := dishes.Dishupdate(redisclient, dishtoupdate)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

func dishalsolist(httpwriter http.ResponseWriter, req *http.Request) {

	var dishlist = dishes.GetAll(redisclient)

	json.NewEncoder(httpwriter).Encode(&dishlist)
}

type rediscachevalues struct {
	MongoDBLocation string
	MongoDBDatabase string
	APIServerPort   string
	APIServerIP     string
}

func getcachedvalues(httpwriter http.ResponseWriter, req *http.Request) {

	var rv = new(rediscachevalues)

	rv.MongoDBLocation, _ = redisclient.Get("MongoDB.Location").Result()
	rv.MongoDBDatabase, _ = redisclient.Get("MongoDB.Database").Result()
	rv.APIServerPort, _ = redisclient.Get("APIServer.Port").Result()
	rv.APIServerIP, _ = redisclient.Get("APIServer.IP").Result()

	json.NewEncoder(httpwriter).Encode(&rv)
}
