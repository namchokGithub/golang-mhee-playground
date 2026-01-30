package vat

import (
	"net/http"

	"proundmhee/internal/shared"

	"github.com/gin-gonic/gin"
)

type CalcRequest struct {
	Amount float64 `json:"amount" binding:"required"`
	Rate   float64 `json:"rate"` // default 7
}

func Register(r *gin.RouterGroup) {
	r.POST("/calc", calc)
}

func calc(c *gin.Context) {
	var req CalcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	if req.Rate == 0 {
		req.Rate = 7
	}

	vat, total := Calc(req.Amount, req.Rate)
	shared.OK(c, gin.H{
		"amount": req.Amount,
		"rate":   req.Rate,
		"vat":    vat,
		"total":  total,
	})
}
