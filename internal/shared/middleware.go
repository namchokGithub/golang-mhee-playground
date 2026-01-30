package shared

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		lat := time.Since(start)

		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.Printf("%s %s -> %d (%s)", method, path, status, lat)
	}
}
