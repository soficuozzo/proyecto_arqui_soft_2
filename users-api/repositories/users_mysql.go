package repositories

import (
	//"errors"
	"fmt"
	"log"
	"proyecto_arqui_soft_2/users-api/dao"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type MySQL struct {
	db *gorm.DB
}

var (
	migrate = []interface{}{
		dao.Usuario{},
	}
)

func NewMySQL(config MySQLConfig) MySQL {
	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Database)

	// Open connection to MySQL using GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %s", err.Error())
	}

	// Automigrate structs to Gorm
	for _, target := range migrate {
		if err := db.AutoMigrate(target); err != nil {
			log.Fatalf("error automigrating structs: %s", err.Error())
		}
	}

	return MySQL{
		db: db,
	}
}
