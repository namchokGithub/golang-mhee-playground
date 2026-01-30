package app

import (
	"net/http"

	"proundmhee/internal/app/routes"
	"proundmhee/internal/shared"

	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {
	r := gin.New()

	// middleware พื้นฐาน
	r.Use(shared.RequestLogger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"I'm healthy": true})
	})

	// รวม routes ของทุก module
	routes.Register(r)

	return r
}
