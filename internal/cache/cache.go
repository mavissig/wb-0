package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"wb/internal/model"
)

type TypeCache map[string]model.Order

var cache__ TypeCache

var syn sync.Once

func GetCacheInstance(params ...map[string]model.Order) TypeCache {
	syn.Do(func() {
		cache__ = make(TypeCache, 100)
		if len(params) > 0 && len(params[0]) != 0 {
			cache__ = params[0]
		}
		fmt.Println(len(cache__))
	})
	return cache__
}

func (c TypeCache) Consume(data []byte) error {
	var (
		order model.Order
		cache = make(TypeCache)
	)

	err := json.Unmarshal(data, &order)
	if err != nil {
		log.Println("unmarshal error: ", err)
	}

	cache = GetCacheInstance()
	cache[order.OrderUID] = order

	fmt.Println("[cache add]: ", order.OrderUID)
	return nil
}
