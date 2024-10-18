package repositories

import (
	"github.com/jinzhu/gorm"
	//"github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var (
	db  *gorm.DB
	err error
)

func init() {
	// DB Connections Paramters
	DBName := "users" //Nombre de la base de datos local de ustedes
	DBUser := "root"          //usuario de la base de datos, habitualmente root
	DBPass := "root"              //password del root en la instalacion
	DBHost := "127.0.0.1"         //host de la base de datos. hbitualmente 127.0.0.1
	// ------------------------

	// abrimos la base de datos
	db, err = gorm.Open("mysql", DBUser+":"+DBPass+"@tcp("+DBHost+":3306)/"+DBName+"?charset=utf8&parseTime=True")

	if err != nil {
		log.Info("Connection Failed to Open")
		log.Fatal(err)
	} else {
		log.Info("Connection Established")
	}

	// We need to add all CLients that we build


}

func StartDbEngine() {
	// We need to migrate all classes model.

	log.Info("Finishing Migration Database Tables")
}
