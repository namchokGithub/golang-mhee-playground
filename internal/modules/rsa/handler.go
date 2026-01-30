package rsa

import (
	"net/http"

	"proundmhee/internal/infra/di"
	"proundmhee/internal/shared"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	log *zap.Logger
	svc *Service
}

func NewHandler(deps *di.Deps) *Handler {
	return &Handler{
		log: deps.Logger().Named("vat"),
		svc: NewService(),
	}
}

func (h *Handler) Register(r *gin.RouterGroup) {
	r.POST("/calc", h.toHex)
}

type ToHexRequest struct {
	Text string `json:"text" binding:"required"`
}

func (h *Handler) toHex(c *gin.Context) {
	var req ToHexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("bad_request", zap.Error(err))
		shared.ToFail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	shared.ToSuccess(c, gin.H{"hex": h.svc.ToHex(req.Text)})
}
