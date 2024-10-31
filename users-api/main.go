package main

import (
	"log"
	controllers "proyecto_arqui_soft_2/users-api/controllers"
	repositories "proyecto_arqui_soft_2/users-api/repositories"
	services "proyecto_arqui_soft_2/users-api/services"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// MySQL
	mySQLRepo := repositories.NewMySQL(
		repositories.MySQLConfig{
			Host:     "mysql",
			Port:     "3306",
			Database: "users-api",
			Username: "root",
			Password: "root",
		},
	)

	// Cache
	cacheRepo := repositories.NewCache(repositories.CacheConfig{
		TTL: 1 * time.Minute,
	})

	// Memcached
	memcachedRepo := repositories.NewMemcached(repositories.MemcachedConfig{
		Host: "memcached",
		Port: "11211",
	})

	// Services
	service := services.NewService(mySQLRepo, cacheRepo, memcachedRepo, jwtTokenizer)

	// Handlers
	controller := controllers.NewController(service)

	// Create router
	router := gin.Default()

	// URL mappings

	// Run application
	if err := router.Run(":8080"); err != nil {
		log.Panicf("Error running application: %v", err)
	}

}
