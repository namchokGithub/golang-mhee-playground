package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func main() {
	for i := 0; i < 2; i++ {
		now := time.Now()
		fmt.Println(GenerateOrderNoGen2WithBranch(now, 123+i, 23333+i), ":len(", len(GenerateOrderNoGen2WithBranch(now, 123+i, 23333+i)), ")")
		fmt.Println(GenerateOrderNoGen2(now), ":len(", len(GenerateOrderNoGen2(now)), ")")
		fmt.Println(GeneratePaymentRefCode2(now), ":len(", len(GeneratePaymentRefCode2(now)), ")")
		fmt.Println(GeneratePaymentRequestId12char2(now), ":len(", len(GeneratePaymentRequestId12char2(now)), ")")

		fmt.Println(GenerateBranchTransactionNo(now, 123, 10, "C"), ":len(", len(GenerateBranchTransactionNo(now, 123, 10, "C")), ")")
		fmt.Println(GenerateBranchTransactionNo(now, 123, i, "WDF"), ":len(", len(GenerateBranchTransactionNo(now, 123, i, "WDF")), ")")
		fmt.Println(GeneratePaymentRefCode2WDF(now), ":len(", len(GeneratePaymentRefCode2WDF(now)), ")")
		fmt.Println(GeneratePaymentRequestId12char2WDF(now), ":len(", len(GeneratePaymentRequestId12char2WDF(now)), ")")
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	}

	// isDuplicated := make(map[string]bool)
	// for i := 0; i < 5; i++ {
	// 	go func() {
	// 		now := time.Now()
	// 		// no := GenerateOrderNo3(now, 123, 23333)
	// 		// _, exists := isDuplicated[no]
	// 		// if exists {
	// 		// 	isDuplicated[no] = true
	// 		// }
	// 		// fmt.fmt.Println(no)

	// 		fmt.Println(GenerateOrderNo3(now, 123, 23333), ":len(", len(GenerateOrderNo3(now, 123, 23333)), ")")
	// 		fmt.Println(GeneratePaymentRequestId12char2(now), ":len(", len(GeneratePaymentRequestId12char2(now)), ")")
	// 		fmt.Println(GeneratePaymentRefCode2(now), ":len(", len(GeneratePaymentRefCode2(now)), ")")
	// 		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	// 	}()
	// }
	// time.Sleep(3 * time.Second)
	// count := 0
	// for k, v := range isDuplicated {
	// 	if v {
	// 		fmt.fmt.Println(k, " is duplicated")
	// 		count++
	// 	}
	// }
	// fmt.fmt.Println("done: ", count)
}

// func GenerateOrderNo(now time.Time) string {
// 	nanoPart := fmt.Sprintf("%06d", (now.UnixNano()/int64(time.Nanosecond))%1_000_000)
// 	orderNo := "OT" + now.Format("060102") + nanoPart + now.Format("1504")
// 	return orderNo
// }

// func GenerateOrderNo1(now time.Time) string {
// 	nanoPart := fmt.Sprintf("%06d", now.UnixNano()/int64(time.Microsecond)%1_000_000)
// 	orderNo := "OT" + now.Format("060102") + nanoPart + now.Format("1504")
// 	return orderNo
// }

// func GenerateOrderNo2(now time.Time, branchID, machineID uint) string {
// 	t := now.Format("05.000000")
// 	number := strings.ReplaceAll(t, ".", "")
// 	randStr := strconv.Itoa(rand.Intn(9)+1) + number[6:]
// 	result, _ := strconv.ParseInt(randStr, 10, 64)
// 	str := strconv.Itoa(int(result))
// 	// "OT" + str[2:] + now.Format("060102") + now.Format("1504")
// 	orderNo := fmt.Sprintf("OT%s%s%s%06d%06d", str, now.Format("060102"), now.Format("1504"), branchID, machineID)
// 	return orderNo
// }

// func GeneratePaymentRequestId12char1(currentTime time.Time) string {
// 	t := currentTime.Format("150405.0000")
// 	number := strings.ReplaceAll(t, ".", "")
// 	randStr := strconv.Itoa(rand.Intn(9)+1) + number[1:]
// 	result, _ := strconv.ParseInt(randStr, 10, 64)
// 	str := strconv.Itoa(int(result))
// 	txRef := "OTR" + str[1:]
// 	return txRef
// }

//	func GeneratePaymentRefCode1(now time.Time) string {
//		nanoPart := fmt.Sprintf("%06d", now.UnixNano()/int64(time.Microsecond)%1_000_000)
//		return "OTR" + now.Format("060102") + nanoPart + now.Format("1504")
//	}
func GenerateOrderNoGen2WithBranch(now time.Time, branchID, machineID int) string {
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
func GenerateOrderNoGen2(now time.Time) string {
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

func GeneratePaymentRefCode2(now time.Time) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 3)
	for i := range b {
		b[i] = letters[int(rand.Int31n(int32(len(letters))))]
	}
	nanoPart := fmt.Sprintf("%06d", now.UnixNano()/int64(time.Microsecond)%1_000_000)
	return "OTR" + now.Format("060102") + now.Format("1504") + string(b) + nanoPart[2:]
}

func GeneratePaymentRequestId12char2(currentTime time.Time) string {
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
func GenerateBranchTransactionNo(now time.Time, branchID, index int, prefix string) string {
	timestamp := now.UnixMilli()
	trimTimestamp := strconv.Itoa(int(timestamp))[8:13]
	yearStr := strconv.Itoa(now.Year())
	last2DigitsOfYear := yearStr[len(yearStr)-2:]
	tranNo := fmt.Sprintf("%s%04d%s%02d%02d%02d%02d%s", prefix, branchID, last2DigitsOfYear, now.Month(), now.Day(), now.Hour(), now.Minute(), trimTimestamp)
	return fmt.Sprintf("%s%02d", tranNo, index)
}

func GeneratePaymentRefCode2WDF(now time.Time) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 3)
	for i := range b {
		b[i] = letters[int(rand.Int31n(int32(len(letters))))]
	}
	nanoPart := fmt.Sprintf("%06d", now.UnixNano()/int64(time.Microsecond)%1_000_000)
	return "WDR" + now.Format("060102") + now.Format("1504") + string(b) + nanoPart[2:]
}

func GeneratePaymentRequestId12char2WDF(currentTime time.Time) string {
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

// // OT + YYMMDD + hhmm + Random char(1) + branchID(5) + machineID(6) + Random number(6)
// // OT-250612-1733-Q-00132-023342-366748
// OT2506121733Q00132023342366748 :len( 30 )

// // OTR + YYMMDD + hhmm + Random char(1) + unix nano(4)
// // OTR-250612-1733-D-6748
// OTR2506121733DIK6748 :len( 20 )

// // OTR + unix nano(6) + Random char(3)
// // OTR-573305-NVA
// OTR373305NVA :len( 12 )

// {
// 	// ของพี่รัต
// 	// C + รหัสสาขา + ปี + เดือน + วัน + ชั่วโมง + นาที + UnixMilli + length of order in cart
// 	// C-0123-25-06-12-17-33-85466-10
// 	C012325061217338546610 :len( 22 )

// 	// WDF + รหัสสาขา + ปี + เดือน + วัน + ชั่วโมง + นาที + UnixMilli + index of order in cart
// 	// WDF-0123-25-06-12-17-33-85466-09
// 	WDF012325061217338546609 :len( 24 )
// }

// // WDR + YYMMDD + hhmm + Random char(3) + unix nano(4)
// // WDR-250612-1733-VEH-6748
// WDR2506121733VEH6748 :len( 20 )

// // WDR + unix nano(6) + Random char(3)
// // WDR-573305-VHO
// WDR573305VHO :len( 12 )

// NOTE:
// {
// 	Random number จาก Nano Second ของเวลาปัจจุบัน
// 	branch จะลองรับได้มากสุด 99,999
// 	machine จะลองรับได้มากสุด 999,999

// 	OTR, WDR
// 	สำหรับ Payment
// 	ที่ขนาด 12 จะสำหรับ Paotang and KTB Next

// 	ของ WDF ไม่เหมือน OT
// 	- สาขานึงจะไม่สามารถสร้างซ้ำกันได้? ถ้าซ้ำ order จะไม่เท่ากัน? ถ้าเท่ากัน?
// 	- จะไม่มี Machine
// }
