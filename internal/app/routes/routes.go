package routes

import (
	"proundmhee/internal/infra/di"
	"proundmhee/internal/modules/rsa"
	"proundmhee/internal/modules/vat"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, deps *di.Deps) {
	api := r.Group("/api")

	vat.Register(api.Group("/vat"), deps)
	rsa.Register(api.Group("/rsa"), deps)
}
