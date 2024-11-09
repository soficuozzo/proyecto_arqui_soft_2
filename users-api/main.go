package main

import (
	"proyecto_arqui_soft_2/users-api/app"
	"proyecto_arqui_soft_2/users-api/repositories"
)

func main() {
	repositories.StartDbEngine()
	app.StartRoute()

}
