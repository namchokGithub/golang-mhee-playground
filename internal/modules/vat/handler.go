package vat

import (
	"net/http"
	"proundmhee/internal/infra/di"

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
	r.POST("/calc", h.calc)
}

type CalcRequest struct {
	Amount float64 `json:"amount" binding:"required"`
	Rate   float64 `json:"rate"`
}

func (h *Handler) calc(c *gin.Context) {
	var req CalcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("bad_request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}
	if req.Rate == 0 {
		req.Rate = 7
	}

	vat, total := h.svc.Calc(req.Amount, req.Rate)
	h.log.Info("calc", zap.Float64("amount", req.Amount), zap.Float64("rate", req.Rate))

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"vat":    vat,
		"total":  total,
		"amount": req.Amount,
		"rate":   req.Rate,
	})
}
