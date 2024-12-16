package app

import (
	"github.com/Lixerus/auth-service-task/internal/config"
	"github.com/Lixerus/auth-service-task/internal/database"
	"github.com/Lixerus/auth-service-task/internal/models"
	"github.com/Lixerus/auth-service-task/internal/routers"

	"github.com/gin-gonic/gin"
)

func InitDeps() {
	config.LoadEnv()
	config.LoadConfig()
	database.DBInit()
	models.InitModels()
}

func InitApp() *gin.Engine {
	r := gin.Default()
	routers.SetupRoutes(r)
	return r
}
