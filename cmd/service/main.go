package main

import (
	"log"
	"wb/internal/cache"
	"wb/internal/database"
	"wb/internal/model"
	"wb/internal/server"
	"wb/pkg/nats"
)

func init() {
	var (
		db     database.Database
		orders = make(map[string]model.Order)
	)
	db.Connect()
	orders, _ = db.GetAllOrders()
	cache.GetCacheInstance(orders)
	cacheServiceInit()
	dbServiceInit()
}

func main() {
	server.RunServer()
}

func cacheServiceInit() {
	var srv nats.Service
	err := srv.Connect("consumer_cache")
	if err != nil {
		log.Fatal("cache service connect: ", err)
	}

	err = srv.Subscribe("uploadFile", cache.GetCacheInstance())
	if err != nil {
		log.Fatal("cache service subscribe: ", err)
	}
}

func dbServiceInit() {
	var (
		db  database.Database
		srv nats.Service
	)

	err := srv.Connect("consumer_db")
	if err != nil {
		log.Fatal("db service connect: ", err)
	}

	err = srv.Subscribe("uploadFile", &db)
}
