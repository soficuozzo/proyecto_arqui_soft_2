package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"proyecto_arqui_soft_2/search-api/clients/queues"
	controllers "proyecto_arqui_soft_2/search-api/controllers"
	repositories "proyecto_arqui_soft_2/search-api/repositories"
	services "proyecto_arqui_soft_2/search-api/services"
)

func main() {
	// Solr
	solrRepo := repositories.NewSolr(repositories.SolrConfig{
		Host:       "localhost",   // Solr host
		Port:       "8983",   // Solr port
		Collection: "cursos", // Collection name
	})

	// Rabbit
	eventsQueue := queues.NewRabbit(queues.RabbitConfig{
		Host:      "localhost",
		Port:      "5672",
		Username:  "user",
		Password:  "password",
		QueueName: "cursos-news",
	})

	// Cursos API
	cursosAPI := repositories.NewHTTP(repositories.HTTPConfig{
		Host: "cursos-api",
		Port: "8081",
	})

	// Services
	service := services.NewService(solrRepo, cursosAPI)

	// Controllers
	controller := controllers.NewController(service)

	// Launch rabbit consumer
	if err := eventsQueue.StartConsumer(service.HandleCursoNew); err != nil {
		log.Fatalf("Error running consumer: %v", err)
	}

	// Create router
	router := gin.Default()
	router.GET("/search", controller.Search)
	if err := router.Run(":8082"); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}