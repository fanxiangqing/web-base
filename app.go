package app

import (
	"github.com/fanxiangqing/web-base/internal/lib/utils/env"
	"github.com/fanxiangqing/web-base/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Run(r *gin.Engine, port string) error {
	setBaseRouter(r)

	return r.Run(":" + port)
}

func setBaseRouter(r *gin.Engine) {
	if env.Env.EnvMode == env.EnvModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.LoggingJson())
	r.Use(middleware.Recovering())
}
