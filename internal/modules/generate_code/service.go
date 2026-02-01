package generate_code

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) GenerateOrderNoGen2WithBranch(now time.Time, branchID, machineID int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 1)
	for i := range b {
		b[i] = letters[int(rand.Int31n(int32(len(letters))))]
	}
	t := now.Format("05.000000")
	number := strings.ReplaceAll(t, ".", "")
	randStr := strconv.Itoa(rand.Intn(9)+1) + number[3:]
	result, _ := strconv.ParseInt(randStr, 10, 64)
	str := strconv.Itoa(int(result))
	orderNo := fmt.Sprintf("OT%s%s%s%05d%06d%s", now.Format("060102"), now.Format("1504"), string(b), branchID, machineID, str)
	return orderNo
}

func (s *Service) GenerateOrderNoGen2(now time.Time) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 1)
	for i := range b {
		b[i] = letters[int(rand.Int31n(int32(len(letters))))]
	}
	t := now.Format("05.000000")
	number := strings.ReplaceAll(t, ".", "")
	randStr := strconv.Itoa(rand.Intn(9)+1) + number[2:]
	result, _ := strconv.ParseInt(randStr, 10, 64)
	str := strconv.Itoa(int(result))
	orderNo := fmt.Sprintf("OT%s%s%s%s", now.Format("060102"), now.Format("1504"), string(b), str)
	return orderNo
}

func (s *Service) GeneratePaymentRefCode2(now time.Time) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 3)
	for i := range b {
		b[i] = letters[int(rand.Int31n(int32(len(letters))))]
	}
	nanoPart := fmt.Sprintf("%06d", now.UnixNano()/int64(time.Microsecond)%1_000_000)
	return "OTR" + now.Format("060102") + now.Format("1504") + string(b) + nanoPart[2:]
}

func (s *Service) GeneratePaymentRequestId12char2(currentTime time.Time) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 3)
	for i := range b {
		b[i] = letters[int(rand.Int31n(int32(len(letters))))]
	}
	t := currentTime.Format("150405.0000")
	number := strings.ReplaceAll(t, ".", "")
	randStr := strconv.Itoa(rand.Intn(9)+1) + number[1:]
	result, _ := strconv.ParseInt(randStr, 10, 64)
	str := strconv.Itoa(int(result))
	txRef := "OTR" + str[:6] + string(b)
	return txRef
}

// prefix/รหัสสาขา/ปี/เดือน/วัน/ชั่วโมง/นาที/UnixMilli
func (s *Service) GenerateBranchTransactionNo(now time.Time, branchID, index int, prefix string) string {
	timestamp := now.UnixMilli()
	trimTimestamp := strconv.Itoa(int(timestamp))[8:13]
	yearStr := strconv.Itoa(now.Year())
	last2DigitsOfYear := yearStr[len(yearStr)-2:]
	tranNo := fmt.Sprintf("%s%04d%s%02d%02d%02d%02d%s", prefix, branchID, last2DigitsOfYear, now.Month(), now.Day(), now.Hour(), now.Minute(), trimTimestamp)
	return fmt.Sprintf("%s%02d", tranNo, index)
}

func (s *Service) GeneratePaymentRefCode2WDF(now time.Time) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 3)
	for i := range b {
		b[i] = letters[int(rand.Int31n(int32(len(letters))))]
	}
	nanoPart := fmt.Sprintf("%06d", now.UnixNano()/int64(time.Microsecond)%1_000_000)
	return "WDR" + now.Format("060102") + now.Format("1504") + string(b) + nanoPart[2:]
}

func (s *Service) GeneratePaymentRequestId12char2WDF(currentTime time.Time) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 3)
	for i := range b {
		b[i] = letters[int(rand.Int31n(int32(len(letters))))]
	}
	t := currentTime.Format("150405.0000")
	number := strings.ReplaceAll(t, ".", "")
	randStr := strconv.Itoa(rand.Intn(9)+1) + number[1:]
	result, _ := strconv.ParseInt(randStr, 10, 64)
	str := strconv.Itoa(int(result))
	txRef := "WDR" + str[:6] + string(b)
	return txRef
}
