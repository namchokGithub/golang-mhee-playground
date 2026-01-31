package vat

import (
	"net/http"
	"proundmhee/internal/infra/di"
	"proundmhee/internal/shared"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	log *zap.Logger
	svc Calculator
}

func NewHandler(deps *di.Deps, svc Calculator) *Handler {
	return &Handler{
		log: deps.Logger().Named("vat"),
		svc: svc,
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
		shared.ToFail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	if req.Rate == 0 {
		req.Rate = 7
	}

	vat, total, _ := h.svc.Calculate(req.Amount, req.Rate)
	h.log.Info("calc", zap.Float64("amount", req.Amount), zap.Float64("rate", req.Rate))

	shared.ToSuccess(c, gin.H{
		"ok":     true,
		"vat":    vat,
		"total":  total,
		"amount": req.Amount,
		"rate":   req.Rate,
	})
}
