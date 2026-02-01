package generate_code

import "time"

type CodeGenerator interface {
	GenerateOrderNoGen2WithBranch(now time.Time, branchID, machineID int) string
	GenerateOrderNoGen2(now time.Time) string
	GeneratePaymentRefCode2(now time.Time) string
	GeneratePaymentRequestId12char2(currentTime time.Time) string
	GenerateBranchTransactionNo(now time.Time, branchID, index int, prefix string) string
	GeneratePaymentRefCode2WDF(now time.Time) string
	GeneratePaymentRequestId12char2WDF(currentTime time.Time) string
}
