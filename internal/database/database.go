package database

import (
	"fmt"
	"log"

	"github.com/Lixerus/auth-service-task/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBInit() {
	cfg := config.DBConfig
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.HOST, cfg.USERNAME, cfg.PASSWORD, cfg.NAME, cfg.PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database" + err.Error())
	}
	DB = db
}
