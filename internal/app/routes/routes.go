package routes

import (
	"proundmhee/internal/modules/rsa"
	"proundmhee/internal/modules/vat"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	api := r.Group("/api")

	vat.Register(api.Group("/vat"))
	rsa.Register(api.Group("/rsa"))
}
