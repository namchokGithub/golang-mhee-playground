package generate_code_testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"proundmhee/internal/infra/di"
	"proundmhee/internal/infra/logger"
	"proundmhee/internal/modules/generate_code"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

type fakeGenerator struct {
	code string
}

func (f fakeGenerator) GenerateOrderNoGen2WithBranch(now time.Time, branchID, machineID int) string {
	return ""
}

func (f fakeGenerator) GenerateOrderNoGen2(now time.Time) string { return "" }

func (f fakeGenerator) GeneratePaymentRefCode2(now time.Time) string { return "" }

func (f fakeGenerator) GeneratePaymentRequestId12char2(currentTime time.Time) string { return "" }

func (f fakeGenerator) GenerateBranchTransactionNo(now time.Time, branchID, index int, prefix string) string {
	return f.code
}

func (f fakeGenerator) GeneratePaymentRefCode2WDF(now time.Time) string { return "" }

func (f fakeGenerator) GeneratePaymentRequestId12char2WDF(currentTime time.Time) string {
	return ""
}

func TestHandler_Generate_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	log, _ := logger.New(logger.Config{
		Env:   "dev",
		Level: zapcore.WarnLevel,
	})
	deps := &di.Deps{Log: log}

	svc := fakeGenerator{code: "TXN-123"}
	h := generate_code.NewHandler(deps, svc)
	grp := r.Group("/api/generate")
	h.Register(grp)

	req := httptest.NewRequest(http.MethodGet, "/api/generate/generate", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, true, resp["ok"])

	data, ok := resp["data"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, true, data["ok"])
	require.Equal(t, "TXN-123", data["code"])
}

func TestHandler_Generate_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	log, _ := logger.New(logger.Config{
		Env:   "dev",
		Level: zapcore.WarnLevel,
	})
	deps := &di.Deps{Log: log}

	svc := fakeGenerator{code: ""}
	h := generate_code.NewHandler(deps, svc)
	grp := r.Group("/api/generate")
	h.Register(grp)

	req := httptest.NewRequest(http.MethodGet, "/api/generate/generate", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, false, resp["ok"])

	errObj, ok := resp["error"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "INTERNAL_SERVER_ERROR", errObj["code"])
}
