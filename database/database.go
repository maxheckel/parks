package database

import (
	"fmt"
	"log"
	"os"

	"github.com/maxheckel/parks/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDb() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "db"
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		dbHost,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migrations")
	db.AutoMigrate(&models.City{})
	db.AutoMigrate(&models.Park{})
	db.AutoMigrate(&models.Tour{})
	db.AutoMigrate(&models.Day{})
	db.AutoMigrate(&models.DayPark{})
	err = db.SetupJoinTable(&models.Day{}, "Parks", &models.DayPark{})
	if err != nil {
		panic(err)
	}
	DB = Dbinstance{
		Db: db,
	}
}
