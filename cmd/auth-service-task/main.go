package main

import (
	"github.com/Lixerus/auth-service-task/internal/app"
)

func main() {
	app.InitDeps()
	app.InitApp().Run("0.0.0.0:8080")
}
