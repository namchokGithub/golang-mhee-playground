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
	svc Rsa
}

func NewHandler(deps *di.Deps, svc Rsa) *Handler {
	return &Handler{
		log: deps.Logger().Named("vat"),
		svc: svc,
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

	hex, err := h.svc.ToHex(req.Text, false)
	if err != nil {
		h.log.Warn("Internal Server Error", zap.Error(err))
		shared.ToFail(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	shared.ToSuccess(c, gin.H{"hex": hex})
}
