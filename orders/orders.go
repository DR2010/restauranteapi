package orders

import (
	"fmt"
	"log"
	helper "restauranteapi/helper"

	"github.com/go-redis/redis"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Order is what the client wants
type Order struct {
	SystemID     bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	ID           string        // random ID for order, yet to define algorithm
	ClientName   string        // Client Name
	ClientID     string        // Client ID in case they logon
	Date         string        // Order Date
	Time         string        // Order Time
	Status       string        // Open, Completed, Cancelled
	EatMode      string        // EatIn, TakeAway, Delivery
	Foodeatplace string        // EatIn, TakeAway, Delivery
	TotalGeral   string        // Delivery phone number
	Items        []Item
}

// Item represents a single item of an order
type Item struct {
	ID         string // Sequential number of the item
	PratoName  string // Dish ID or unique name from "Dishes"
	GlutenFree string // Just Yes or No in case the dish has gluten free options
	DiaryFree  string // Just Yes or No in case the dish has this option
	Price      string // Individual price
	Quantidade string // Individual price
	Total      string // Total Price
	Tax        string // GST
}

// SearchCriteria is what the client wants
type SearchCriteria struct {
	ID                   string // random ID for order, yet to define algorithm
	ClientName           string // Client Name
	ClientID             string // Client ID in case they logon
	Date                 string // Order Date
	Time                 string // Order Time
	Status               string // Open, Completed, Cancelled
	EatMode              string // EatIn, TakeAway, Delivery
	DeliveryMode         string // Internal, UberEats,
	DeliveryFee          string // Delivery Fee
	DeliveryLocation     string // Address
	DeliveryContactPhone string // Delivery phone number
}

// Add is for export
func Add(redisclient *redis.Client, objtoinsert Order) helper.Resultado {

	database := new(helper.DatabaseX)
	database.Collection = "orders"
	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Insert(objtoinsert)

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Order added"
	res.IsSuccessful = "Y"

	return res
}

// Find is to find stuff
func Find(redisclient *redis.Client, objtofind string) (Order, string) {

	database := new(helper.DatabaseX)
	database.Collection = "orders"
	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	objkey := objtofind
	objnull := Order{}

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(database.Database).C(database.Collection)

	result := []Order{}
	err1 := c.Find(bson.M{"id": objkey}).All(&result)
	if err1 != nil {
		log.Fatal(err1)
	}

	var numrecsel = len(result)

	if numrecsel <= 0 {
		return objnull, "404 Not found"
	}

	return result[0], "200 OK"
}

// Getall works
func Getall(redisclient *redis.Client) []Order {

	database := new(helper.DatabaseX)

	database.Collection = "orders"

	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	fmt.Println("database.Location")
	fmt.Println(database.Location)

	session, err := mgo.Dial(database.Location)

	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(database.Database).C(database.Collection)

	var results []Order

	err = c.Find(nil).All(&results)
	if err != nil {
		// TODO: Do something about the error
	} else {
		return results
	}

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// Update is
func Update(redisclient *redis.Client, objtoupdate Order) helper.Resultado {

	database := new(helper.DatabaseX)
	database.Collection = "orders"
	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Update(bson.M{"name": objtoupdate.ID}, objtoupdate)

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Something Happened"
	res.IsSuccessful = "Y"

	return res
}

// Delete is
func Delete(redisclient *redis.Client, objtodeletekey string) helper.Resultado {

	database := new(helper.DatabaseX)
	database.Collection = "dishes"
	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()
	database.Collection = "dishes"

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Remove(bson.M{"ID": objtodeletekey})

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Dish deleted successfully"
	res.IsSuccessful = "Y"

	return res
}
