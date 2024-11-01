package main

import (
    "context"
    "fmt"
    "log"
    "proyecto_arqui_soft_2/cursos-api/repositories"
)

func main() {
    config := repositories.MongoConfig{
        Host:       "localhost",   
        Port:       "27017",       
        Username:   "sofia",  
        Password:   "1234",  
        Database:   "cursos_db",
        Collection: "cursos",
    }

    mongoRepo := repositories.NewMongo(config)
    ctx := context.Background()

    if err := mongoRepo.TestConnection(ctx); err != nil {
        log.Fatalf("MongoDB connection failed: %v", err)
    } else {
        fmt.Println("MongoDB connection test successful!")
    }
}

