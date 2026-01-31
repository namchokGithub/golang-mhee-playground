package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	http.HandleFunc("/run", runHandler)
	http.HandleFunc("/mock-api", mockAPIHandler)

	log.Println("Server started on :9000")
	http.ListenAndServe(":9000", nil)
}

// /run endpoint to trigger concurrent calls
func runHandler(w http.ResponseWriter, r *http.Request) {
	concurrentCalls := 100
	apiURL := "http://localhost:9001/generate"
	// apiURL := "http://localhost:8080/mock-api"

	var wg sync.WaitGroup
	wg.Add(concurrentCalls)

	results := make([]string, concurrentCalls)
	errors := make([]error, concurrentCalls)

	client := &http.Client{Timeout: 5 * time.Second}

	for i := 0; i < concurrentCalls; i++ {
		go func(i int) {
			defer wg.Done()

			req, _ := http.NewRequestWithContext(context.Background(), "GET", apiURL, nil)
			resp, err := client.Do(req)
			if err != nil {
				errors[i] = err
				return
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			results[i] = fmt.Sprintf("Call %d: %s", i+1, string(body))
		}(i)
	}

	wg.Wait()

	for i := 0; i < concurrentCalls; i++ {
		if errors[i] != nil {
			fmt.Fprintf(w, "Call %d failed: %v\n", i+1, errors[i])
		} else {
			fmt.Fprintf(w, "%s\n", results[i])
		}
	}
}

// /mock-api endpoint that simulates an API
func mockAPIHandler(w http.ResponseWriter, r *http.Request) {
	// time.Sleep(500 * time.Millisecond) // simulate delay
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
