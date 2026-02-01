package generate_code_testing

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"proundmhee/internal/modules/generate_code"

	"github.com/stretchr/testify/require"
)

func TestService_GenerateOrderNoGen2WithBranch(t *testing.T) {
	svc := generate_code.NewService()
	now := time.Date(2026, time.February, 1, 14, 5, 30, 123456000, time.UTC)

	got := svc.GenerateOrderNoGen2WithBranch(now, 123, 45)

	require.Len(t, got, 30)
	require.True(t, regexp.MustCompile(`^OT[0-9]{6}[0-9]{4}[A-Z]{1}[0-9]{5}[0-9]{6}[0-9]{6}$`).MatchString(got))
	require.Equal(t, "OT2602011405", got[:12])
}

func TestService_GenerateOrderNoGen2(t *testing.T) {
	svc := generate_code.NewService()
	now := time.Date(2026, time.February, 1, 9, 4, 5, 987654000, time.UTC)

	got := svc.GenerateOrderNoGen2(now)

	require.Len(t, got, 20)
	require.True(t, regexp.MustCompile(`^OT[0-9]{6}[0-9]{4}[A-Z]{1}[0-9]{7}$`).MatchString(got))
	require.Equal(t, "OT2602010904", got[:12])
}

func TestService_GeneratePaymentRefCode2(t *testing.T) {
	svc := generate_code.NewService()
	now := time.Date(2026, time.February, 1, 9, 4, 5, 123456000, time.UTC)

	got := svc.GeneratePaymentRefCode2(now)

	require.Len(t, got, 20)
	require.True(t, regexp.MustCompile(`^OTR[0-9]{6}[0-9]{4}[A-Z]{3}[0-9]{4}$`).MatchString(got))
	require.Equal(t, "OTR2602010904", got[:13])
}

func TestService_GeneratePaymentRequestId12char2(t *testing.T) {
	svc := generate_code.NewService()
	now := time.Date(2026, time.February, 1, 9, 4, 5, 654321000, time.UTC)

	got := svc.GeneratePaymentRequestId12char2(now)

	require.Len(t, got, 12)
	require.True(t, regexp.MustCompile(`^OTR[0-9]{6}[A-Z]{3}$`).MatchString(got))
}

func TestService_GenerateBranchTransactionNo(t *testing.T) {
	svc := generate_code.NewService()
	now := time.Date(2026, time.February, 1, 14, 5, 30, 123000000, time.UTC)

	got := svc.GenerateBranchTransactionNo(now, 12, 7, "C")

	timestamp := now.UnixMilli()
	trimTimestamp := fmt.Sprintf("%013d", timestamp)[8:13]
	expected := fmt.Sprintf("C%04d26%02d%02d%02d%02d%s%02d", 12, now.Month(), now.Day(), now.Hour(), now.Minute(), trimTimestamp, 7)

	require.Equal(t, expected, got)
}

func TestService_GeneratePaymentRefCode2WDF(t *testing.T) {
	svc := generate_code.NewService()
	now := time.Date(2026, time.February, 1, 9, 4, 5, 123456000, time.UTC)

	got := svc.GeneratePaymentRefCode2WDF(now)

	require.Len(t, got, 20)
	require.True(t, regexp.MustCompile(`^WDR[0-9]{6}[0-9]{4}[A-Z]{3}[0-9]{4}$`).MatchString(got))
	require.Equal(t, "WDR2602010904", got[:13])
}

func TestService_GeneratePaymentRequestId12char2WDF(t *testing.T) {
	svc := generate_code.NewService()
	now := time.Date(2026, time.February, 1, 9, 4, 5, 654321000, time.UTC)

	got := svc.GeneratePaymentRequestId12char2WDF(now)

	require.Len(t, got, 12)
	require.True(t, regexp.MustCompile(`^WDR[0-9]{6}[A-Z]{3}$`).MatchString(got))
}
