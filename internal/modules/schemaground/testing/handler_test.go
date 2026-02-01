package schemaground_testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"proundmhee/internal/infra/di"
	"proundmhee/internal/infra/logger"
	"proundmhee/internal/modules/schemaground"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

type fakeComparator struct {
	resp schemaground.CompareResponse
	err  error
}

func (f fakeComparator) Compare(schema string) (schemaground.CompareResponse, error) {
	return f.resp, f.err
}

func TestHandler_Compare_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	log, _ := logger.New(logger.Config{
		Env:   "dev",
		Level: zapcore.WarnLevel,
	})
	deps := &di.Deps{Log: log}

	svc := fakeComparator{
		resp: schemaground.CompareResponse{
			Schema: "public",
			AtoB: schemaground.SchemaDiffSummary{
				DBA: "A",
				DBB: "B",
				TableDiffs: []schemaground.TableDiff{
					{Table: "users", OnlyInA: true},
				},
			},
			AtoC: schemaground.SchemaDiffSummary{
				DBA:        "A",
				DBB:        "C",
				TableDiffs: nil,
			},
		},
	}
	h := schemaground.NewHandler(deps, svc)
	grp := r.Group("/api/schemaground")
	h.Register(grp)

	body := map[string]any{"schema": "public"}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/schemaground/compare", bytes.NewReader(b))
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
}

func TestHandler_Compare_BadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	log, _ := logger.New(logger.Config{
		Env:   "dev",
		Level: zapcore.WarnLevel,
	})
	deps := &di.Deps{Log: log}

	svc := fakeComparator{}
	h := schemaground.NewHandler(deps, svc)
	grp := r.Group("/api/schemaground")
	h.Register(grp)

	req := httptest.NewRequest(http.MethodPost, "/api/schemaground/compare", bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
