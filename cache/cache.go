package cache

import (
	"strconv"

	"github.com/pufferpanel/apufferi/cache"
	"github.com/pufferpanel/apufferi/config"
)

func CreateCache() *cache.MemoryCache {
	capacity, err := strconv.Atoi(config.Get("console-buffer"))
	if err != nil {
		capacity = 50
	}
	return &cache.MemoryCache{
		Buffer:   make([]cache.Message, 0),
		Capacity: capacity,
	}
}
