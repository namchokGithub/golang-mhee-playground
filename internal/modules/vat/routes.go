package vat

import (
	"proundmhee/internal/infra/di"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup, deps *di.Deps) {
	h := NewHandler(deps)
	h.Register(r)
}
