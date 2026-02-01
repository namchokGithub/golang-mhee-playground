package schemaground

import (
	"net/http"
	"proundmhee/internal/infra/di"
	"proundmhee/internal/shared"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	log *zap.Logger
	svc Comparator
}

func NewHandler(deps *di.Deps, svc Comparator) *Handler {
	return &Handler{
		log: deps.Logger().Named("schemaground"),
		svc: svc,
	}
}

func (h *Handler) Register(r *gin.RouterGroup) {
	r.POST("/compare", h.compare)
}

type CompareRequest struct {
	Schema string `json:"schema"`
}

func (h *Handler) compare(c *gin.Context) {
	var req CompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("bad_request", zap.Error(err))
		shared.ToFail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	resp, err := h.svc.Compare(req.Schema)
	if err != nil {
		h.log.Warn("Internal Server Error", zap.Error(err))
		shared.ToFail(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	shared.ToSuccess(c, gin.H{
		"ok":     true,
		"result": resp,
	})
}
