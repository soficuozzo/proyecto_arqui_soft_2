package main

import (
	"context"
	"log"

	"proyecto_arqui_soft_2/cursos-api/controllers"
	"proyecto_arqui_soft_2/cursos-api/repositories"
	"proyecto_arqui_soft_2/cursos-api/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Configuración del repositorio principal
	mainRepository := repositories.NewMongo(repositories.MongoConfig{
		Host:       "localhost",
		Port:       "27017",
		Username:   "root",
		Password:   "root",
		Database:   "cursos_db",
		Collection: "cursos",
	})

	InscripcionRepository := repositories.NewMySQL(
		repositories.MySQLConfig{
			Host:     "127.0.0.1",
			Port:     "3306",
			Database: "users-api",
			Username: "root",
			Password: "root1234",
		},
	)

	// Prueba de conexión a MongoDB
	ctx := context.Background()
	if err := mainRepository.TestConnection(ctx); err != nil {
		log.Fatalf("Conexión a MongoDB fallida: %v", err)
	} else {
		log.Println("¡Conexión a MongoDB exitosa!")
	}

	cursoService := services.NewCursoService(mainRepository, InscripcionRepository)

	// Creación del controlador de cursos
	cursoController := controllers.NewCursoController(cursoService)

	// Configuración de rutas
	router := gin.Default()
	router.GET("/cursos/:id", cursoController.GetCursoByID)
	router.POST("/cursos", cursoController.Create)
	router.PUT("/cursos/:id", cursoController.Update)
	router.DELETE("/cursos/:id", cursoController.Delete)
	router.POST("/inscripcion", cursoController.CrearInscripcion)
	router.GET("/usuario/miscursos/:id", cursoController.GetInscripcionByUserId)
	// Inicia el servidor en el puerto 8081
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
