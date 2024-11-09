package config

import "time"

const (
	MySQLHost     = "127.0.0.1"
	MySQLPort     = "3306"
	MySQLDatabase = "users-api"
	MySQLUsername = "root"
	MySQLPassword = "root1234"

	CacheDuration = 30 * time.Second

	MemcachedHost = "localhost"
	MemcachedPort = "11211"
)
