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
			Host:     "127.0.0.1",
			Port:     "3306",
			Database: "users-api",
			Username: "root",
			Password: "root1234",
		},
	)

	// Cache
	cacheRepo := repositories.NewCache(repositories.CacheConfig{
		TTL: 1 * time.Minute,
	})

	// Memcached
	memcachedRepo := repositories.NewMemcached(repositories.MemcachedConfig{
		Host: "localhost",
		Port: "11211",
	})

	// Services

	service := services.NewService(mySQLRepo, cacheRepo, memcachedRepo)
	// Handlers
	controlleruser := controllers.NewController(service)

	// Create router
	router := gin.Default()

	// URL mappings
	router.GET("/users/:id", controlleruser.GetUsuariobyID)
	router.GET("/users/email/:email", controlleruser.GetUsuariobyEmail)
	router.POST("/users/create", controlleruser.CrearUsuario)
	router.POST("/login", controlleruser.Login)

	// Run application
	if err := router.Run(":8080"); err != nil {
		log.Panicf("Error running application: %v", err)
	}

}
