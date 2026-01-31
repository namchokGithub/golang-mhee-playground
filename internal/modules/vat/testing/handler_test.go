package vat_testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"proundmhee/internal/infra/di"
	"proundmhee/internal/infra/logger"
	"proundmhee/internal/modules/vat"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

type fakeCalc struct {
	vat   float64
	total float64
	err   error
}

func (f fakeCalc) Calculate(amount, rate float64) (float64, float64, error) {
	return f.vat, f.total, f.err
}

func TestHandler_Calc_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// router
	r := gin.New()

	// deps ที่จำเป็น
	// logger
	log, _ := logger.New(logger.Config{
		Env:   "dev",
		Level: zapcore.WarnLevel,
	})
	// รวม dependencies ของแอป (logger, db, config)
	deps := &di.Deps{Log: log}

	svc := fakeCalc{vat: 70, total: 1070, err: nil}
	h := vat.NewHandler(deps, svc)

	grp := r.Group("/api/vat")
	h.Register(grp)

	// request
	body := map[string]any{"amount": 1000, "rate": 7}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/vat/calc", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// serve
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	// assert response (ถ้าคุณมีรูปแบบ response กลางก็ parse ได้)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, true, resp["ok"])
}

func TestHandler_Calc_BadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// deps ที่จำเป็น
	// logger
	log, _ := logger.New(logger.Config{
		Env:   "dev",
		Level: zapcore.WarnLevel,
	})
	// รวม dependencies ของแอป (logger, db, config)
	deps := &di.Deps{Log: log}
	svc := fakeCalc{}
	h := vat.NewHandler(deps, svc)

	grp := r.Group("/api/vat")
	h.Register(grp)

	req := httptest.NewRequest(http.MethodPost, "/api/vat/calc", bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
