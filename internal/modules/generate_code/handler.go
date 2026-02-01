package generate_code

import (
	"net/http"
	"proundmhee/internal/infra/di"
	"proundmhee/internal/shared"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	log *zap.Logger
	svc CodeGenerator
}

func NewHandler(deps *di.Deps, svc CodeGenerator) *Handler {
	return &Handler{
		log: deps.Logger().Named("generate_code"),
		svc: svc,
	}
}

func (h *Handler) Register(r *gin.RouterGroup) {
	r.GET("", h.generate)
}

func (h *Handler) generate(c *gin.Context) {
	code := h.svc.GenerateBranchTransactionNo(time.Now(), 123, 23333, "C")

	if code == "" {
		h.log.Warn("Internal Server Error")
		shared.ToFail(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal Server Error")
		return
	}

	shared.ToSuccess(c, gin.H{
		"ok":   true,
		"code": code,
	})
}
