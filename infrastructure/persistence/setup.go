package persistence

import (
	"fmt"
	"github.com/akwanmaroso/ddd-drugs/domain/entity"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func DBConn() (*gorm.DB, error) {
	var err error
	err = godotenv.Load(os.ExpandEnv("./../../.env"))
	if err != nil {
		log.Fatalf("error getting env %v\n", err)
	}
	return LocalDatabase()
}

func LocalDatabase() (*gorm.DB, error) {
	dbDriver := os.Getenv("TEST_DB_DRIVER")
	dbHost := os.Getenv("TEST_DB_HOST")
	dbUser := os.Getenv("TEST_DB_USER")
	dbPassword := os.Getenv("TEST_DB_PASSWORD")
	dbName := os.Getenv("TEST_DB_NAME")
	dbPort := os.Getenv("TEST_DB_PORT")

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
	conn, err := gorm.Open(dbDriver, DBURL)
	if err != nil {
		return nil, err
	} else {
		log.Println("CONNECTED TO: ", dbDriver)
	}

	err = conn.DropTableIfExists(&entity.User{}, &entity.Drug{}).Error
	if err != nil {
		return nil, err
	}

	err = conn.Debug().AutoMigrate(entity.User{}, entity.Drug{}).Error
	if err != nil {
		return nil, err
	}

	return conn, nil
}
