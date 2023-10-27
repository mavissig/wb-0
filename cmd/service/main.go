package main

import (
	"log"
	"wb/internal/cache"
	"wb/internal/server"
	"wb/pkg/nats"
)

func init() {
	cache.GetCacheInstance()
	cacheServiceInit()
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
