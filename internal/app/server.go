package app

import (
	"net/http"

	"proundmhee/internal/app/routes"
	"proundmhee/internal/infra/di"
	"proundmhee/internal/infra/logger"
	"proundmhee/internal/shared"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

func NewServer(env string) (*gin.Engine, *di.Deps, error) {
	r := gin.New()

	// logger
	log, _ := logger.New(logger.Config{
		Env:   "dev",
		Level: zapcore.WarnLevel,
	})

	// รวม dependencies ของแอป (logger, db, config)
	deps := &di.Deps{Log: log}

	// middleware พื้นฐาน
	r.Use(logger.GinMiddleware(log))
	r.Use(shared.RequestLogger())
	r.Use(gin.Recovery())
	r.Use(logger.GinMiddleware(log))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"I'm healthy": true})
	})

	// รวม routes ของทุก module
	routes.RegisterRoutes(r, deps)

	return r, deps, nil
}
