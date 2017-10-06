// Package main is the main package
// -------------------------------------
// .../restauranteapi/restauranteapi.go
// -------------------------------------
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

type rediscachevalues struct {
	MongoDBLocation string
	MongoDBDatabase string
	APIServerPort   string
	APIServerIP     string
}

// Cache represents the cache data
type Cache struct {
	Key   string // cache key
	Value string // value in cache
}

func getcachedvalues(httpwriter http.ResponseWriter, req *http.Request) {

	var rv = new(rediscachevalues)

	rv.MongoDBLocation, _ = redisclient.Get("API.MongoDB.Location").Result()
	rv.MongoDBDatabase, _ = redisclient.Get("API.MongoDB.Database").Result()
	rv.APIServerPort, _ = redisclient.Get("API.APIServer.Port").Result()
	rv.APIServerIP, _ = redisclient.Get("API.APIServer.IPAddress").Result()

	keys := make([]Cache, 4)
	keys[0].Key = "API.MongoDB.Location"
	keys[0].Value = rv.MongoDBLocation

	keys[1].Key = "API.MongoDB.Database"
	keys[1].Value = rv.MongoDBDatabase

	keys[2].Key = "API.APIServer.Port"
	keys[2].Value = rv.APIServerPort

	keys[3].Key = "API.APIServer.IPAddress"
	keys[3].Value = rv.APIServerIP

	json.NewEncoder(httpwriter).Encode(&keys)
}
