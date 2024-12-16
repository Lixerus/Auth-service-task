package routers

import (
	"github.com/Lixerus/auth-service-task/internal/middleware"
	"github.com/Lixerus/auth-service-task/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/credentials", services.GetCredentials)
	r.POST("/refresh", middleware.RequireCookies, services.RefreshAuthToken)
}
