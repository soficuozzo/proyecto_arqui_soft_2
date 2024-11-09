package main

import (
	"context"
	"log"
	"proyecto_arqui_soft_2/cursos-api/clients"
	"proyecto_arqui_soft_2/cursos-api/controllers"
	"proyecto_arqui_soft_2/cursos-api/repositories"
	"proyecto_arqui_soft_2/cursos-api/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Configuración de MongoDB
	mainRepository := repositories.NewMongo(repositories.MongoConfig{
		Host:       "localhost",
		Port:       "27017",
		Username:   "root",
		Password:   "root",
		Database:   "cursos_db",
		Collection: "cursos",
	})

	// Prueba de conexión a MongoDB
	ctx := context.Background()
	if err := mainRepository.TestConnection(ctx); err != nil {
		log.Fatalf("Conexión a MongoDB fallida: %v", err)
	} else {
		log.Println("¡Conexión a MongoDB exitosa!")
	}

	// Configuración de RabbitMQ
	rabbitClient := clients.NewRabbit(clients.RabbitConfig{
		Host:      "localhost",
		Port:      "5672",
		Username:  "user",
		Password:  "password",
		QueueName: "curso_actualizado",
	})
	defer rabbitClient.Close()

	// Crear servicio con RabbitMQ y MongoDB
	cursoService := services.NewCursoService(mainRepository, rabbitClient)

	// Crear controlador de cursos
	cursoController := controllers.NewCursoController(cursoService)

	// Configuración de rutas
	router := gin.Default()
	router.GET("/cursos/:id", cursoController.GetCursoByID)
	router.POST("/cursos", cursoController.Create)
	router.PUT("/cursos/:id", cursoController.Update)
	router.DELETE("/cursos/:id", cursoController.Delete)

	// Inicia el servidor en el puerto 8081
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
