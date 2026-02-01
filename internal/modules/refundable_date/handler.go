package refundable_date

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
	svc Checker
}

func NewHandler(deps *di.Deps, svc Checker) *Handler {
	return &Handler{
		log: deps.Logger().Named("refundable_date"),
		svc: svc,
	}
}

func (h *Handler) Register(r *gin.RouterGroup) {
	r.POST("/check", h.check)
}

type CheckRequest struct {
	Date  string `json:"date" binding:"required"`
	Today string `json:"today"`
}

func (h *Handler) check(c *gin.Context) {
	var req CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("bad_request", zap.Error(err))
		shared.ToFail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	date, err := time.ParseInLocation("2006-01-02", req.Date, time.Local)
	if err != nil {
		h.log.Warn("bad_request", zap.Error(err))
		shared.ToFail(c, http.StatusBadRequest, "BAD_REQUEST", "invalid date format")
		return
	}

	today := time.Now()
	if req.Today != "" {
		today, err = time.ParseInLocation("2006-01-02", req.Today, time.Local)
		if err != nil {
			h.log.Warn("bad_request", zap.Error(err))
			shared.ToFail(c, http.StatusBadRequest, "BAD_REQUEST", "invalid today format")
			return
		}
	}

	refundable := h.svc.CheckRefundable(date, today)
	shared.ToSuccess(c, gin.H{
		"ok":         true,
		"refundable": refundable,
		"date":       date.Format("2006-01-02"),
		"today":      today.Format("2006-01-02"),
	})
}
