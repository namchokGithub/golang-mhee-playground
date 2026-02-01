package refundable_date_testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"proundmhee/internal/infra/di"
	"proundmhee/internal/infra/logger"
	"proundmhee/internal/modules/refundable_date"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

type fakeChecker struct {
	result bool
}

func (f fakeChecker) CheckRefundable(date, today time.Time) bool {
	return f.result
}

func TestHandler_Check_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	log, _ := logger.New(logger.Config{
		Env:   "dev",
		Level: zapcore.WarnLevel,
	})
	deps := &di.Deps{Log: log}

	svc := fakeChecker{result: true}
	h := refundable_date.NewHandler(deps, svc)
	grp := r.Group("/api/refundable")
	h.Register(grp)

	body := map[string]any{
		"date":  "2025-09-02",
		"today": "2025-09-02",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/refundable/check", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, true, resp["ok"])

	data, ok := resp["data"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, true, data["ok"])
	require.Equal(t, true, data["refundable"])
	require.Equal(t, "2025-09-02", data["date"])
	require.Equal(t, "2025-09-02", data["today"])
}

func TestHandler_Check_BadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	log, _ := logger.New(logger.Config{
		Env:   "dev",
		Level: zapcore.WarnLevel,
	})
	deps := &di.Deps{Log: log}

	svc := fakeChecker{result: false}
	h := refundable_date.NewHandler(deps, svc)
	grp := r.Group("/api/refundable")
	h.Register(grp)

	req := httptest.NewRequest(http.MethodPost, "/api/refundable/check", bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Check_InvalidDateFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	log, _ := logger.New(logger.Config{
		Env:   "dev",
		Level: zapcore.WarnLevel,
	})
	deps := &di.Deps{Log: log}

	svc := fakeChecker{result: false}
	h := refundable_date.NewHandler(deps, svc)
	grp := r.Group("/api/refundable")
	h.Register(grp)

	body := map[string]any{
		"date":  "2025/09/02",
		"today": "2025-09-02",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/refundable/check", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
