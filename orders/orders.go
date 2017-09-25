package orders

import (
	"log"
	helper "restauranteapi/helper"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Order is what the client wants
type Order struct {
	SystemID             bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	ID                   string        // random ID for order, yet to define algorithm
	ClientName           string        // Client Name
	ClientID             string        // Client ID in case they logon
	Date                 string        // Order Date
	Time                 string        // Order Time
	Status               string        // Open, Completed, Cancelled
	EatMode              string        // EatIn, TakeAway, Delivery
	DeliveryMode         string        // Internal, UberEats,
	DeliveryFee          string        // Delivery Fee
	DeliveryLocation     string        // Address
	DeliveryContactPhone string        // Delivery phone number
	Items                Item
}

// Item represents a single item of an order
type Item struct {
	ID         string // Sequential number of the item
	DishID     string // Dish ID or unique name from "Dishes"
	GlutenFree string // Just Yes or No in case the dish has gluten free options
	DiaryFree  string // Just Yes or No in case the dish has this option
	Price      string // Individual price
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

// Add an order
func Add(database helper.DatabaseX, orderInsert Order) helper.Resultado {

	database.Collection = "orders"

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Insert(orderInsert)

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Order placed successfully"
	res.IsSuccessful = "Y"

	return res
}

// Search for orders allow multiple criteria
// Client ID, Date, Status
func Search(database helper.DatabaseX, searchcriteria SearchCriteria) Order {

	// Searching for orders can be a pain
	// Customer ID or Name
	// It can find one or all
	// Order may have a status (open, completed, cancelled)

	database.Collection = "orders"

	ordernull := Order{}

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(database.Database).C(database.Collection)

	// Various search combinations
	// Let's start with Client ID and Status = "Open"

	result := []Order{}
	err1 := c.Find(bson.M{"ClientID": searchcriteria.ClientID}).All(&result)
	if err1 != nil {
		log.Fatal(err1)
	}

	var numrecsel = len(result)

	if numrecsel <= 0 {
		return ordernull
	}

	return result[0]
}

// GetAll works
func GetAll(database helper.DatabaseX) []Order {

	database.Collection = "orders"

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

// Update is for full replacement of order when it's status is Open
// I am not sure if I should protect the other orders
// It uses the order ID to search
func Update(database helper.DatabaseX, orderUpdate Order) helper.Resultado {

	database.Collection = "order"

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Update(bson.M{"ClientID": orderUpdate.ID}, orderUpdate)

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Order has been updated."
	res.IsSuccessful = "Y"

	return res
}

// Delete order
func Delete(database helper.DatabaseX, orderDelete Order) helper.Resultado {

	database.Collection = "orders"

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Remove(bson.M{"ClientID": orderDelete.ID})

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Order deleted"
	res.IsSuccessful = "Y"

	return res
}
