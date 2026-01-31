package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/generate", generateHandler)
	fmt.Println("Service B listening on :9001")
	http.ListenAndServe(":9001", nil)
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, GenerateOrderNo3(time.Now(), 121, 32111))
}

func GenerateOrderNo3(now time.Time, branchID, machineID int) string {
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
