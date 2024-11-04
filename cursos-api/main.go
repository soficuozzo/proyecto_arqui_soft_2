package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"proyecto_arqui_soft_2/cursos-api/controllers"
	"proyecto_arqui_soft_2/cursos-api/repositories"
	"proyecto_arqui_soft_2/cursos-api/services"
)

func main() {
	// Configuración del repositorio principal
	mainRepository := repositories.NewMongo(repositories.MongoConfig{
		Host:       "localhost",
		Port:       "27017",
		Username:   "sofia",
		Password:   "1234",
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


	cursoService := services.NewCursoService(mainRepository)

	// Creación del controlador de cursos
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
