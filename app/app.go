package app

import (
	_ "github.com/fanxiangqing/web-base/internal/lib/logger"
	"github.com/fanxiangqing/web-base/internal/lib/utils/env"
	"github.com/fanxiangqing/web-base/internal/middleware"
	"github.com/gin-gonic/gin"
)

type App struct {
	Port   string
	Engine *gin.Engine
}

func NewApp(port string) *App {
	r := gin.New()
	setBaseRouter(r)
	return &App{
		Engine: r,
		Port:   port,
	}
}

func (app *App) Run(port string) error {
	return app.Engine.Run(":" + port)
}

func setBaseRouter(r *gin.Engine) {
	if env.Env.EnvMode == env.EnvModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.LoggingJson())
	r.Use(middleware.Recovering())
}
