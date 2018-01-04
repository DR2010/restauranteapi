// Package btcmarkets is a dish for packages
// -------------------------------------
// .../restauranteapi/dishes/btcmarkets.go
// -------------------------------------
package btcmarkets

import (
	"fmt"
	"log"
	helper "restauranteapi/helper"
	"strconv"

	"github.com/go-redis/redis"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// BalanceCrypto is to be exported
type BalanceCrypto struct {
	Balance        string // balance
	Currency       string // Currency
	CotacaoAtual   string // Cotacao
	ValueInCashAUD string // Value in AUD
	DateTime       string // date time
}

// CryptoCotacaoAdd is for export
func CryptoCotacaoAdd(redisclient *redis.Client, cryptoInsert BalanceCrypto) helper.Resultado {

	database := new(helper.DatabaseX)
	database.Collection = "btccotacao"
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

	err = collection.Insert(cryptoInsert)

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Crypto Cotacao added"
	res.IsSuccessful = "Y"

	return res
}

// GetAll works
func GetAll(redisclient *redis.Client, currency string, rows string) []BalanceCrypto {

	rowsi, err := strconv.Atoi(rows)
	if rowsi == 0 {
		rowsi = 100
	}

	database := new(helper.DatabaseX)

	database.Collection = "btccotacao"

	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	// database.Database = "restaurante"
	// database.Location = "192.168.2.180"

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

	var results []BalanceCrypto

	// err1 := c.Find(bson.M{"name": dishName}).All(&result)
	// db.getCollection('btccotacao').find({}).sort({_id:-1}).limit(15);
	// err = c.Find(bson.M{"currency": currency}).All(&results)
	// err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&results)
	err = c.Find(bson.M{"currency": currency}).Sort("-_id").Limit(rowsi).All(&results)
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
