package generate_code

import (
	"proundmhee/internal/infra/di"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup, deps *di.Deps) {
	svc := NewService()        // concrete
	h := NewHandler(deps, svc) // inject ผ่าน interface
	h.Register(r)
}
