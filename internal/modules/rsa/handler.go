package rsa

import (
	"net/http"

	"proundmhee/internal/shared"

	"github.com/gin-gonic/gin"
)

type ToHexRequest struct {
	Text string `json:"text" binding:"required"`
}

func Register(r *gin.RouterGroup) {
	r.POST("/to-hex", toHex)
}

func toHex(c *gin.Context) {
	var req ToHexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	shared.OK(c, gin.H{"hex": ToHex(req.Text)})
}
