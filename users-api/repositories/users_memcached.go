package repositories

import (
	//"encoding/json"
	//"errors"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	//"proyecto_arqui_soft_2/users-api/dao"
)

type MemcachedConfig struct {
	Host string
	Port string
}

type Memcached struct {
	client *memcache.Client
}

func usernameKey(username string) string {
	return fmt.Sprintf("username:%s", username)
}

func NewMemcached(config MemcachedConfig) Memcached {
	// Connect to Memcached
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	client := memcache.New(address)

	return Memcached{client: client}
}
