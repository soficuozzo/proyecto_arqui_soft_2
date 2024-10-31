package repositories

import (
	//"fmt"
	"time"

	"github.com/karlseguin/ccache"
	//"proyecto_arqui_soft_2/users-api/dao"
)

type CacheConfig struct {
	TTL time.Duration // Cache expiration time
}

type Cache struct {
	client *ccache.Cache
	ttl    time.Duration
}

func NewCache(config CacheConfig) Cache {
	// Initialize ccache with default settings
	cache := ccache.New(ccache.Configure())
	return Cache{
		client: cache,
		ttl:    config.TTL,
	}
}
