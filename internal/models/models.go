package models

import (
	"log"

	"github.com/Lixerus/auth-service-task/internal/database"
)

type UserCredentials struct {
	ID                 string `gorm:"primary_key"`
	RefreshToken       string `gorm:"unique"`
	PartialAccessToken string `gorm:"index"`
}

func InitModels() {
	err := database.DB.AutoMigrate(&UserCredentials{})
	if err != nil {
		log.Fatal("Unable to init models" + err.Error())
	}
}
