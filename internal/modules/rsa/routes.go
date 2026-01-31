package rsa

import (
	"proundmhee/internal/infra/di"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup, deps *di.Deps) {
	svc := NewService()
	h := NewHandler(deps, svc)
	h.Register(r)
}
